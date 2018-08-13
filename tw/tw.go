package tw

import (
	"fmt"
	"net/url"

	"github.com/minodisk/rescraper/http"
)

type Client struct {
	authenticityToken string
	authToken         string
	csrfID            string
}

func NewClient(authenticityToken, authToken, csrfID string) (*Client, error) {
	if authenticityToken == "" || authToken == "" || csrfID == "" {
		return nil, fmt.Errorf("tw: empty token is not allowed")
	}
	return &Client{authenticityToken, authToken, csrfID}, nil
}

func (c *Client) Scrape(u string) error {
	fmt.Printf("tw: start scraping: %s\n", u)

	values := url.Values{}
	values.Set("authenticity_token", c.authenticityToken)
	values.Set("url", u)
	values.Set("platform", "Swift-12")

	cookies := http.NewCookies()
	cookies.Set("auth_token", c.authToken)
	cookies.Set("csrf_id", c.csrfID)

	_, err := http.PostFormWithCookies("https://cards-dev.twitter.com/validator/validate", values, cookies)
	if err != nil {
		return err
	}

	fmt.Printf("tw: success to scrape: %s\n", u)
	return nil
}
