// Package ltdynamicfarm is a Go client + MCP tool surface for
// ltdynamicfarm.com (Lawyers Title "Dynamic Farm").
//
// There is no public API. This client targets the same ASP.NET endpoints
// the website itself hits, starting with a form-encoded POST to
// /Login.aspx and a sticky ASP.NET_SessionId cookie.
//
// The site is hosted on Microsoft-IIS/10.0 and does not appear to use
// CDN-level bot protection as of 2026-05.
package ltdynamicfarm

import (
	"net/http"
	"net/http/cookiejar"
	"sync"
	"time"
)

const (
	baseURL          = "https://ltdynamicfarm.com"
	loginPath        = "/Login.aspx"
	defaultUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36"
	defaultRetries   = 3
	defaultRetryBase = 500 * time.Millisecond
	// affiliatePrefix is the hidden form value the site requires alongside
	// Email + PassWord. "LDF" = Lawyers Title Dynamic Farm.
	affiliatePrefix = "LDF"
)

// Client talks to ltdynamicfarm.com.
type Client struct {
	auth       Auth
	httpClient *http.Client
	userAgent  string
	maxRetries int
	retryBase  time.Duration
	minGap     time.Duration

	gapMu     sync.Mutex
	lastReqAt time.Time

	loginOnce sync.Once
	loginErr  error

	authMu sync.RWMutex
}

// Option configures a Client.
type Option func(*Client)

// WithUserAgent overrides the default browser User-Agent string.
func WithUserAgent(ua string) Option { return func(c *Client) { c.userAgent = ua } }

// WithRetry sets the maximum retry count and base backoff duration.
func WithRetry(maxRetries int, base time.Duration) Option {
	return func(c *Client) {
		c.maxRetries = maxRetries
		c.retryBase = base
	}
}

// WithHTTPClient overrides the default http.Client. Nil is ignored.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) {
		if hc != nil {
			c.httpClient = hc
		}
	}
}

// WithMinRequestGap sets the minimum time between consecutive requests.
// Defaults to 400ms.
func WithMinRequestGap(d time.Duration) Option {
	return func(c *Client) { c.minGap = d }
}

// New creates a new client. Either (Email + Password) or SessionID must
// be supplied on Auth.
func New(auth Auth, opts ...Option) (*Client, error) {
	if auth.SessionID == "" && (auth.Email == "" || auth.Password == "") {
		return nil, ErrInvalidAuth
	}
	jar, _ := cookiejar.New(nil)
	c := &Client{
		auth:       auth,
		httpClient: &http.Client{Timeout: 30 * time.Second, Jar: jar},
		userAgent:  defaultUserAgent,
		maxRetries: defaultRetries,
		retryBase:  defaultRetryBase,
		minGap:     400 * time.Millisecond,
	}
	for _, o := range opts {
		o(c)
	}
	return c, nil
}

// AuthSnapshot returns the current cached auth credentials. Useful for
// persisting an authenticated session to disk between runs.
func (c *Client) AuthSnapshot() Auth {
	c.authMu.RLock()
	defer c.authMu.RUnlock()
	return c.auth
}
