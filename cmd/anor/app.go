package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/aliml92/anor"
	"github.com/aliml92/anor/brevo"
	"github.com/aliml92/anor/html/templates/shared/header"
	"github.com/aliml92/anor/html/templates/shared/header/components"
	"github.com/aliml92/anor/middlewares"
	redissession "github.com/aliml92/anor/redis/session"
	"github.com/aliml92/anor/session"
	"github.com/aliml92/anor/storefront/cart"
	"github.com/aliml92/anor/storefront/checkout"
	"github.com/gorilla/sessions"
	_ "github.com/joho/godotenv/autoload"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"github.com/samber/oops"
	sloghttp "github.com/samber/slog-http"
	"io"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/aliml92/anor/config"
	"github.com/aliml92/anor/html"
	"github.com/aliml92/anor/postgres"
	addressRepo "github.com/aliml92/anor/postgres/repository/address"
	cartRepo "github.com/aliml92/anor/postgres/repository/cart"
	categoryRepo "github.com/aliml92/anor/postgres/repository/category"
	featuredSelectionRepo "github.com/aliml92/anor/postgres/repository/featured_selection"
	orderRepo "github.com/aliml92/anor/postgres/repository/order"
	paymentRepo "github.com/aliml92/anor/postgres/repository/payment"
	productRepo "github.com/aliml92/anor/postgres/repository/product"
	userRepo "github.com/aliml92/anor/postgres/repository/user"
	"github.com/aliml92/anor/redis"
	rs "github.com/aliml92/anor/redis/cache"
	authCache "github.com/aliml92/anor/redis/cache/auth"
	"github.com/aliml92/anor/storefront/auth"
	"github.com/aliml92/anor/storefront/product"
	"github.com/aliml92/anor/storefront/user"
	"github.com/aliml92/anor/typesense"
)

func run(ctx context.Context, w io.Writer, args []string) error {
	// get config
	fmt.Println("Starting anor...")
	cfg, err := config.New()
	if err != nil {
		return err
	}

	// get db connection pool
	dbPool, err := postgres.NewDatabasePool(ctx, &cfg.Database)
	if err != nil {
		return err
	}
	defer dbPool.Close()

	// get typesense client
	tsClient, err := typesense.NewClient(ctx, cfg.Typesense)
	if err != nil {
		return err
	}
	if err := typesense.InitCollections(ctx, tsClient); err != nil {
		return err
	}

	// get redis client
	redisClient, err := redis.NewClient(ctx, cfg.Redis)
	if err != nil {
		return err
	}
	categoryCacher := rs.NewCategoryCacher(redisClient)
	rpCacher := authCache.NewResetPasswordTokenCacher(redisClient)
	otpCacher := authCache.NewSignupConfirmationOTPCacher(redisClient)
	searcher := typesense.NewSearcher(tsClient)

	// setup brevo emailer
	brevoEmailer := brevo.NewEmailer(cfg.Email)

	redisStore := redissession.NewRedisStore(redisClient).WithPrefix("keys:session:")

	// setup session manager
	sessionManager := redissession.NewManager(
		redissession.WithCookieName(cfg.Session.CookieName),
		redissession.WithAuthLifetime(cfg.Session.AuthLifetime),
		redissession.WithGuestLifetime(cfg.Session.GuestLifetime),
		redissession.WithGuestSkipPaths(map[string][]string{
			"/auth/signin":     {"POST"},
			"/auth/signup":     {"POST"},
			"/static":          {"GET"},
			"/auth/static":     {"GET"},
			"/products/static": {"GET"},
			"/category/static": {"GET"},
			"/user/static":     {"GET"},
			"/checkout/static": {"GET"},
			"/favicon.ico":     {"GET"},
		}),
		redissession.WithStore(redisStore),
		redissession.WithCodec(redissession.MessagePackCodec{}),
	)

	store := sessions.NewCookieStore([]byte("My_Secure_key"))
	store.MaxAge(86400 * 30)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = false
	store.Options.SameSite = http.SameSiteLaxMode

	gothic.Store = store
	// setup goth store
	oauthCfg := cfg.GoogleOAuth
	goth.UseProviders(
		google.New(oauthCfg.ClientID, oauthCfg.ClientSecret, oauthCfg.RedirectURL, oauthCfg.Scopes...),
	)

	userRepository := userRepo.NewRepository(dbPool)
	productRepository := productRepo.NewRepository(dbPool)
	categoryRepository := categoryRepo.NewRepository(dbPool)
	featuredSelectionRepository := featuredSelectionRepo.NewRepository(dbPool)
	cartRepository := cartRepo.NewRepository(dbPool)
	orderRepository := orderRepo.NewRepository(dbPool)
	addressRepository := addressRepo.NewRepository(dbPool)
	paymentRepository := paymentRepo.NewRepository(dbPool)

	userService := postgres.NewUserService(userRepository, cartRepository, orderRepository)
	productService := postgres.NewProductService(productRepository, categoryRepository)
	categoryService := postgres.NewCategoryService(categoryRepository, categoryCacher)
	featuredSelectionService := postgres.NewFeaturedSelectionService(featuredSelectionRepository, productRepository, categoryRepository)
	cartService := postgres.NewCartService(cartRepository, productRepository)
	addressService := postgres.NewAddressService(addressRepository)
	orderService := postgres.NewOrderService(orderRepository)
	stripePaymentService := postgres.NewStipePaymentService(paymentRepository)

	authServiceCfg := auth.ServiceConfig{
		UserService:                 userService,
		Emailer:                     brevoEmailer,
		SignupConfirmationOTPCacher: otpCacher,
		ResetPasswordTokenCacher:    rpCacher,
	}
	authService := auth.NewAuthService(authServiceCfg)

	view := html.NewView()
	getHeaderDataFunc := func(ctx context.Context) (header.Base, error) {
		u := session.UserFromContext(ctx)
		rc, err := categoryService.GetRootCategories(ctx)
		if err != nil {
			return header.Base{}, err
		}

		h := header.Base{
			User:           *u,
			RootCategories: rc,
		}

		if u.CartID != 0 {
			cartItemCount, err := cartService.CountItems(ctx, u.CartID)
			if err != nil {
				return header.Base{}, err
			}
			h.CartNavItem = components.CartNavItem{CartItemsCount: int(cartItemCount)}
		}

		return h, nil
	}

	logger := newSlogLogger(cfg.Logger)
	slog.SetDefault(logger)
	oops.StackTraceMaxDepth = 2

	authCfg := &auth.HandlerConfig{
		AuthService: authService,
		CartService: cartService,
		Session:     sessionManager,
		View:        view,
		Logger:      logger,
	}

	authHandler := auth.NewHandler(authCfg)

	productCfg := &product.HandlerConfig{
		UserService:             userService,
		ProductService:          productService,
		CategoryService:         categoryService,
		FeatureSelectionService: featuredSelectionService,
		CartService:             cartService,
		Session:                 sessionManager,
		Searcher:                searcher,
		View:                    view,
		Logger:                  logger,
		GetHeaderDataFunc:       getHeaderDataFunc,
	}
	productcatalogHandler := product.NewHandler(productCfg)

	userCfg := &user.HandlerConfig{
		UserService:       userService,
		OrderService:      orderService,
		AddressService:    addressService,
		Session:           sessionManager,
		View:              view,
		Logger:            logger,
		Config:            cfg,
		GetHeaderDataFunc: getHeaderDataFunc,
	}
	userHandler := user.NewHandler(userCfg)

	cartConfig := &cart.HandlerConfig{
		UserService:       userService,
		CartService:       cartService,
		CategoryService:   categoryService,
		Session:           sessionManager,
		View:              view,
		Logger:            logger,
		GetHeaderDataFunc: getHeaderDataFunc,
	}
	cartHandler := cart.NewHandler(cartConfig)

	checkoutHandlerCfg := checkout.HandlerConfig{
		UserService:          userService,
		CartService:          cartService,
		OrderService:         orderService,
		CategorySvc:          categoryService,
		AddressService:       addressService,
		StripePaymentService: stripePaymentService,
		View:                 view,
		SessionManager:       sessionManager,
		Logger:               logger,
		Config:               cfg,
	}
	checkoutHandler := checkout.NewHandler(checkoutHandlerCfg)

	mux := anor.NewRouter()

	addStaticRoutes(mux)

	//corsConfig := middlewares.DefaultCORSConfig()
	//corsConfig.AllowedOrigins = []string{"*"}

	mux.Use(
		//middlewares.RequestID,
		sloghttp.New(logger.WithGroup("http")),
		sloghttp.Recovery,
		sessionManager.LoadAndSave,
		sessionManager.LoadAndSaveGuest,
		sessionManager.LoadUser,
		//middlewares.CORSMiddleware(corsConfig),
	)

	auth.RegisterRoutes(authHandler, mux)
	cart.RegisterRoutes(cartHandler, mux)
	product.RegisterRoutes(productcatalogHandler, mux)
	mux.Use(middlewares.RequireAuth)
	user.RegisterRoutes(userHandler, mux)
	checkout.RegisterRoutes(checkoutHandler, mux)

	var handler http.Handler = mux

	// create server and run
	httpServer := http.Server{
		Addr:     net.JoinHostPort(cfg.Server.Host, cfg.Server.Port),
		Handler:  handler,
		ErrorLog: newServerLogger(cfg.Logger),
	}
	go func() {
		log.Printf("listening on %s\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
		}
	}()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		if err := httpServer.Shutdown(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
		}
	}()
	wg.Wait()
	return nil
}

func addStaticRoutes(mux *anor.Router) {
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	mux.Handle("GET /auth/static/", http.StripPrefix("/auth/static/", http.FileServer(http.Dir("web/static"))))
	mux.Handle("GET /user/static/", http.StripPrefix("/user/static/", http.FileServer(http.Dir("web/static"))))
	mux.Handle("GET /products/static/", http.StripPrefix("/products/static/", http.FileServer(http.Dir("web/static"))))
	mux.Handle("GET /categories/static/", http.StripPrefix("/categories/static/", http.FileServer(http.Dir("web/static"))))
	mux.Handle("GET /checkout/static/", http.StripPrefix("/checkout/static/", http.FileServer(http.Dir("web/static"))))
	mux.Handle("GET /favicon.ico", http.HandlerFunc(faviconHandler))
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/static/images/favicon.ico")
}

func newServerLogger(cfg config.LoggerConfig) *log.Logger {
	timeFormat, ok := timeFormats[cfg.TimeFormat]
	if !ok {
		timeFormat = time.StampMilli
	}

	handlerOptions := &slog.HandlerOptions{
		Level:     slog.LevelError,
		AddSource: cfg.AddSource,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{
					Key:   a.Key,
					Value: slog.StringValue(a.Value.Time().Format(timeFormat)),
				}
			}
			return a
		},
	}

	var handler slog.Handler
	switch strings.ToLower(cfg.Format) {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, handlerOptions)
	case "text":
		handler = slog.NewTextHandler(os.Stdout, handlerOptions)
	default:
		handler = slog.NewTextHandler(os.Stdout, handlerOptions)
	}
	logger := slog.NewLogLogger(handler, slog.LevelError)

	return logger
}

func newSlogLogger(cfg config.LoggerConfig) *slog.Logger {
	var level slog.Level
	switch strings.ToLower(cfg.Level) {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo // Default to Info if invalid level is specified
	}

	timeFormat, ok := timeFormats[cfg.TimeFormat]
	if !ok {
		timeFormat = time.StampMilli
	}

	handlerOptions := &slog.HandlerOptions{
		Level:     level,
		AddSource: cfg.AddSource,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{
					Key:   a.Key,
					Value: slog.StringValue(a.Value.Time().Format(timeFormat)),
				}
			}
			return a
		},
	}

	var handler slog.Handler
	switch strings.ToLower(cfg.Format) {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, handlerOptions)
	case "text":
		handler = slog.NewTextHandler(os.Stdout, handlerOptions)
	default:
		handler = slog.NewTextHandler(os.Stdout, handlerOptions)
	}

	return slog.New(handler)
}

var timeFormats = map[string]string{
	"Layout":      time.Layout,
	"ANSIC":       time.ANSIC,
	"UnixDate":    time.UnixDate,
	"RubyDate":    time.RubyDate,
	"RFC822":      time.RFC822,
	"RFC822Z":     time.RFC822Z,
	"RFC850":      time.RFC850,
	"RFC1123":     time.RFC1123,
	"RFC1123Z":    time.RFC1123Z,
	"RFC3339":     time.RFC3339,
	"RFC3339Nano": time.RFC3339Nano,
	"Kitchen":     time.Kitchen,
	"Stamp":       time.Stamp,
	"StampMilli":  time.StampMilli,
	"StampMicro":  time.StampMicro,
	"StampNano":   time.StampNano,
	"DateTime":    time.DateTime,
	"DateOnly":    time.DateOnly,
	"TimeOnly":    time.TimeOnly,
}
