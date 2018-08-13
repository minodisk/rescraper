package http

import orig "net/http"

type Cookies []orig.Cookie

func NewCookies() *Cookies {
	return &Cookies{}
}

func (c *Cookies) Set(name, value string) {
	*c = append(*c, orig.Cookie{Name: name, Value: value})
}
