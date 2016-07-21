package vk

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/codeship/go-retro"
)

const (
	vkAuthorize   = "https://oauth.vk.com/authorize"
	vkAccessToken = "https://oauth.vk.com/access_token"
)

var (
	Debug = false
	// Version of VK API
	Version = "5.53"
	// APIURL is a base to make API calls
	APIURL = "https://api.vk.com/method/"
	// HTTPS defines if use https instead of http. 1 - use https. 0 - use http
	HTTPS                = 1
	errTooManyTooManyReq = retro.NewStaticRetryableError(
		errors.New("Too too many requests: even 'retro' used up"), 35, 1)
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

type resolveCaptcha struct {
	id, key string
	sync.RWMutex
}

func (r *resolveCaptcha) set(id, key string) {
	r.Lock()
	r.id = id
	r.key = key
	r.Unlock()
}

func (r *resolveCaptcha) clear() {
	r.Lock()
	r.id = ""
	r.key = ""
	r.Unlock()
}

func (r *resolveCaptcha) get(v url.Values) {
	r.RLock()
	if r.id != "" {
		v.Set("captcha_sid", r.id)
		v.Set("captcha_key", r.key)
	}
	r.RUnlock()
}

var captcha = new(resolveCaptcha)

func SetCaptcha(id, key string) {
	captcha.set(id, key)
}

func ClearCaptcha() {
	captcha.clear()
}

type delayer struct {
	d    time.Duration
	next time.Time
	sync.Mutex
}

func NewDelayer(d time.Duration) *delayer {
	return &delayer{d: d, next: time.Now()}
}

func (d *delayer) Wait() {
	d.Lock()
	time.Sleep(d.next.Sub(time.Now()))
	d.next = time.Now().Add(d.d)
	d.Unlock()
}

var gdelay = NewDelayer(time.Millisecond * 400) // 3 APIs call per seconds

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
	ru, err := url.Parse(vkAuthorize)
	if err != nil {
		return nil
	}
	atu, err := url.Parse(vkAccessToken)
	if err != nil {
		return nil
	}
	return &API{
		AppID:           appID,
		Secret:          secret,
		Scope:           scope,
		callbackURL:     callbackURL,
		requestTokenURL: ru,
		accessTokenURL:  atu,
	}
}

func (api *API) NewSession(tok string) *Session {
	if tok == "" {
		tok = api.AccessToken
	}
	return &Session{AccessToken: tok}
}

func NewSession(tok string) *Session {
	if tok == "" {
		return nil
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

// save for multi-goroutines
func (s *Session) CallAPI(method string, params url.Values, out interface{}) error {
	endpoint := s.getAPIURL(method)
	query := endpoint.Query()
	captcha.get(params)
	for k, v := range params {
		if len(v) > 0 {
			query.Set(k, v[0])
		}
	}
	endpoint.RawQuery = query.Encode()
	err := retro.DoWithRetry(func() error {
		var (
			err      error
			resp     *http.Response
			response struct {
				Err      *Error          `json:"error"`
				Response json.RawMessage `json:"response"`
			}
		)
		if Debug {
			fmt.Printf("vk call: %s\n", endpoint.String())
		}
		gdelay.Wait()
		if resp, err = http.Get(endpoint.String()); err != nil {
			return err
		}
		defer resp.Body.Close()
		if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return err
		}
		if response.Err != nil && response.Err.Code != 0 {
			if response.Err.Code == ErrTooManyReq {
				return errTooManyTooManyReq
			}
			return response.Err
		}
		if Debug {
			fmt.Printf("vk api resp: %s\n", string(response.Response))
		}
		if err = json.Unmarshal(response.Response, out); err != nil {
			return err
		}
		return nil
	})
	return err
}

type ApiList struct {
	Count int         `json:"count"`
	Items interface{} `json:"items"`
}

func PublicAPI(method string, params url.Values, out interface{}) error {
	q := url.Values{
		"v":     {Version},
		"https": {strconv.Itoa(HTTPS)},
	}.Encode()
	endpoint, err := url.Parse(fmt.Sprint(APIURL, method, "?", q))
	if err != nil {
		return err
	}
	query := endpoint.Query()
	for k, v := range params {
		if len(v) > 0 {
			query.Set(k, v[0])
		}
	}
	endpoint.RawQuery = query.Encode()
	var (
		resp     *http.Response
		response struct {
			Err      *Error          `json:"error"`
			Response json.RawMessage `json:"response"`
		}
	)
	if resp, err = http.Get(endpoint.String()); err != nil {
		return err
	}
	defer resp.Body.Close()
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}
	if response.Err != nil && response.Err.Code != 0 {
		return response.Err
	}
	if err = json.Unmarshal(response.Response, out); err != nil {
		return err
	}
	return nil
}
