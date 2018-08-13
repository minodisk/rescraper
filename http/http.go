package http

import (
	"bytes"
	"io"
	"io/ioutil"
	_http "net/http"
	_url "net/url"

	"github.com/minodisk/rescraper/errs"
	"github.com/pkg/errors"
)

func PostFormWithCookies(url string, data _url.Values, cookies *Cookies) (io.ReadCloser, error) {
	req, err := _http.NewRequest(_http.MethodPost, url, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, errors.Wrap(err, "can not create request")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	if cookies != nil {
		for _, cookie := range *cookies {
			req.Header.Add("Cookie", cookie.String())
		}
	}

	return Do(req)
}

func Get(url string) (io.ReadCloser, error) {
	req, err := _http.NewRequest(_http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "can not create request")
	}
	return Do(req)
}

func Head(url string) (io.ReadCloser, error) {
	req, err := _http.NewRequest(_http.MethodHead, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "can not create request")
	}
	return Do(req)
}

func Do(req *_http.Request) (io.ReadCloser, error) {
	cli := &_http.Client{}
	cli.CheckRedirect = ignoreRedirect
	res, err := cli.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "can not send HTTP request")
	}

	if res.StatusCode >= 300 {
		buf, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, errors.Wrap(err, "can not read response body")
		}
		return nil, errs.NewHTTPError(res.StatusCode, string(buf))
	}

	return res.Body, nil
}

func ignoreRedirect(req *_http.Request, via []*_http.Request) error {
	return _http.ErrUseLastResponse
}
