package scraper

import (
	"io"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/minodisk/rescraper/http"
	"github.com/pkg/errors"
)

func Scrape(u string) ([]string, error) {
	body, err := http.Get(u)
	if err != nil {
		return nil, errors.Wrapf(err, "can not fetch '%s'", u)
	}

	hrefs, err := parse(body)
	if err != nil {
		return nil, err
	}

	return filter(u, hrefs)
}

func parse(body io.ReadCloser) ([]string, error) {
	defer func() {
		if e := body.Close(); e != nil {
			panic(e)
		}
	}()

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}
	return doc.Find("a").FilterFunction(func(i int, s *goquery.Selection) bool {
		href, exists := s.Attr("href")
		if href == "" || !exists {
			return false
		}
		return true
	}).Map(func(i int, s *goquery.Selection) string {
		href, _ := s.Attr("href")
		return href
	}), nil
}

func filter(u string, hrefs []string) ([]string, error) {
	base, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	validHrefs := []string{}
	for _, u := range hrefs {
		ref, err := url.Parse(u)
		if err != nil {
			return nil, err
		}
		to := base.ResolveReference(ref)
		if base.Scheme != to.Scheme || base.Host != to.Host {
			// different origin
			continue
		}
		validHrefs = append(validHrefs, to.String())
	}

	return validHrefs, nil
}
