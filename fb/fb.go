package fb

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Client struct {
	accessToken string
}

func NewClient(t string) (*Client, error) {
	fmt.Println(t)
	if t == "" {
		return nil, errors.New("empty access token")
	}
	return &Client{t}, nil
}

func (c *Client) Scrape(u string) error {
	fmt.Println("Scrape", u)
	res, err := http.PostForm("https://graph.facebook.com", url.Values{"id": {u}, "scrape": {"true"}, "access_token": {c.accessToken}})
	if err != nil {
		return err
	}
	if res.StatusCode >= 300 {
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		return errors.New(string(body))
	}
	return nil
}
