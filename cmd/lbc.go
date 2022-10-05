package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"github.com/xefiji/lbc/lbc"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Warn().Msg("no env file loaded")
	}
}

func main() {
	if err := run(); err != nil {
		log.Error().Err(err).Msg("program exiting with error")
		os.Exit(1)
	}
}

func run() error {
	return lbc.Crawl(
		lbc.WithRapidAPI(
			env("API_URL", ""),
			env("RAPIDAPI_HOST", ""),
			env("RAPIDAPI_KEY", ""),
		),
		lbc.WithRedis(
			env("REDIS_HOST", "localhost"),
			env("REDIS_PORT", "6379"),
			env("REDIS_PASSWORD", ""),
			env("REDIS_DB", "0"),
		),
		lbc.WithUsers(env("USERS", "")),
		lbc.WithMailJet(
			env("MAILJET_KEY", ""),
			env("MAILJET_SECRET", ""),
			env("MAIL_FROM", ""),
		),
	)
}

func env(name, fallback string) string {
	if val, ok := os.LookupEnv(name); ok {
		return val
	}

	return fallback
}
