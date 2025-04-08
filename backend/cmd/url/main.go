package main

import (
	"REST_project/config"
	"REST_project/internal/handlers/create-handlers"
	"REST_project/internal/handlers/logger"
	"REST_project/internal/handlers/register-handlers"
	"REST_project/internal/storage"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
	"github.com/rs/cors"
	"os/signal"
	"syscall"
	"time"
)

/*
migrate: migrate create -ext sql -dir migrations -seq create_users_table
*/

const (
	envProd = "prod"
	envDev  = "Dev"
)

func main() {
	//TODO: config +++
	//TODO: init logger +++
	//TODO: data base (users, posts, comments){
	//	1. migrations ++
	//	2. connect to postgres+docker ++
	// }
	//TODO: route and serve home page ++
	//TODO: сделать обработку ошибки повтороной записи в базу данных при регистрации
	log := SetupLogger(envDev) //наверное, стоит добавить в config
	log.Info(
		"starting server",
		slog.String("env", envDev),
		slog.String("version", "0.0.1"),
	)
	cfg := config.MustLoad()
	db, err := storage.New(cfg.DBConf)
	if err != nil {
		log.Error("failed to init storage", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})
		os.Exit(1)
	}
	if err = storage.RunMigrations(db.DB); err != nil {
		log.Error("failed to make migrations", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})
	}

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(corsMiddleware.Handler)

	// Маршруты регистрации
	router.Route("/register", func(r chi.Router) {
		r.Post("/enterprise", register_handlers.RegisterEnterprise(log, db))
		r.Get("/enterprise", register_handlers.GetEnterprises(log, db))
		r.Post("/event", register_handlers.RegisterEvent(log, db))
		r.Get("/event", register_handlers.GetEvents(log, db)) 
		r.Post("/user", register_handlers.RegisterUser(log, db))
		r.Get("/user", register_handlers.GetUsers(log, db))
	})


	// Маршруты для работы с постами и комментариями
	router.Route("/api", func(r chi.Router) {
		r.Post("/posts", create_handlers.CreatePost(log, db))
		r.Get("/posts", create_handlers.GetPosts(log, db)) 
		r.Post("/comments", create_handlers.CreateComment(log, db))
		r.Get("/comments", create_handlers.GetComments(log, db)) 
	})

	// Health check endpoint
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Info("starting server", slog.String("address", cfg.ServConf.HostREST))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         cfg.ServConf.HostREST,
		Handler:      router,
		ReadTimeout:  cfg.ServConf.Timeout,
		WriteTimeout: cfg.ServConf.Timeout,
		IdleTimeout:  cfg.ServConf.Timeout * 3,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server")
		}
	}()

	log.Info("server started")

	<-done
	log.Info("stopping server")

	// TODO: move timeout to config
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", slog.Attr{
			Key:   "Error",
			Value: slog.StringValue(err.Error()),
		})
		return
	}
	log.Info("gracefully stopped")
}

func SetupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	default:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
