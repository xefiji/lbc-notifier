package lbc

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"os"

	"github.com/rs/zerolog/log"
)

const test = false

type crawler struct {
	cfg  config
	repo repo
}

type repo interface {
	save(ad Ad) (Ad, error)
	get(id int64) (Ad, error)
}

func Crawl(opts ...Option) error {
	cfg := new(config)
	for _, opt := range opts {
		if err := opt(cfg); err != nil {
			return err
		}
	}

	crawler := crawler{
		cfg: *cfg,
		repo: newRepository(
			cfg.RedisHost,
			cfg.RedisPort,
			cfg.RedisPassword,
			cfg.RedisDB,
		),
	}

	result, err := crawler.fetch(test)
	if err != nil {
		return err
	}

	log.Info().Int("total", result.Total).Msg("results found")

	for _, ad := range result.Ads {
		log.Debug().Int64("id", ad.ListID).Msg("processing ad")

		if crawler.has(ad.ListID) {
			log.Debug().Int("id", int(ad.ListID)).Msg("ads exist")

			continue
		}

		if err := crawler.save(ad); err != nil {
			log.Error().Err(err).Interface("ad", ad).Msg("error while saving add")

			continue
		}

		log.Debug().Int("id", int(ad.ListID)).Msg("ads added")

		if err := crawler.notify(ad); err != nil {
			log.Error().Err(err).Interface("ad", ad).Msg("error while notifying for add")

			continue
		}
	}

	return nil
}

func (c *crawler) fetch(test bool) (Result, error) {
	wd, _ := os.Getwd()

	if test {
		response, err := os.ReadFile(fmt.Sprintf("%s/docs/sample_response.json", wd))
		if err != nil {
			return Result{}, err
		}

		var result Result
		if err := json.Unmarshal(response, &result); err != nil {
			return Result{}, err
		}

		return result, nil
	}

	request, err := os.ReadFile(fmt.Sprintf("%s/docs/request.json", wd))
	if err != nil {
		return Result{}, err
	}

	payload := strings.NewReader(string(request))
	req, _ := http.NewRequest(http.MethodPost, c.cfg.APIUrl, payload)

	req.Header.Add("content-type", "application/json")
	req.Header.Add("X-RapidAPI-Key", c.cfg.APIKey)
	req.Header.Add("X-RapidAPI-Host", c.cfg.APIHost)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return Result{}, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return Result{}, err
	}

	var result Result
	if err = json.Unmarshal(body, &result); err != nil {
		return Result{}, err
	}

	return result, nil
}

func (c *crawler) has(id int64) bool {
	ad, err := c.repo.get(id)
	if err != nil {
		log.Error().Err(err).Int64("id", id).Msg("error while checking if ad exists")

		return false
	}

	return ad.ListID != 0
}

func (c *crawler) save(ad Ad) error {
	_, err := c.repo.save(ad)

	return err
}

func (c *crawler) notify(ad Ad) error {
	// notify users.
	return nil
}

// func main() {

// 	request, err := os.ReadFile("request.json")
// 	if err != nil {
// 		log.Fatal().Err(err).Msg("error parsing request")
// 	}

// 	response, err := os.ReadFile("sample_response.json")
// 	if err != nil {
// 		log.Fatal().Err(err).Msg("error parsing response")
// 	}

// 	var result Result
// 	err = json.Unmarshal(response, &result)
// 	if err != nil {
// 		log.Fatal().Err(err).Msg("error unmarshalling response")
// 	}

// 	for _, ad := range result.Ads {
// 		fmt.Println(ad.PriceCents / 100)
// 	}

// 	fmt.Println(result.Ads[0].PriceCents / 100)

// 	os.Exit(1)

// 	url := os.Getenv("API_URL")
// 	payload := strings.NewReader(string(request))
// 	req, _ := http.NewRequest("POST", url, payload)

// 	req.Header.Add("content-type", "application/json")
// 	req.Header.Add("X-RapidAPI-Key", os.Getenv("RAPIDAPI_KEY"))
// 	req.Header.Add("X-RapidAPI-Host", os.Getenv("RAPIDAPI_HOST"))

// 	res, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		log.Error().Err(err).Msg("error sending request")
// 	}

// 	defer res.Body.Close()

// 	body, _ := io.ReadAll(res.Body)
// 	fmt.Println(res)
// 	fmt.Println(string(body))
// }
