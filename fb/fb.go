package fb

import (
	"net/http"
	"net/url"
)

type Client struct {
	accessToken string
}

func NewClient(t string) *Client {
	return &Client{t}
}

func (c *Client) Scrape(u string) error {
	form := url.Values{}
	form.Add("id", u)
	form.Add("scrape", "true")
	form.Add("access_token", c.accessToken)
	return http.Post("https://graph.facebook.com", "", form)
}
