package tw

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"

	"github.com/minodisk/rescraper/errs"
	"github.com/pkg/errors"
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

	req, err := http.NewRequest(http.MethodPost, "https://cards-dev.twitter.com/validator/validate", bytes.NewBufferString(values.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("Cookie", fmt.Sprintf("auth_token=%s", c.authToken))
	req.Header.Add("Cookie", fmt.Sprintf("csrf_id=%s", c.csrfID))

	cli := &http.Client{}
	res, err := cli.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode >= 300 {
		return errors.Wrapf(errs.NewHTTPError(res.StatusCode, res.Body), "tw: fail to scrape '%s'", u)
	}

	fmt.Printf("tw: success to scrape: %s\n", u)
	return nil
}
