package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"vault/internal/config"
	"vault/internal/db"
	"vault/internal/root"
	"vault/internal/user"
	"vault/pkg/database/postgresql"
	"vault/pkg/lib/logger/sl"
	mwLogger "vault/pkg/lib/middleware"
	"vault/pkg/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.MustLoad()

	log := logger.CreateLogger(cfg.Env)

	log.Info("logger successfully setup", slog.String("env", cfg.Env))
	log.Debug("debug messages are available")
	log.Info("info messages are available")
	log.Warn("warn messages are available")
	log.Error("error messages are available")

	databaseConfig := postgresql.DatabaseConfig{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		Name:     cfg.Database.Name,
		SSLMode:  cfg.Database.SSLMode,
		Attemps:  cfg.Database.Attemps,
		Timeout:  cfg.Database.Timeout,
		Delay:    cfg.Database.Delay,
	}

	pool := postgresql.NewClient(context.TODO(), log, databaseConfig)
	if pool == nil {
		os.Exit(1)
	}
	log.Info("database successfully conected")

	dbClient := db.NewClient(pool)

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(mwLogger.New(log))
	log.Info("middleware successfully conected")

	router.Route("/root", root.AddRootRouter(router, dbClient, log, cfg))
	router.Route("/user", user.AddUserRouter(router, dbClient, log, cfg))

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.HTTPServer.Host, cfg.HTTPServer.Port),
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	log.Info("Server statup", "host", cfg.HTTPServer.Host, "port", cfg.HTTPServer.Port)
	if err := srv.ListenAndServe(); err != nil {
		log.Error("Failed to start server", sl.Err(err))
		return
	}
}
