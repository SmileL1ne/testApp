package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"testApp/internal/handlers"
	"testApp/internal/repository"
	"testApp/internal/service"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func Run() {
	// Initialized logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Load environment variables to from .env
	err := godotenv.Load() // Add path to to .env file when main moves to cmd/app
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Get env variables from .env
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSslMode := os.Getenv("DB_SSLMODE")

	dsn := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		dbHost, dbPort, dbName, dbUser, dbPassword, dbSslMode)

	// Database connection initialization
	conn, err := openDB(dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	r := repository.NewRepository(conn)
	s := service.NewService(r)

	// Server initialization
	server := &http.Server{
		Addr:     "127.0.0.1:" + os.Getenv("PORT"),
		Handler:  handlers.NewRouter(logger, s),
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	// Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

		sig := <-sigCh
		logger.Info("signal received", "signal", sig.String())

		if err := server.Shutdown(context.Background()); err != nil {
			logger.Error("server shutdown", "error", err)
			os.Exit(1)
		}
		if err := conn.Close(context.Background()); err != nil {
			logger.Error("connection close", "error", err)
			os.Exit(1)
		}

		os.Exit(0)
	}()

	// Starting the server
	logger.Info("starting the server", "addr", server.Addr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("listen and serve", "error", err.Error())
		os.Exit(1)
	}
}

func openDB(dsn string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(context.Background()); err != nil {
		return nil, err
	}

	return conn, nil
}
