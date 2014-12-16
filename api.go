package vk

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var (
	// Version of VK API
	Version = "5.27"
	// APIURL is a base to make API calls
	APIURL       = "https://api.vk.com/method/"
	reqTokURL, _ = url.Parse("https://oauth.vk.com/authorize")
	accTokURL, _ = url.Parse("https://oauth.vk.com/access_token")
	// HTTPS defines if use https instead of http. 1 - use https. 0 - use http
	HTTPS = 1
	Debug = false
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
	UserID      int
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

func (s *Session) CallAPI(method string, params url.Values, out interface{}) error {
	endpoint := s.getAPIURL(method)
	query := endpoint.Query()
	for k, v := range params {
		if len(v) > 0 {
			query.Set(k, v[0])
		}
	}
	endpoint.RawQuery = query.Encode()

	var (
		err      error
		resp     *http.Response
		response struct {
			Error struct {
				Code int    `json:"error_code"`
				Msg  string `json:"error_msg"`
			} `json:"error"`
			Response json.RawMessage `json:"response"`
		}
	)
	//response.Response = out

	if Debug {
		fmt.Printf("vk call: %s\n", endpoint.String())
	}
	if resp, err = http.Get(endpoint.String()); err != nil {
		return err
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}
	if response.Error.Code != 0 {
		return errors.New(response.Error.Msg)
	}
	if Debug {
		fmt.Printf("vk api resp: %s\n", string(response.Response))
	}
	if err = json.Unmarshal(response.Response, out); err != nil {
		return err
	}
	return nil
}

type ApiList struct {
	Count int         `json:"count"`
	Items interface{} `json:"items"`
}
