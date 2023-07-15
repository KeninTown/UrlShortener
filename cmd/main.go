package main

import (
	"net/http"
	"os"
	"urlShortener/internal/config"
	save "urlShortener/internal/httpServer/handlers/url"
	mwLogger "urlShortener/internal/httpServer/middleware"
	slogpretty "urlShortener/internal/lib/loggers/handlers/slogPretty"
	"urlShortener/internal/lib/loggers/sl"

	"urlShortener/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/exp/slog"
)

const (
	LOCAL string = "local"
	DEV   string = "dev"
	PROD  string = "prod"
)

func main() {
	// считать конфиг
	cfg := config.MustLoad()

	//TODO сделать логгер
	log := initLogger(cfg.Env)

	//TODO создать DB в sqlite
	db, err := storage.New(cfg.StoragePath)
	if err != nil {
		log.Error("Failed to init storage", sl.Error(err))
		os.Exit(1)
	}

	//TODO init router

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Post("/url", save.New(log, db))

	srv := &http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}
	_ = srv
	// fmt.Printf("%+v", cfg)
	log.Info("Server listening", slog.String("addres", cfg.HttpServer.Address))

	// TODO: run server
	if err = srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")
}

func initLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case LOCAL:
		log = setupPrettySlog()
	case DEV:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case PROD:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
