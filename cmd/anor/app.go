package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/aliml92/anor"
	"github.com/aliml92/anor/brevo"
	"github.com/aliml92/anor/storefront/cart"
	"github.com/aliml92/anor/storefront/checkout"
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
	cartRepo "github.com/aliml92/anor/postgres/repository/cart"
	categoryRepo "github.com/aliml92/anor/postgres/repository/category"
	orderRepo "github.com/aliml92/anor/postgres/repository/order"
	productRepo "github.com/aliml92/anor/postgres/repository/product"
	userRepo "github.com/aliml92/anor/postgres/repository/user"
	"github.com/aliml92/anor/redis"
	rs "github.com/aliml92/anor/redis/cache"
	authCache "github.com/aliml92/anor/redis/cache/auth"
	"github.com/aliml92/anor/redis/cache/session"
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

	// setup session manager
	sessionManager := session.NewManager(cfg.Session, redisClient)

	userRepository := userRepo.NewRepository(dbPool)
	productRepository := productRepo.NewRepository(dbPool)
	categoryRepository := categoryRepo.NewRepository(dbPool)
	cartRepository := cartRepo.NewRepository(dbPool)
	orderRepository := orderRepo.NewRepository(dbPool)

	userService := postgres.NewUserService(userRepository, cartRepository, orderRepository)
	productService := postgres.NewProductService(productRepository, categoryRepository)
	categoryService := postgres.NewCategoryService(categoryRepository, categoryCacher)
	cartService := postgres.NewCartService(cartRepository, productRepository)
	orderService := postgres.NewOrderService(orderRepository)
	authService := auth.NewAuthService(userService, brevoEmailer, otpCacher, rpCacher)

	view := html.NewView()

	logger := newSlogLogger(cfg.Logger)
	slog.SetDefault(logger)

	authHandler := auth.NewHandler(authService, cartService, view, sessionManager, logger)
	productcatalogHandler := product.NewHandler(userService, productService, categoryService, cartService, searcher, view, logger, sessionManager)
	userHandler := user.NewHandler(userService, cartService, view, sessionManager, logger)
	cartHandler := cart.NewHandler(userService, cartService, view, sessionManager, logger)
	checkoutHandler := checkout.NewHandler(userService, cartService, orderService, view, sessionManager, logger, cfg)

	mux := anor.NewRouter()

	addStaticRoutes(mux)

	mux.Use(
		//middlewares.RequestID,
		sloghttp.Recovery,
		sloghttp.New(logger.WithGroup("http")),
		sessionManager.Auth.LoadAndSave,
		userHandler.AuthInjector,
		authHandler.SessionMiddleware,
	)

	auth.RegisterRoutes(authHandler, mux)
	product.RegisterRoutes(productcatalogHandler, mux)
	user.RegisterRoutes(userHandler, mux)
	cart.RegisterRoutes(cartHandler, mux)
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
