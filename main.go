package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/minodisk/rescraper/fb"
)

var (
	cli *fb.Client
	m   = new(sync.Mutex)
)

type ListedError struct {
	url string
}

func NewListedError(u string) *ListedError {
	return &ListedError{u}
}

func (e *ListedError) Error() string {
	return fmt.Sprintf("%s: %s", "already listed", e.url)
}

type BadStatusError struct {
	statusCode int
}

func NewBadStatusError(s int) *BadStatusError {
	return &BadStatusError{s}
}

func (e *BadStatusError) Error() string {
	return fmt.Sprintf("bad staus: %d", e.statusCode)
}

type Errors []error

func NewErrors() Errors {
	return Errors{}
}

func (e Errors) Append(err error) Errors {
	return append(e, err)
}

func (e Errors) Error() string {
	msgs := make([]string, len(e))
	for i, err := range e {
		msgs[i] = err.Error()
	}
	return strings.Join(msgs, "\n")
}

var cache = []string{}

func main() {
	var err error
	cli, err = fb.NewClient(os.Getenv("FB_ACCESS_TOKEN"))
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(1)
		return
	}

	flag.Parse()
	root := flag.Arg(0)
	if err := process(root); err != nil {
		if _, e := os.Stderr.WriteString(err.Error()); e != nil {
			panic(e)
		}
		os.Exit(1)
	}
}

func process(u string) error {
	normalizedURL, err := normalizeURL(u)
	if err != nil {
		return err
	}
	if err := parallel([]string{normalizedURL}); err != nil {
		return err
	}

	fmt.Println(len(cache), cache)

	return nil
}

var id = 0

func normalizeURL(u string) (string, error) {
	obj, err := url.Parse(u)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s://%s%s", obj.Scheme, obj.Host, obj.RequestURI()), nil
}

func filterURLs(urls []string) ([]string, error) {
	unlistedURLs := []string{}
outer:
	for _, u := range urls {
		normalizedURL, err := normalizeURL(u)
		if err != nil {
			return nil, err
		}
		for _, listedURL := range cache {
			if normalizedURL == listedURL {
				// already listed
				continue outer
			}
		}
		unlistedURLs = append(unlistedURLs, normalizedURL)
	}
	return unlistedURLs, nil
}

func parallel(urls []string) error {
	for _, u := range urls {
		cache = append(cache, u)
	}

	count := 0
	errChan := make(chan error, 1)
	doneChan := make(chan struct{}, 1)

	for _, u := range urls {
		count++
		go func(errChan chan error, u string) {
			defer func() {
				doneChan <- struct{}{}
			}()

			m.Lock()
			go func(u string) {
				defer m.Unlock()
				time.Sleep(2 * time.Second)
				if err := cli.Scrape(u); err != nil {
					errChan <- err
				}
			}(u)

			us, err := scrape(u)
			if err != nil {
				errChan <- err
				return
			}

			filteredURLs, err := filterURLs(us)
			if err != nil {
				errChan <- err
				return
			}
			if len(filteredURLs) == 0 {
				return
			}
			if err := parallel(filteredURLs); err != nil {
				errChan <- err
			}
		}(errChan, u)
	}

	errs := NewErrors()
loop:
	for {
		select {
		case err := <-errChan:
			if err != nil {
				errs = append(errs, err)
			}
		case <-doneChan:
			count--
			if count == 0 {
				break loop
			}
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func scrape(u string) ([]string, error) {
	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	if res.StatusCode >= 300 {
		return nil, NewBadStatusError(res.StatusCode)
	}
	defer func() {
		if e := res.Body.Close(); e != nil {
			panic(e)
		}
	}()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	hrefs := doc.Find("a").FilterFunction(func(i int, s *goquery.Selection) bool {
		_, exists := s.Attr("href")
		if !exists {
			return false
		}
		return true
	}).Map(func(i int, s *goquery.Selection) string {
		href, _ := s.Attr("href")
		return href
	})

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
