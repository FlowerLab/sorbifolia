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
	/* #nosec G107 */
	resp, err := http.Get(u)
	if err != nil {
		log.Panicln(err)
	}
	defer func() { _ = resp.Body.Close() }()

	var sm *sitemap.Sitemap
	if sm, err = sitemap.Parse(resp.Body); err != nil {
		log.Panicln(err)
	}

	return sm
}
