package main

import (
	"log"
	"net/http"

	"go.x2ox.com/sorbifolia/sitemap"
)

func importURL(u string) *sitemap.Sitemap {
	if u == "" {
		return nil
	}
	resp, err := http.Get(u)
	if err != nil {
		log.Panicln(err)
	}
	defer resp.Body.Close()

	sm, err := sitemap.Parse(resp.Body)
	if err != nil {
		log.Panicln(err)
	}

	return sm
}
