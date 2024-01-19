package app

import (
	"errors"
	"log/slog"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

const (
	_defaultAttempts = 20
	_defaultTimeout  = time.Second
)

func init() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Add path to to .env file when main moves to cmd/app
	if err := godotenv.Load(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	dbURL, ok := os.LookupEnv("DB_URL")
	if !ok || len(dbURL) == 0 {
		logger.Error("migrate", "error", "environment variable not declared: DB_dbURL")
		os.Exit(1)
	}

	sslMode, ok := os.LookupEnv("SSL_MODE")
	if !ok || len(sslMode) == 0 {
		sslMode = "disable"
	}

	dbURL += "?sslmode=" + sslMode

	var (
		attempts = _defaultAttempts
		err      error
		m        *migrate.Migrate
	)

	for attempts > 0 {
		m, err = migrate.New("file://migrations", dbURL)
		if err == nil {
			break
		}

		logger.Warn("migrate", "postgres is trying to connect, attempts left", attempts)
		time.Sleep(_defaultTimeout)
		attempts--
	}

	if err != nil {
		logger.Error("migrate", "postgres connection error", err.Error())
		os.Exit(1)
	}

	err = m.Up()
	defer m.Close()

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.Error("migrate", "up error", err.Error())
		os.Exit(1)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		logger.Warn("migrate", "up", "no change")
		return
	}

	logger.Info("migrate", "up", "success")
}
