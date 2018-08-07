package scraper

import (
	"io"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/minodisk/rescraper/errs"
	"github.com/pkg/errors"
)

func Scrape(u string) ([]string, error) {
	body, err := fetch(u)
	if err != nil {
		return nil, err
	}

	hrefs, err := parse(body)
	if err != nil {
		return nil, err
	}

	return filter(u, hrefs)
}

func fetch(u string) (io.ReadCloser, error) {
	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	if res.StatusCode >= 300 {
		return nil, errors.Wrapf(errs.NewHTTPError(res.StatusCode, res.Body), "fail to fetch: %s", u)
	}
	return res.Body, nil
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
