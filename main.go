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

var (
	fbClient   *fb.Client
	twClient   *tw.Client
	fbMutex    = new(sync.Mutex)
	twMutex    = new(sync.Mutex)
	cacheMutex = new(sync.RWMutex)
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

	normalizedURL, err := normalizeURL(u)
	if err != nil {
		return err
	}

	wg := &sync.WaitGroup{}

	if err := process(wg, []string{normalizedURL}); err != nil {
		return err
	}

	wg.Wait()

	fmt.Printf("%s\nfinally scraped: %d URLs\n%s\n%s\n", strings.Repeat("-", 50), len(cache), strings.Repeat("-", 50), strings.Join(cache, "\n"))

	return nil
}

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
		cacheMutex.RLock()
		for _, listedURL := range cache {
			if normalizedURL == listedURL {
				// already listed
				cacheMutex.RUnlock()
				continue outer
			}
		}
		cacheMutex.RUnlock()
		unlistedURLs = append(unlistedURLs, normalizedURL)
	}
	return unlistedURLs, nil
}

func process(wg *sync.WaitGroup, urls []string) error {
	cacheMutex.Lock()
	for _, u := range urls {
		cache = append(cache, u)
	}
	cacheMutex.Unlock()

	for _, u := range urls {
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
		time.Sleep(18 * time.Second)
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
		time.Sleep(2 * time.Second)
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

	filteredURLs, err := filterURLs(us)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(filteredURLs) == 0 {
		return
	}
	if err := process(wg, filteredURLs); err != nil {
		fmt.Println(err)
	}
}
