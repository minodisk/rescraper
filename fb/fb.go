package fb

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/minodisk/rescraper/errs"
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

	req, err := http.NewRequest(http.MethodPost, "https://graph.facebook.com", bytes.NewBufferString(values.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	cli := &http.Client{}
	res, err := cli.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode >= 300 {
		return errors.Wrapf(errs.NewHTTPError(res.StatusCode, res.Body), "fb: fail to scrape '%s'", u)
	}

	fmt.Printf("fb: sccess to scrape: %s\n", u)
	return nil
}
