package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"github.com/xefiji/lbc/lbc"
)

func init() {
	if err := godotenv.Load(); err == nil {
		log.Info().Msg("env file loaded")
	}
}

func main() {
	if err := run(); err != nil {
		log.Error().Err(err).Msg("program exiting with error")
		os.Exit(1)
	}
}

func run() error {
	return lbc.Listen(
		lbc.WithRedis(
			lbc.Env("REDIS_HOST", "localhost"),
			lbc.Env("REDIS_PORT", "6379"),
			lbc.Env("REDIS_PASSWORD", ""),
			lbc.Env("REDIS_DB", "0"),
		),
		lbc.WithServer(
			lbc.Env("SERVER_PORT", "8080"),
		),
	)
}
