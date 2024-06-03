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
	"strings"
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
	self.visit(to, 0)
}

func (self *Crawler) done() {
	now := time.Now()
	self.endedAt = &now
}

func (self *Crawler) visit(to string, depth int64) {
	url, err := url.Parse(to)

	if err != nil {
		panic(err)
	}

	path := fmt.Sprintf(
		"%s%s",
		url.Hostname(),
		url.Path,
	)

	if v, _ := visited[path]; v {
		return
	}

	visited[path] = true
	startedAt := time.Now()
	res, err := self.get(to)

	if err != nil || res.StatusCode != http.StatusOK {
		return
	}

	defer res.Body.Close()
	content, err := io.ReadAll(res.Body)

	if err != nil {
		self.log.Println(err)
		return
	}

	raw := html.UnescapeString(string(content))
	title := self.parseTitle(raw)
	urls, err := self.parseUrls(url, raw)

	if err != nil {
		panic(err)
	}

	self.amqp.Publish("pages", "upsert", map[string]any{
		"title":      title,
		"url":        to,
		"address":    res.Header.Get("X-RemoteAddress"),
		"size":       len(content),
		"elapse_ms":  time.Now().UnixMilli() - startedAt.UnixMilli(),
		"link_count": len(urls),
	})

	for _, url := range urls {
		defer self.visit(url, depth+1)
	}
}

func (self *Crawler) parseTitle(html string) string {
	matches := self.title.FindAllStringSubmatch(html, -1)

	if len(matches) == 0 {
		return ""
	}

	return matches[len(matches)-1][1]
}

func (self *Crawler) parseUrls(url *url.URL, html string) ([]string, error) {
	urls := []string{}
	matches := self.anchor.FindAllStringSubmatch(html, -1)

	for _, match := range matches {
		matchUrl, err := url.Parse(strings.TrimSpace(match[1]))

		if err != nil {
			self.log.Println(err)
			continue
		}

		if matchUrl.IsAbs() {
			urls = append(urls, matchUrl.String())
		} else {
			urls = append(urls, url.Scheme+"://"+url.Host+matchUrl.String())
		}
	}

	return urls, nil
}

func (self *Crawler) get(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return nil, err
	}

	address := ""

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), &httptrace.ClientTrace{
		GotConn: func(info httptrace.GotConnInfo) {
			address = info.Conn.RemoteAddr().String()
		},
	}))

	res, err := self.http.Do(req)

	if err != nil {
		return nil, err
	}

	res.Header.Set("X-RemoteAddress", address)
	return res, nil
}
