package fb

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/minodisk/rescraper/http"
	"github.com/pkg/errors"
)

type Client struct {
	accessTokens []string
	index        int
}

func NewClient(token string) (*Client, error) {
	if token == "" {
		return nil, fmt.Errorf("fb: empty token is not allowed")
	}
	return &Client{strings.Split(token, ","), 0}, nil
}

func (c *Client) Scrape(u string) error {
	fmt.Printf("fb: start scraping: %s\n", u)

	values := url.Values{}
	values.Set("id", u)
	values.Set("scrape", "true")
	values.Set("access_token", c.accessTokens[c.index])
	c.index = (c.index + 1) % len(c.accessTokens)

	_, err := http.PostFormWithCookies("https://graph.facebook.com", values, nil)
	if err != nil {
		return errors.Wrap(err, "can not post form")
	}

	fmt.Printf("fb: sccess to scrape: %s\n", u)
	return nil
}
