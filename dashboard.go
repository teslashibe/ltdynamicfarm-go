package ltdynamicfarm

import (
	"context"
	"errors"
	"net/url"
	"regexp"
	"strings"
)

var titleRe = regexp.MustCompile(`(?is)<title>\s*(.*?)\s*</title>`)

// GetDashboard fetches /genericpage.aspx?page=dashboard.aspx and returns
// a normalized envelope. The HTML body itself is not included — call
// GetPageHTML for the raw markup.
func (c *Client) GetDashboard(ctx context.Context) (*DashboardPage, error) {
	body, status, err := c.getBytes(ctx, "/genericpage.aspx",
		url.Values{"page": []string{"dashboard.aspx"}})
	if err != nil {
		return nil, err
	}
	p := &DashboardPage{
		URL:          baseURL + "/genericpage.aspx?page=dashboard.aspx",
		StatusCode:   status,
		ContentBytes: len(body),
	}
	if m := titleRe.FindStringSubmatch(string(body)); len(m) > 1 {
		p.Title = strings.TrimSpace(m[1])
	}
	if m := dashboardGreetingRe.FindStringSubmatch(string(body)); len(m) > 1 {
		p.Greeting = strings.TrimSpace(m[1])
	}
	return p, nil
}

// GetPageHTML fetches an arbitrary GenericPage and returns raw HTML.
// page is a value like "dashboard.aspx", "farmlist.aspx",
// "probatelist", "ecampaignhistory".
//
// Note: page availability is site-variant dependent. Some documented
// pages (e.g. "saleshistory.aspx") are missing their server-side .html
// template on certain site variants (e.g. homeprofile-us) and will 5xx;
// in that case GetPageHTML returns ErrPageUnavailable rather than the
// raw HTML 500 body, so callers can distinguish "page unavailable" from
// auth/network errors.
//
// This is the catch-all reader that gives the agent access to any
// authenticated Dynamic Farm page without us having to write per-page
// parsers up-front.
func (c *Client) GetPageHTML(ctx context.Context, page string) (string, error) {
	body, status, err := c.getBytes(ctx, "/genericpage.aspx",
		url.Values{"page": []string{page}})
	if err != nil {
		// A 5xx surfaces as *HTTPError after retries are exhausted. Map it
		// to a typed ErrPageUnavailable so callers don't get a raw 500 body.
		var httpErr *HTTPError
		if errors.As(err, &httpErr) && httpErr.StatusCode >= 500 {
			return "", ErrPageUnavailable
		}
		return "", err
	}
	if status >= 500 {
		return "", ErrPageUnavailable
	}
	return string(body), nil
}
