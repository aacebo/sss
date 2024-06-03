package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/aacebo/sss/amqp"
	"github.com/aacebo/sss/models"
	"github.com/aacebo/sss/postgres"
	"github.com/aacebo/sss/repos"
	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
)

var visited = map[string]bool{}

func main() {
	pg := postgres.New()
	queue := amqp.New()
	domains := repos.NewDomain(pg)
	pages := repos.NewPage(pg)

	go func() {
		crawler := NewCrawler()
		crawler.Run("https://www.reddit.com")
		elapse := crawler.Elapse()
		log.Println(fmt.Sprintf(
			"%f:%f:%f",
			elapse.Hours(),
			elapse.Minutes(),
			elapse.Seconds(),
		))
	}()

	queue.Consume("pages", "upsert", func(d amqp091.Delivery) {
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
			domain = models.Domain{
				ID:        uuid.NewString(),
				Name:      domainName,
				Extension: ext,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			domain = domains.Create(domain)
		} else {
			domain = domains.Update(domain)
		}

		path := url.Hostname() + url.Path

		if path[:4] == "www." {
			path = path[4:]
		}

		log.Println(path)
		page, exists := pages.GetOne(path)

		var title *string

		if v, ok := data["title"].(string); ok {
			title = &v
		}

		if !exists {
			page = models.Page{
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
			}

			page = pages.Create(page)
		} else {
			page.Title = title
			page.Address = data["address"].(string)
			page.Size = int64(data["size"].(int))
			page.ElapseMs = data["elapse_ms"].(int64)
			page.LinkCount = data["link_count"].(int)
			page = pages.Update(page)
		}

		d.Ack(false)
	})
}
