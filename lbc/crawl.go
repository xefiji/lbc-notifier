package lbc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"text/template"

	"os"

	"github.com/mailjet/mailjet-apiv3-go"
	"github.com/rs/zerolog/log"
)

const test = true

type crawler struct {
	cfg  config
	repo repo
}

type email struct {
	Title string
	Ads   []Ad
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

	log.Info().Int("total", result.Total).Msg("nb results found")

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
			log.Error().Err(err).Int64("id", ad.ListID).Msg("error while notifying for add")

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
	mailjetClient := mailjet.NewMailjetClient(
		c.cfg.MailJetKey,
		c.cfg.MailJetSecret,
	)

	title := fmt.Sprintf("Nouvelle annonce à %s", ad.Location.CityLabel)
	tmpl := template.Must(template.ParseFiles("docs/email.html"))

	data := email{
		Title: title,
		Ads:   []Ad{ad},
	}

	var html bytes.Buffer
	if err := tmpl.Execute(&html, data); err != nil {
		return err
	}

	messagesInfo := make([]mailjet.InfoMessagesV31, 0)

	for _, email := range c.cfg.Users {
		message := mailjet.InfoMessagesV31{
			From: &mailjet.RecipientV31{
				Email: c.cfg.MailFrom,
				Name:  "LBC Crawl",
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: email,
				},
			},
			Subject:  title,
			HTMLPart: html.String(),
		}

		messagesInfo = append(messagesInfo, message)
	}

	messages := mailjet.MessagesV31{Info: messagesInfo}

	_, err := mailjetClient.SendMailV31(&messages)
	if err != nil {
		return err
	}

	return nil
}
