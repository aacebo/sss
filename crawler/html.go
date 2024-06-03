package main

import (
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"
)

type Html struct {
	url     *url.URL
	address string
	content []byte
	status  int
	log     *log.Logger
	anchor  *regexp.Regexp
	title   *regexp.Regexp
}

func NewHtml(url *url.URL) *Html {
	return &Html{
		url:     url,
		content: nil,
		log:     log.New(os.Stdout, "html ", log.Ldate|log.Ltime|log.Lshortfile),
		anchor:  regexp.MustCompile("<a.*?href=\"(.*?)\""),
		title:   regexp.MustCompile("<title>(.*?)</title>"),
	}
}

func (self Html) Len() int {
	return len(self.content)
}

func (self Html) Url() *url.URL {
	return self.url
}

func (self Html) UrlString() string {
	return self.url.Hostname() + self.url.Path
}

func (self Html) Address() string {
	return self.address
}

func (self Html) Status() int {
	return self.status
}

func (self Html) Title() string {
	matches := self.title.FindAllStringSubmatch(string(self.content), -1)

	if len(matches) == 0 {
		return ""
	}

	return matches[len(matches)-1][1]
}

func (self Html) Urls() []string {
	urls := []string{}
	matches := self.anchor.FindAllStringSubmatch(string(self.content), -1)
	visited := map[string]bool{}

	for _, match := range matches {
		matchUrl, err := url.Parse(strings.TrimSpace(match[1]))

		if err != nil {
			self.log.Println(err)
			continue
		}

		var parsed string

		if matchUrl.IsAbs() {
			parsed = matchUrl.String()
		} else {
			parsed = self.url.Scheme + "://" + self.url.Host + matchUrl.String()
		}

		if _, ok := visited[parsed]; !ok {
			urls = append(urls, parsed)
			visited[parsed] = true
		}
	}

	return urls
}
