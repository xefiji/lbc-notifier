package lbc

import (
	"os"
	"strconv"
	"strings"
)

type Option func(*config) error

type config struct {
	APIUrl        string
	APIKey        string
	APIHost       string
	RedisHost     string
	RedisPassword string
	RedisPort     int
	RedisDB       int
	Users         []string
	MailJetKey    string
	MailJetSecret string
	MailFrom      string
	ShouldExecute bool
	ServerPort    string
}

func WithRedis(host, port, password, db string) Option {
	portInt, _ := strconv.Atoi(port)
	dbInt, _ := strconv.Atoi(db)

	return func(cfg *config) error {
		cfg.RedisHost = host
		cfg.RedisPassword = password
		cfg.RedisPort = portInt
		cfg.RedisDB = dbInt

		return nil
	}
}

func WithRapidAPI(url, host, key string) Option {
	return func(cfg *config) error {
		cfg.APIHost = host
		cfg.APIKey = key
		cfg.APIUrl = url

		return nil
	}
}

func WithUsers(users string) Option {
	return func(cfg *config) error {
		cfg.Users = strings.Split(users, ",")

		return nil
	}
}

func WithMailJet(key, secret, from string) Option {
	return func(cfg *config) error {
		cfg.MailJetKey = key
		cfg.MailJetSecret = secret
		cfg.MailFrom = from

		return nil
	}
}

func WithDryRun(dryRun string) Option {
	return func(cfg *config) error {
		boolValue, err := strconv.ParseBool(dryRun)
		if err != nil {
			return err
		}

		cfg.ShouldExecute = !boolValue

		return nil
	}
}

func WithServer(port string) Option {
	return func(cfg *config) error {
		cfg.ServerPort = port

		return nil
	}
}

func Env(name, fallback string) string {
	if val, ok := os.LookupEnv(name); ok {
		return val
	}

	return fallback
}
