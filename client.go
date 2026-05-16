package ltdynamicfarm

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// getBytes fetches a path under baseURL and returns the raw response body.
// Useful for HTML endpoints that have no JSON equivalent.
func (c *Client) getBytes(ctx context.Context, path string, query url.Values) ([]byte, int, error) {
	full := baseURL + path
	if len(query) > 0 {
		full += "?" + query.Encode()
	}
	return c.makeRequest(ctx, http.MethodGet, full, nil, "")
}

// postForm submits a form-encoded POST body.
func (c *Client) postForm(ctx context.Context, path string, form url.Values) ([]byte, int, error) {
	body := []byte(form.Encode())
	return c.makeRequest(ctx, http.MethodPost, baseURL+path, body, "application/x-www-form-urlencoded")
}

// makeRequest performs an HTTP request with automatic retry on 429 and
// 5xx, after lazily bootstrapping a session if needed.
func (c *Client) makeRequest(ctx context.Context, method, rawURL string, body []byte, contentType string) ([]byte, int, error) {
	if err := c.ensureLoggedIn(ctx); err != nil {
		return nil, 0, err
	}
	return c.doRetried(ctx, method, rawURL, body, contentType)
}

func (c *Client) doRetried(ctx context.Context, method, rawURL string, body []byte, contentType string) ([]byte, int, error) {
	var lastErr error
	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		if attempt > 0 {
			wait := c.backoff(attempt)
			select {
			case <-ctx.Done():
				return nil, 0, ctx.Err()
			case <-time.After(wait):
			}
		}
		raw, status, err := c.doRequest(ctx, method, rawURL, body, contentType)
		if err == nil {
			return raw, status, nil
		}
		lastErr = err
		if errors.Is(err, ErrRateLimited) {
			continue
		}
		var httpErr *HTTPError
		if errors.As(err, &httpErr) && httpErr.StatusCode >= 500 {
			continue
		}
		return nil, status, err
	}
	return nil, 0, lastErr
}

// ensureLoggedIn calls Login lazily if Auth was constructed with email
// + password but no SessionID yet.
func (c *Client) ensureLoggedIn(ctx context.Context) error {
	c.authMu.RLock()
	has := c.auth.SessionID != ""
	c.authMu.RUnlock()
	if has {
		return nil
	}
	c.loginOnce.Do(func() {
		_, c.loginErr = c.Login(ctx)
	})
	return c.loginErr
}

// doRequest performs a single HTTP round-trip. It returns body, status,
// and an error (which may also carry a non-zero status when the response
// was non-2xx).
func (c *Client) doRequest(ctx context.Context, method, rawURL string, body []byte, contentType string) ([]byte, int, error) {
	c.waitForGap(ctx)
	if ctx.Err() != nil {
		return nil, 0, ctx.Err()
	}

	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, rawURL, bodyReader)
	if err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrRequestFailed, err)
	}
	c.setCommonHeaders(req, contentType)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrRequestFailed, err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("%w: reading body: %v", ErrRequestFailed, err)
	}

	c.captureSessionCookie(resp)

	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusNoContent,
		http.StatusMovedPermanently, http.StatusFound, http.StatusSeeOther,
		http.StatusTemporaryRedirect, http.StatusPermanentRedirect:
		return raw, resp.StatusCode, nil
	case http.StatusUnauthorized:
		return nil, resp.StatusCode, ErrUnauthorized
	case http.StatusForbidden:
		return nil, resp.StatusCode, ErrForbidden
	case http.StatusNotFound:
		return nil, resp.StatusCode, ErrNotFound
	case http.StatusTooManyRequests:
		wait := parseRetryAfter(resp.Header.Get("Retry-After"), 60*time.Second)
		c.gapMu.Lock()
		if earliest := time.Now().Add(wait); c.lastReqAt.Before(earliest) {
			c.lastReqAt = earliest
		}
		c.gapMu.Unlock()
		return nil, resp.StatusCode, fmt.Errorf("%w: retry after %s", ErrRateLimited, wait)
	default:
		return nil, resp.StatusCode, &HTTPError{StatusCode: resp.StatusCode, Body: truncate(string(raw), 256)}
	}
}

func (c *Client) setCommonHeaders(req *http.Request, contentType string) {
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Referer", baseURL+"/index_dynamicfarm.shtml")
	req.Header.Set("Origin", baseURL)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	if cookie := c.cookieString(); cookie != "" {
		req.Header.Set("Cookie", cookie)
	}
}

// cookieString assembles the Cookie header from the cached SessionID
// plus anything the jar has accumulated.
func (c *Client) cookieString() string {
	var parts []string
	c.authMu.RLock()
	if c.auth.SessionID != "" {
		parts = append(parts, "ASP.NET_SessionId="+c.auth.SessionID)
	}
	c.authMu.RUnlock()

	if c.httpClient != nil && c.httpClient.Jar != nil {
		u, _ := url.Parse(baseURL)
		for _, ck := range c.httpClient.Jar.Cookies(u) {
			if ck.Name == "ASP.NET_SessionId" {
				continue // already added above (or about to be by Set-Cookie capture)
			}
			parts = append(parts, ck.Name+"="+ck.Value)
		}
	}
	return strings.Join(parts, "; ")
}

// captureSessionCookie reads Set-Cookie headers from a response and
// caches the ASP.NET_SessionId value onto Auth.
func (c *Client) captureSessionCookie(resp *http.Response) {
	c.authMu.Lock()
	defer c.authMu.Unlock()
	for _, ck := range resp.Cookies() {
		if ck.Name == "ASP.NET_SessionId" && ck.Value != "" {
			c.auth.SessionID = ck.Value
		}
	}
}

func (c *Client) backoff(attempt int) time.Duration {
	return time.Duration(math.Pow(2, float64(attempt-1))) * c.retryBase
}

func (c *Client) waitForGap(ctx context.Context) {
	c.gapMu.Lock()
	now := time.Now()
	next := c.lastReqAt.Add(c.minGap)
	if now.After(next) {
		next = now
	}
	c.lastReqAt = next
	c.gapMu.Unlock()
	if wait := time.Until(next); wait > 0 {
		select {
		case <-ctx.Done():
		case <-time.After(wait):
		}
	}
}

func parseRetryAfter(val string, fallback time.Duration) time.Duration {
	if val == "" {
		return fallback
	}
	trimmed := strings.TrimSpace(val)
	if t, err := http.ParseTime(trimmed); err == nil {
		if d := time.Until(t); d > 0 {
			return d
		}
	}
	return fallback
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
