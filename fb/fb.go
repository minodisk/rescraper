package fb

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/minodisk/rescraper/errs"
	"github.com/pkg/errors"
)

type Client struct {
	accessToken string
}

func NewClient(token string) (*Client, error) {
	if token == "" {
		return nil, fmt.Errorf("fb: empty token is not allowed")
	}
	return &Client{token}, nil
}

func (c *Client) Scrape(u string) error {
	fmt.Printf("fb: start scraping: %s\n", u)

	res, err := http.PostForm("https://graph.facebook.com", url.Values{"id": {u}, "scrape": {"true"}, "access_token": {c.accessToken}})
	if err != nil {
		return err
	}
	if res.StatusCode >= 300 {
		return errors.Wrapf(errs.NewHTTPError(res.StatusCode, res.Body), "fb: fail to scrape '%s'", u)
	}

	fmt.Printf("fb: sccess to scrape: %s\n", u)
	return nil
}
