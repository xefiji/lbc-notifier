package lbc

import "strconv"

type Option func(*config) error

type config struct {
	APIUrl        string
	APIKey        string
	APIHost       string
	RedisHost     string
	RedisPassword string
	RedisPort     int
	RedisDB       int
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
