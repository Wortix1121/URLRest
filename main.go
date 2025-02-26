package main

import (
	"api/internal/config"
	"api/internal/httpServer/handlers/del"
	"api/internal/httpServer/handlers/redirect"
	"api/internal/httpServer/handlers/save"
	"api/internal/storage/postgres/rest"
	"os"

	"fmt"
	"net/http"

	_ "api/docs"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	chiSwagger "github.com/swaggo/http-swagger/v2"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// @title Doc API(URL)
// @version 0.1.1
// @Description Url REST

// @contact.name Andrey
// @contract.url https://github.com/Wortix1121/URLRest

// @host localhost:8000
// @BasePath
func main() {
	// TODO: inti config - Cleanenv
	cfg := config.MustLoad()
	fmt.Println(cfg)

	// TODO: init logger - Zap
	log := setupLogger(cfg.Env)

	log.Info("Start", zap.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	// TODO: init storage - postgre

	// newdb, err := postgres.NewTableUrls(context.Background())
	// if err != nil {
	// 	log.Fatal("failed to create table", zap.Error(err))
	// 	os.Exit(1)
	// }

	// _ = newdb

	// exampleSave, err := rest.SaveUrl("https://google.com", "google")
	// if err != nil {
	// 	log.Fatal("failed to request db", zap.Error(err))
	// 	os.Exit(1)
	// }
	// // "https://google.com", "google"

	// _ = exampleSave

	// exampleSave, err := rest.SaveUrl("https://www.google.ru/?hl=ru", "google", context.Background())
	// if err != nil {
	// 	log.Fatal("failed to request db", zap.Error(err))
	// 	os.Exit(1)
	// } else {
	// 	log.Info("Nice request to db")
	// }

	// a := exampleSave

	// b, _ := json.Marshal(a)
	// fmt.Println(string(b))

	// TODO: init storage - sqlite

	// storage, err := sqlite.New(cfg.StoragePath)
	// if err != nil {
	// 	log.Fatal("failed to init storage", zap.Error(err))
	// 	os.Exit(1)
	// }

	// _ = storage

	storage, err := rest.Connect()
	if err != nil {
		log.Fatal("failed to init storage", zap.Error(err))
		os.Exit(1)
	}

	// TODO: init router - Chi

	router := chi.NewRouter()

	// middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	// router.Use(middleware.New(log))
	router.Use(middleware.Recoverer)
	// router.Use(middleware.URLFormat)

	// Authorization

	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-rest", map[string]string{
			cfg.Authorization.User: cfg.Authorization.Password,
			cfg.Authorization.User: cfg.Authorization.Password,
			cfg.Authorization.User: cfg.Authorization.Password,
		}))

		r.Post("/save", (save.SaveHand(log, storage)))
		r.Delete("/urldel/{alias}", del.DelHand(log, storage))
	})

	router.Get("/swagger/*", chiSwagger.Handler(
		chiSwagger.URL("http://localhost:8000/swagger/doc.json"), //The url pointing to API definition
	))
	//handlers
	router.Get("/{alias}", redirect.GetHand(log, storage))
	// router.Post("/save", save.SaveHand(log, storage))
	// router.Delete("/url/{alias}", redirect.DelHand(log, storage))

	log.Info("Starting server", zap.String("address", cfg.Address))

	serv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := serv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")

	// TODO: init server -
}

func setupLogger(env string) *zap.Logger {
	var log *zap.Logger

	switch env {
	case envLocal:
		log, _ = zap.NewDevelopment()
		defer log.Sync()

	case envDev:
		log = zap.NewExample()
		defer log.Sync()

	case envProd:
		log, _ = zap.NewProduction()
		defer log.Sync()

	}

	return log

}
