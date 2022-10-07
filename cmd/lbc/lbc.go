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
	return lbc.Crawl(
		lbc.WithDryRun(
			lbc.Env("DRY_RUN", "true"),
		),
		lbc.WithRapidAPI(
			lbc.Env("API_URL", ""),
			lbc.Env("RAPIDAPI_HOST", ""),
			lbc.Env("RAPIDAPI_KEY", ""),
		),
		lbc.WithRedis(
			lbc.Env("REDIS_HOST", "localhost"),
			lbc.Env("REDIS_PORT", "6379"),
			lbc.Env("REDIS_PASSWORD", ""),
			lbc.Env("REDIS_DB", "0"),
		),
		lbc.WithUsers(lbc.Env("USERS", "")),
		lbc.WithMailJet(
			lbc.Env("MAILJET_KEY", ""),
			lbc.Env("MAILJET_SECRET", ""),
			lbc.Env("MAIL_FROM", ""),
		),
	)
}
