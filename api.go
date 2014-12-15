package vk

import (
	"net/url"
	"strconv"
	"time"
)

var (
	// Version of VK API
	Version = "5.12"
	// APIURL is a base to make API calls
	APIURL = "https://api.vk.com/method/"
	// HTTPS defines if use https instead of http. 1 - use https. 0 - use http
	HTTPS = 1
)

// API holds data to use for communication
type API struct {
	AppID           string
	Secret          string
	Scope           []Scope
	AccessToken     string
	Expiry          time.Time
	UserID          string
	UserEmail       string
	callbackURL     *url.URL
	requestTokenURL *url.URL
	accessTokenURL  *url.URL
}

// NewAPI creates instance of API
func NewAPI(appID, secret string, scope []Scope, callback string) *API {
	var err error
	var callbackURL *url.URL

	if appID == "" {
		return nil
	}
	if secret == "" {
		return nil
	}
	if callbackURL, err = url.Parse(callback); err != nil {
		return nil
	}
	reqTokURL, _ := url.Parse("https://oauth.vk.com/authorize")
	accTokURL, _ := url.Parse("https://oauth.vk.com/access_token")

	return &API{
		AppID:           appID,
		Secret:          secret,
		Scope:           scope,
		callbackURL:     callbackURL,
		requestTokenURL: reqTokURL,
		accessTokenURL:  accTokURL,
	}
}

func (api *API) NewSession(tok string) *Session {
	if tok == "" {
		tok = api.AccessToken
	}
	return &Session{AccessToken: tok}
}

type Session struct {
	AccessToken string
	UserID      string
	UserEmail   string
}

// getAPIURL prepares URL instance with defined method
func (s *Session) getAPIURL(method string) *url.URL {
	q := url.Values{
		"v":            {Version},
		"https":        {strconv.Itoa(HTTPS)},
		"access_token": {s.AccessToken},
	}.Encode()
	apiURL, _ := url.Parse(APIURL + method + "?" + q)
	return apiURL
}
