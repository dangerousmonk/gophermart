package config

import (
	"flag"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

const (
	defaultServerAddr = "localhost:8080"
	defaultJWTSecret  = "b6e2490a47c14cb7a1732aed3ba3f3c5"
)

type Config struct {
	ServerAddr  string
	DatabaseDSN string
	JWTSecret   string
	AccrualAddr string
}

func InitConfig() *Config {
	cfg := &Config{}
	if err := godotenv.Load(); err != nil {
		slog.Error("Unable to load envs from file", slog.Any("error", err))

	}
	// Чтение флагов командной строки
	flag.StringVar(&cfg.ServerAddr, "a", cfg.ServerAddr, "Address to run server")
	flag.StringVar(&cfg.AccrualAddr, "r", cfg.AccrualAddr, "Address to run accrual system")
	flag.StringVar(&cfg.DatabaseDSN, "d", cfg.DatabaseDSN, "Database DSN")
	flag.Parse()

	// Чтение переменных окружения
	addr := os.Getenv("RUN_ADDRESS")
	if addr != "" {
		cfg.ServerAddr = addr
	}

	dbDSN := os.Getenv("DATABASE_URI")
	if dbDSN != "" {
		cfg.DatabaseDSN = dbDSN
	}

	accrualAddr := os.Getenv("ACCRUAL_SYSTEM_ADDRESS")
	if accrualAddr != "" {
		cfg.AccrualAddr = "http://" + accrualAddr
	}

	// Инициализация переменных по умолчанию
	if cfg.ServerAddr == "" {
		cfg.ServerAddr = defaultServerAddr
	}

	cfg.JWTSecret = defaultJWTSecret
	return cfg
}
