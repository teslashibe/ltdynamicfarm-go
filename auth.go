package ltdynamicfarm

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// Login posts the configured Email + Password to /Login.aspx exactly
// the way the website's loginForm does:
//
//	POST /Login.aspx
//	Content-Type: application/x-www-form-urlencoded
//	prefix=LDF&Email=...&PassWord=...
//
// On success the server responds 302 → /genericpage.aspx?page=dashboard.aspx
// and sets an ASP.NET_SessionId cookie. We disable redirect-following for
// this single round-trip so curl-style POST-redirect-with-no-body 411s
// don't bury the cookie.
func (c *Client) Login(ctx context.Context) (*User, error) {
	c.authMu.RLock()
	email := c.auth.Email
	password := c.auth.Password
	prefix := c.auth.Prefix
	if prefix == "" {
		prefix = affiliatePrefix
	}
	c.authMu.RUnlock()

	if email == "" || password == "" {
		return nil, fmt.Errorf("%w: email/password required", ErrInvalidAuth)
	}

	form := url.Values{}
	form.Set("prefix", prefix)
	form.Set("Email", email)
	form.Set("PassWord", password)

	// Use the underlying transport with a one-shot client that does NOT
	// auto-follow redirects. The redirect target is what proves login;
	// it carries the session cookie.
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+loginPath,
		strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidParams, err)
	}
	c.setCommonHeaders(req, "application/x-www-form-urlencoded")

	noRedir := *c.httpClient
	noRedir.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	c.waitForGap(ctx)
	resp, err := noRedir.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrLoginFailed, err)
	}
	defer resp.Body.Close()

	c.captureSessionCookie(resp)

	switch resp.StatusCode {
	case http.StatusFound, http.StatusSeeOther, http.StatusMovedPermanently:
		loc := resp.Header.Get("Location")
		if !strings.Contains(strings.ToLower(loc), "dashboard") {
			return nil, fmt.Errorf("%w: unexpected redirect to %q", ErrLoginFailed, loc)
		}
	case http.StatusOK:
		// Some failure modes return 200 with the login form re-rendered.
		// Treat as failure.
		return nil, fmt.Errorf("%w: server returned 200 (likely bad email/password)", ErrLoginFailed)
	default:
		return nil, fmt.Errorf("%w: HTTP %d", ErrLoginFailed, resp.StatusCode)
	}

	c.authMu.RLock()
	sid := c.auth.SessionID
	c.authMu.RUnlock()
	if sid == "" {
		return nil, fmt.Errorf("%w: server did not set ASP.NET_SessionId", ErrLoginFailed)
	}

	return c.GetMe(ctx)
}

// dashboardGreetingRe finds the first-name span the dashboard renders in
// the header: <span class="hidden-xs">Margie</span>
var dashboardGreetingRe = regexp.MustCompile(`(?i)<span class="hidden-xs">\s*([^<\s][^<]*?)\s*</span>`)

// GetMe fetches the dashboard page and extracts the displayed first
// name as a stand-in for "current user". Dynamic Farm exposes no JSON
// /me endpoint.
func (c *Client) GetMe(ctx context.Context) (*User, error) {
	body, _, err := c.getBytes(ctx, "/genericpage.aspx", url.Values{"page": []string{"dashboard.aspx"}})
	if err != nil {
		return nil, fmt.Errorf("GetMe: %w", err)
	}
	html := string(body)
	signedOut := !strings.Contains(strings.ToLower(html), "sign out")

	u := &User{LoggedIn: !signedOut}
	if m := dashboardGreetingRe.FindStringSubmatch(html); len(m) > 1 {
		u.FirstName = strings.TrimSpace(m[1])
	}
	if signedOut {
		return u, fmt.Errorf("%w: dashboard does not show a Sign Out link", ErrUnauthorized)
	}
	return u, nil
}
