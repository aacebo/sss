package main

import "net/url"

func UrlToString(url *url.URL) string {
	path := url.Hostname() + url.Path

	if len(path) > 4 && path[:4] == "www." {
		path = path[4:]
	}

	return path
}
