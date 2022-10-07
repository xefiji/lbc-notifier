package lbc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v9"
)

const (
	adsKey     = "ads"
	configKey  = "config"
	enabledKey = "enabled"
)

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

func (r repository) disable() error {
	ctx := context.Background()
	if err := r.db.Set(ctx, fmt.Sprintf("%s:%s", configKey, enabledKey), false, 0).Err(); err != nil {
		return err
	}

	return nil
}

func (r repository) enable() error {
	ctx := context.Background()
	if err := r.db.Set(ctx, fmt.Sprintf("%s:%s", configKey, enabledKey), true, 0).Err(); err != nil {
		return err
	}

	return nil
}

func (r repository) enabled() (bool, error) {
	ctx := context.Background()

	val, err := r.db.Get(ctx, fmt.Sprintf("%s:%s", configKey, enabledKey)).Result()
	if errors.Is(err, redis.Nil) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	boolValue, err := strconv.ParseBool(val)
	if err != nil {
		return false, err
	}

	return boolValue, nil
}
