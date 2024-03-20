package main

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/aliml92/go-typesense/typesense"

	"github.com/aliml92/anor/auth"
	"github.com/aliml92/anor/config"
	"github.com/aliml92/anor/html"
	"github.com/aliml92/anor/pkg/emailer"
	"github.com/aliml92/anor/pkg/middlewares"
	"github.com/aliml92/anor/postgres"
	"github.com/aliml92/anor/postgres/store"
	"github.com/aliml92/anor/postgres/store/category"
	"github.com/aliml92/anor/postgres/store/product"
	"github.com/aliml92/anor/postgres/store/user"
	"github.com/aliml92/anor/productcatalog"
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
	dbPool, err := store.NewDatabasePool(ctx, &cfg.Database)
	if err != nil {
		return err
	}

	defer dbPool.Close()

	// setup typesense client
	client, _ := typesense.NewClient(nil, "http://localhost:8108")
	client = client.WithAPIKey("xyz")
	searcher := ts.NewSearcher(client)
	// already initialized
	//if err := ts.InitCollections(ctx, client); err != nil {
	//	return err
	//}

	// setup brevo brevoEmailer
	emailTemplate := template.Must(template.ParseGlob("./html/templates/html/email/*.tmpl"))
	brevoEmailer := emailer.NewBrevoEmailer(&cfg.Email, emailTemplate)

	// setup session manager
	// TODO: do not use hardcoded lifetime value
	session := scs.New()
	session.Lifetime = 60 * time.Minute
	session.Store = pgxstore.New(dbPool)

	userStore := user.NewStore(dbPool)
	productStore := product.NewStore(dbPool)
	categoryStore := category.NewStore(dbPool)

	userService := postgres.NewUserService(userStore)
	productService := postgres.NewProductService(productStore, categoryStore)
	categoryService := postgres.NewCategoryService(categoryStore)
	authService := auth.NewAuthService(userService, brevoEmailer)

	// setup templates
	renderer := html.NewRender()

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

	authHandler := auth.NewHandler(authService, renderer, session, logger)
	productcatalogHandler := productcatalog.NewHandler(userService, productService, categoryService, searcher, renderer, logger, session)

	// setup handlers
	mux := http.NewServeMux()

	addStaticRoutes(mux)
	auth.RegisterRoutes(authHandler, mux)
	productcatalog.RegisterRoutes(productcatalogHandler, mux)

	var handler http.Handler = mux
	excluded := []string{
		"/health", "/status", "/static", "/user/static/", "/products/static/", "/categories/static/",
		"/favicon.ico", "/categories/images/",
	}
	handler = middlewares.RequestLogger(handler, logger, excluded...)

	// create server and run
	httpServer := http.Server{
		Addr:     net.JoinHostPort(cfg.Server.Host, cfg.Server.Port),
		Handler:  session.LoadAndSave(handler),
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

func addStaticRoutes(mux *http.ServeMux) {
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	mux.Handle("GET /user/static/", http.StripPrefix("/user/static/", http.FileServer(http.Dir("web/static"))))
	mux.Handle("GET /products/static/", http.StripPrefix("/products/static/", http.FileServer(http.Dir("web/static"))))
	mux.Handle("GET /categories/static/", http.StripPrefix("/categories/static/", http.FileServer(http.Dir("web/static"))))
	mux.Handle("GET /favicon.ico", http.HandlerFunc(faviconHandler))
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/static/images/favicon.ico")
}
