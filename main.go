package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/minodisk/rescraper/fb"
	"github.com/minodisk/rescraper/scraper"
	"github.com/minodisk/rescraper/tw"
)

const (
	fbWait = 1 * time.Second
	twWait = 1 * time.Second
)

var (
	fbClient   *fb.Client
	twClient   *tw.Client
	fbMutex    = new(sync.Mutex)
	twMutex    = new(sync.Mutex)
	cacheMutex = new(sync.Mutex)
	cache      = []string{}
)

func main() {
	if err := _main(); err != nil {
		if _, e := os.Stderr.WriteString(fmt.Sprintf("%s\n", err.Error())); e != nil {
			panic(e)
		}
		os.Exit(1)
		return
	}
}

func _main() error {
	var err error

	fbClient, err = fb.NewClient(os.Getenv("FB_ACCESS_TOKEN"))
	if err != nil {
		return err
	}

	twClient, err = tw.NewClient(
		os.Getenv("TW_AUTHENTICITY_TOKEN"),
		os.Getenv("TW_AUTH_TOKEN"),
		os.Getenv("TW_CSRF_ID"),
	)
	if err != nil {
		return err
	}

	flag.Parse()
	u := flag.Arg(0)

	wg := &sync.WaitGroup{}

	if err := process(wg, []string{u}); err != nil {
		return err
	}

	wg.Wait()

	fmt.Printf("%s\nfinally scraped: %d URLs\n%s\n%s\n", strings.Repeat("-", 50), len(cache), strings.Repeat("-", 50), strings.Join(cache, "\n"))

	return nil
}

func process(wg *sync.WaitGroup, urls []string) error {
	unlistedURLs := []string{}
	cacheMutex.Lock()
outer:
	for _, u := range urls {
		normalizedURL, err := normalizeURL(u)
		if err != nil {
			return err
		}

		for _, listedURL := range cache {
			if normalizedURL == listedURL {
				// already listed
				continue outer
			}
		}

		unlistedURLs = append(unlistedURLs, normalizedURL)
		cache = append(cache, normalizedURL)
	}
	cacheMutex.Unlock()

	for _, u := range unlistedURLs {
		fmt.Println("->", u)

		wg.Add(1)
		go rescrapeFB(wg, u)

		wg.Add(1)
		go rescrapeTW(wg, u)

		wg.Add(1)
		go scrape(wg, u)
	}

	return nil
}

func rescrapeFB(wg *sync.WaitGroup, u string) {
	defer func() {
		time.Sleep(fbWait)
		fbMutex.Unlock()
		wg.Done()
	}()
	fbMutex.Lock()
	if err := fbClient.Scrape(u); err != nil {
		fmt.Println(err)
	}
}

func rescrapeTW(wg *sync.WaitGroup, u string) {
	defer func() {
		time.Sleep(twWait)
		twMutex.Unlock()
		wg.Done()
	}()
	twMutex.Lock()
	if err := twClient.Scrape(u); err != nil {
		fmt.Println(err)
	}
}

func scrape(wg *sync.WaitGroup, u string) {
	defer func() {
		wg.Done()
	}()

	us, err := scraper.Scrape(u)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := process(wg, us); err != nil {
		fmt.Println(err)
	}
}

func normalizeURL(u string) (string, error) {
	obj, err := url.Parse(u)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s://%s%s", obj.Scheme, obj.Host, obj.RequestURI()), nil
}
