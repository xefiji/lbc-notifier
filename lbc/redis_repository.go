package lbc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v9"
)

const adsKey = "ads"

var ErrAdNotFound = errors.New("add was not found")

type repository struct {
	db *redis.Client
}

func newRepository(host string, port int, password string, db int) repository {
	return repository{
		db: redis.NewClient(
			&redis.Options{
				Addr:     fmt.Sprintf("%s:%d", host, port),
				Password: password,
				DB:       db,
			},
		),
	}
}

func (r repository) save(ad Ad) (Ad, error) {
	json, err := json.Marshal(ad)
	if err != nil {
		return ad, err
	}

	ctx := context.Background()
	if err = r.db.Set(ctx, fmt.Sprintf("%s:%d", adsKey, ad.ListID), json, 0).Err(); err != nil {
		return ad, err
	}

	return ad, nil
}

func (r repository) get(id int64) (Ad, error) {
	ctx := context.Background()

	val, err := r.db.Get(ctx, fmt.Sprintf("%s:%d", adsKey, id)).Result()
	if errors.Is(err, redis.Nil) {
		return Ad{}, ErrAdNotFound
	}

	if err != nil {
		return Ad{}, err
	}

	var ad Ad

	err = json.Unmarshal([]byte(val), &ad)
	if err != nil {
		return Ad{}, err
	}

	return ad, nil
}
