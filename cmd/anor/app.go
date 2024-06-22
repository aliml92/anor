package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/aliml92/anor"
	"github.com/aliml92/anor/storefront/cart"
	"github.com/aliml92/anor/storefront/checkout"
	"html/template"
	"io"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/aliml92/go-typesense/typesense"
	"github.com/redis/go-redis/v9"

	"github.com/aliml92/anor/config"
	"github.com/aliml92/anor/html"
	"github.com/aliml92/anor/pkg/emailer"
	"github.com/aliml92/anor/pkg/middlewares"
	"github.com/aliml92/anor/postgres"
	"github.com/aliml92/anor/postgres/repository"
	cartRepo "github.com/aliml92/anor/postgres/repository/cart"
	categoryRepo "github.com/aliml92/anor/postgres/repository/category"
	orderRepo "github.com/aliml92/anor/postgres/repository/order"
	productRepo "github.com/aliml92/anor/postgres/repository/product"
	userRepo "github.com/aliml92/anor/postgres/repository/user"
	rs "github.com/aliml92/anor/redis/cache"
	authCache "github.com/aliml92/anor/redis/cache/auth"
	"github.com/aliml92/anor/redis/cache/session"
	"github.com/aliml92/anor/storefront/auth"
	"github.com/aliml92/anor/storefront/product"
	"github.com/aliml92/anor/storefront/user"
	ts "github.com/aliml92/anor/typesense"
)

const (
	ConfigFilePathEnvVar = "CONFIG_FILEPATH"
)

func run(ctx context.Context, w io.Writer, args []string) error {
	// get config
	cfg, err := getConfig()
	if err != nil {
		return err
	}

	// get db connection pool
	dbPool, err := repository.NewDatabasePool(ctx, &cfg.Database)
	if err != nil {
		return err
	}

	defer dbPool.Close()

	// setup typesense client
	tsClient, _ := typesense.NewClient(nil, "http://localhost:8108")
	tsClient = tsClient.WithAPIKey("xyz")
	redisClient, err := getRedisClient(ctx, cfg.Redis)
	if err != nil {
		return err
	}

	categoryCacher := rs.NewCategoryCacher(redisClient)
	rpCacher := authCache.NewResetPasswordTokenCacher(redisClient)
	otpCacher := authCache.NewSignupConfirmationOTPCacher(redisClient)
	searcher := ts.NewSearcher(tsClient)

	// already initialized
	if err := ts.InitCollections(ctx, tsClient); err != nil {
		return err
	}

	// setup brevo brevoEmailer
	emailTemplate := template.Must(template.ParseGlob("./pkg/emailer/templates/*.gohtml"))
	brevoEmailer := emailer.NewBrevoEmailer(&cfg.Email, emailTemplate)

	// setup session manager
	// TODO: do not use hardcoded lifetime value
	authSession := scs.New()
	authSession.Store = session.NewRedisStore(redisClient).WithPrefix("keys:app:session:user:authenticated:")
	authSession.Codec = session.MessagePackCodec{}
	authSession.Lifetime = 2 * 24 * time.Hour
	authSession.Cookie.Name = "__anor_ust" // anor, authenticated user's session token

	anonSession := scs.New()
	anonSession.Store = session.NewRedisStore(redisClient).WithPrefix("keys:app:session:user:anonymous:")
	anonSession.Codec = session.MessagePackCodec{}
	anonSession.Lifetime = 5 * 24 * time.Hour
	anonSession.Cookie.Name = "__anor_gst" // anor, guest's session token

	sessionManager := session.NewManager(authSession, anonSession)

	userRepository := userRepo.NewRepository(dbPool)
	productRepository := productRepo.NewRepository(dbPool)
	categoryRepository := categoryRepo.NewRepository(dbPool)
	cartRepository := cartRepo.NewRepository(dbPool)
	orderRepository := orderRepo.NewRepository(dbPool)

	userService := postgres.NewUserService(userRepository, cartRepository, orderRepository)
	productService := postgres.NewProductService(productRepository, categoryRepository)
	categoryService := postgres.NewCategoryService(categoryRepository, categoryCacher)
	authService := auth.NewAuthService(userService, brevoEmailer, otpCacher, rpCacher)
	cartService := postgres.NewCartService(cartRepository, productRepository)
	orderService := postgres.NewOrderService(orderRepository)

	// setup loggers
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: false,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				a.Value = slog.Int64Value(time.Now().Unix())
			}
			return a
		},
	})
	srvLogger := slog.NewLogLogger(logHandler, slog.LevelError)
	logger := slog.New(logHandler)

	view := html.NewView()

	authHandler := auth.NewHandler(authService, cartService, view, sessionManager, logger)
	productcatalogHandler := product.NewHandler(userService, productService, categoryService, cartService, searcher, view, logger, sessionManager)
	userHandler := user.NewHandler(userService, cartService, view, sessionManager, logger)
	cartHandler := cart.NewHandler(userService, cartService, view, sessionManager, logger)
	checkoutHandler := checkout.NewHandler(userService, cartService, orderService, view, sessionManager, logger)

	mux := anor.NewRouter()

	addStaticRoutes(mux)

	responseLogger := middlewares.NewRequestLogger(logger, nil)
	mux.Use(middlewares.RequestID, responseLogger.Log, authSession.LoadAndSave, userHandler.AuthInjector, authHandler.SessionMiddleware)
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
		ErrorLog: srvLogger,
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

func getRedisClient(ctx context.Context, cfg config.RedisConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Username:     cfg.Username,
		Password:     cfg.Password,
		DB:           cfg.DB,
		MinIdleConns: cfg.MinIdleConns,
		MaxIdleConns: cfg.MaxIdleConns,
	})

	statusCmd := rdb.Ping(ctx)
	if statusCmd.Err() != nil {
		return nil, statusCmd.Err()
	}

	return rdb, nil
}

func getConfig() (*config.Config, error) {
	cfgPath := os.Getenv(ConfigFilePathEnvVar)
	if cfgPath == "" {
		return nil, fmt.Errorf("environment variable %s not set or empty", ConfigFilePathEnvVar)
	}

	cfgFile, err := config.LoadConfigFromFile(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config file: %v", err)
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %v", err)
	}

	return cfg, nil
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
