package ltdynamicfarm

import (
	"context"
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
// "saleshistory.aspx", "probatelist", "ecampaignhistory".
//
// This is the catch-all reader that gives the agent access to any
// authenticated Dynamic Farm page without us having to write per-page
// parsers up-front.
func (c *Client) GetPageHTML(ctx context.Context, page string) (string, error) {
	body, _, err := c.getBytes(ctx, "/genericpage.aspx",
		url.Values{"page": []string{page}})
	if err != nil {
		return "", err
	}
	return string(body), nil
}
