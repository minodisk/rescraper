# rescraper

A CLI tool to crawl **your** site and reset the OGP cache of valid links.

## Usage

```
rescraper http://example.com
```

## Installation

```
go get github.com/minodisk/rescraper
```

## Required Environment Variables

- `FB_ACCESS_TOKEN`
- `TW_AUTHENTICITY_TOKEN`
- `TW_AUTH_TOKEN`
- `TW_CSRF_ID`

### Facebook

1. Access to https://developers.facebook.com/tools/explorer/?classic=0
2. Select *Facebook App*.
3. Select *User Token* in *User of Page*.
4. Set `FB_ACCESS_TOKEN` *Access Token*.

![access_token](https://user-images.githubusercontent.com/514164/44014419-d7519a16-9f06-11e8-8f96-0ce9d0ddaf7a.png)

### Twitter

1. Access to https://cards-dev.twitter.com/validator
2. Open dev tools in browser.
3. Input *Card URL* and press *Preview card*.
4. Click *Name* `validate` in *Network* tab.
4. Open *Headers* tab of `validate`.
    - Set `TW_AUTHENTICITY_TOKEN` `authenticity_token` of *Form Data*.*
5. Open **Cookies* tab of `validate`.
    - Set `TW_AUTH_TOKEN` `auth_token` of Cookie
    - Set `TW_CSRF_ID` `csrf_id` of Cookie

![authenticity_token](https://user-images.githubusercontent.com/514164/44014388-b9a9b96c-9f06-11e8-9d52-8528b0abc653.png)
![auth_token-csrf_id_](https://user-images.githubusercontent.com/514164/44014282-fd8a695c-9f05-11e8-89aa-778d11908a16.png)
