package lbc

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type repoConfig interface {
	disable() error
	enable() error
}

func Listen(opts ...Option) error {
	cfg := new(config)
	for _, opt := range opts {
		if err := opt(cfg); err != nil {
			log.Error().Err(err).Msg("invalid configuration")

			return err
		}
	}

	repo := newRepository(
		cfg.RedisHost,
		cfg.RedisPort,
		cfg.RedisPassword,
		cfg.RedisDB,
	)

	r := mux.NewRouter().StrictSlash(false)
	r.HandleFunc("/disable", func(w http.ResponseWriter, r *http.Request) {
		toggleHandler(w, repo, false)
	}).Methods("GET")

	r.HandleFunc("/enable", func(w http.ResponseWriter, r *http.Request) {
		toggleHandler(w, repo, true)
	}).Methods("GET")

	return serve(r, cfg.ServerPort)
}

func serve(router http.Handler, port string) error {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		Handler:           router,
		ReadHeaderTimeout: 3000,
	}

	sink := make(chan error, 1)

	go func() {
		defer close(sink)
		sink <- srv.ListenAndServe()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	select {
	case <-quit:
		return shutdown(srv, "quit signaled")
	case err := <-sink:
		return err
	}
}

func shutdown(srv *http.Server, from string) error {
	ctx, cancel := context.WithTimeout(context.Background(), (20 * time.Second))
	defer cancel()
	log.Warn().Msg(fmt.Sprintf("shutting down from %s", from))

	return srv.Shutdown(ctx)
}

func toggleHandler(w http.ResponseWriter, repo repoConfig, enabled bool) {
	var message string

	w.Header().Set("Content-Type", "application/json")

	var action func() error
	if enabled {
		action = repo.enable
	} else {
		action = repo.disable
	}

	if err := action(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		message = `{"message": "error occurred"}`

		log.Error().Err(err).Msg("error while toggling")
	} else {
		w.WriteHeader(http.StatusOK)

		message = `{"message": "ok"}`
	}

	if _, err := w.Write([]byte(message)); err != nil {
		log.Error().Err(err).Msg("error while sending response")
	}
}
