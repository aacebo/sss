package main

import (
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"os"
	"regexp"
	"time"

	"github.com/aacebo/sss/amqp"
)

type Crawler struct {
	http      *http.Client
	amqp      *amqp.Client
	log       *log.Logger
	anchor    *regexp.Regexp
	title     *regexp.Regexp
	startedAt *time.Time
	endedAt   *time.Time
}

func NewCrawler() *Crawler {
	return &Crawler{
		http:   &http.Client{Timeout: 5 * time.Second},
		amqp:   amqp.New(),
		log:    log.New(os.Stdout, "crawler ", log.Ldate|log.Ltime|log.Lshortfile),
		anchor: regexp.MustCompile("<a.*?href=\"(.*?)\""),
		title:  regexp.MustCompile("<title>(.*?)</title>"),
	}
}

func (self *Crawler) Elapse() time.Duration {
	return self.endedAt.Sub(*self.startedAt)
}

func (self *Crawler) Run(to string) {
	now := time.Now()
	self.startedAt = &now
	self.endedAt = nil
	defer self.done()
	self.visit(to, to, 0)
}

func (self *Crawler) done() {
	now := time.Now()
	self.endedAt = &now
}

func (self *Crawler) visit(from string, to string, depth int64) {
	url, err := url.Parse(to)

	if err != nil {
		log.Println(err)
		return
	}

	path := fmt.Sprintf(
		"%s => %s%s",
		from,
		url.Hostname(),
		url.Path,
	)

	if visited.Has(path) {
		return
	}

	visited.Set(path, true)
	startedAt := time.Now()
	page, err := self.get(url)

	if err != nil || page.Status() != http.StatusOK {
		return
	}

	urls := page.Urls()
	self.amqp.Publish("pages", "upsert", map[string]any{
		"from_url":   from,
		"title":      page.Title(),
		"url":        to,
		"address":    page.Address(),
		"size":       page.Len(),
		"elapse_ms":  time.Now().UnixMilli() - startedAt.UnixMilli(),
		"link_count": len(urls),
	})

	for _, url := range urls {
		self.visit(to, url, depth+1)
	}
}

func (self *Crawler) get(url *url.URL) (*Html, error) {
	req, err := http.NewRequest(http.MethodGet, url.String(), nil)

	if err != nil {
		return nil, err
	}

	page := NewHtml(url)
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), &httptrace.ClientTrace{
		GotConn: func(info httptrace.GotConnInfo) {
			page.address = info.Conn.RemoteAddr().String()
		},
	}))

	res, err := self.http.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	content, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	page.content = []byte(html.UnescapeString(string(content)))
	page.status = res.StatusCode
	return page, nil
}
