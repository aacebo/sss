package main

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/aacebo/sss/models"
	"github.com/aacebo/sss/repos"
	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
)

func onPageUpsert(pg *sql.DB) func(d amqp091.Delivery) {
	domains := repos.NewDomain(pg)
	pages := repos.NewPage(pg)

	return func(d amqp091.Delivery) {
		data := map[string]any{}
		dec := gob.NewDecoder(bytes.NewBuffer(d.Body))

		if err := dec.Decode(&data); err != nil {
			log.Fatal(err)
		}

		url, err := url.Parse(data["url"].(string))

		if err != nil {
			log.Fatal(err)
		}

		parts := strings.Split(url.Hostname(), ".")

		if parts[0] == "www" {
			parts = parts[1:]
		}

		domainName := parts[0]
		ext := parts[len(parts)-1]

		if len(parts) == 3 {
			domainName = parts[0] + "." + parts[1]
		}

		domain, exists := domains.GetOne(
			domainName,
			ext,
		)

		if !exists {
			domain, err = domains.Create(models.Domain{
				ID:        uuid.NewString(),
				Name:      domainName,
				Extension: ext,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			})
		} else {
			domain, err = domains.Update(domain)
		}

		if err != nil {
			d.Nack(false, true)
			return
		}

		path := UrlToString(url)
		log.Println(fmt.Sprintf(
			"%s %dms",
			path,
			data["elapse_ms"].(int64),
		))

		page, exists := pages.GetOne(path)

		var title *string

		if v, ok := data["title"].(string); ok && v != "" {
			title = &v
		}

		if !exists {
			page, err = pages.Create(models.Page{
				ID:        uuid.NewString(),
				DomainID:  domain.ID,
				Title:     title,
				Url:       path,
				Address:   data["address"].(string),
				Size:      int64(data["size"].(int)),
				ElapseMs:  data["elapse_ms"].(int64),
				LinkCount: data["link_count"].(int),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			})
		} else {
			page.Title = title
			page.Address = data["address"].(string)
			page.Size = int64(data["size"].(int))
			page.ElapseMs = data["elapse_ms"].(int64)
			page.LinkCount = data["link_count"].(int)
			page, err = pages.Update(page)
		}

		if err != nil {
			d.Nack(false, true)
			return
		}

		d.Ack(false)
	}
}
