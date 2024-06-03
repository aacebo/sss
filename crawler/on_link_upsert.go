package main

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"log"
	"net/url"
	"time"

	"github.com/aacebo/sss/models"
	"github.com/aacebo/sss/repos"
	"github.com/rabbitmq/amqp091-go"
)

func onLinkUpsert(pg *sql.DB) func(d amqp091.Delivery) {
	pages := repos.NewPage(pg)
	links := repos.NewLink(pg)

	return func(d amqp091.Delivery) {
		data := map[string]string{}
		dec := gob.NewDecoder(bytes.NewBuffer(d.Body))

		if err := dec.Decode(&data); err != nil {
			log.Fatal(err)
		}

		fromUrl, err := url.Parse(data["from_url"])

		if err != nil {
			log.Fatal(err)
		}

		toUrl, err := url.Parse(data["to_url"])

		if err != nil {
			log.Fatal(err)
		}

		from, exists := pages.GetOne(UrlToString(fromUrl))

		if !exists {
			d.Nack(false, true)
			return
		}

		to, exists := pages.GetOne(UrlToString(toUrl))

		if !exists || from.ID == to.ID {
			d.Nack(false, true)
			return
		}

		if _, exists := links.GetOne(from.ID, to.ID); !exists {
			_, err = links.Create(models.Link{
				FromID:    from.ID,
				ToID:      to.ID,
				CreatedAt: time.Now(),
			})
		}

		if err != nil {
			d.Nack(false, true)
			return
		}

		d.Ack(false)
	}
}
