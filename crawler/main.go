package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/aacebo/sss/amqp"
	"github.com/aacebo/sss/async"
	"github.com/aacebo/sss/postgres"
)

var visited = async.NewMap[string, bool]()

func main() {
	pg := postgres.New()
	queue := amqp.New()
	roots := []string{
		"https://www.reddit.com/r/popular/",
		"https://stackoverflow.com/questions?tab=Votes",
		"https://www.yahoo.com/",
		"https://en.wikipedia.org/",
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		wg := sync.WaitGroup{}

		for _, root := range roots {
			wg.Add(1)

			go func() {
				defer wg.Done()

				crawler := NewCrawler()
				crawler.Run(root)
				elapse := crawler.Elapse()
				log.Println(fmt.Sprintf(
					"%f:%f:%f",
					elapse.Hours(),
					elapse.Minutes(),
					elapse.Seconds(),
				))
			}()
		}

		wg.Wait()
	}()

	go func() {
		defer wg.Done()
		queue.Consume("pages", "upsert", onPageUpsert(pg))
	}()

	go func() {
		defer wg.Done()
		queue.Consume("links", "upsert", onLinkUpsert(pg))
	}()

	wg.Wait()
}
