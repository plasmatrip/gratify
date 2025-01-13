package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

const (
	port     = "8080"
	host     = "localhost"
	accrual  = ""
	database = "postgres://gratify:password@localhost:5432/gratify?sslmode=disable"
)

type Config struct {
	Host        string `env:"RUN_ADDRESS"`
	Database    string `env:"DATABASE_URI"`
	Accrual     string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	TokenSecret string `env:"TOKEN_SECRET"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}

	ex, err := os.Executable()
	if err != nil {
		return nil, err
	}

	if err := godotenv.Load(filepath.Dir(ex) + "/.env"); err != nil {
		return nil, errors.New(".env not found")
	}

	cl := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	// читаем переменные окружения, при ошибке прокидываем ее наверх
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to read environment variable: %w", err)
	}

	var fHost string
	cl.StringVar(&fHost, "a", host+":"+port, "server address host:port")

	var fDatabase string
	cl.StringVar(&fDatabase, "d", database, "")

	var fAccrual string
	cl.StringVar(&fAccrual, "r", accrual, "")

	if err := cl.Parse(os.Args[1:]); err != nil {
		return nil, fmt.Errorf("failed to parse flags: %w", err)
	}

	if _, exist := os.LookupEnv("RUN_ADDRESS"); !exist {
		cfg.Host = fHost
	}

	if _, exist := os.LookupEnv("DATABASE_URI"); !exist {
		cfg.Database = fDatabase
	}

	if _, exist := os.LookupEnv("ACCRUAL_SYSTEM_ADDRESS"); !exist {
		cfg.Accrual = fAccrual
	}

	if err := parseAddress(cfg); err != nil {
		return nil, fmt.Errorf("port parsing error: %w", err)
	}

	return cfg, nil
}

func parseAddress(cfg *Config) error {
	args := strings.Split(cfg.Host, ":")
	if len(args) == 2 {
		if len(args[0]) == 0 || len(args[1]) == 0 {
			cfg.Host = host + ":" + port
			return nil
		}

		_, err := strconv.ParseInt(args[1], 10, 64)
		return err
	}
	cfg.Host = host + ":" + port
	return nil
}
