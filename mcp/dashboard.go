package mcp

import (
	"context"

	ltdf "github.com/teslashibe/ltdynamicfarm-go"
	"github.com/teslashibe/mcptool"
)

// GetDashboardInput is the typed input for ltdynamicfarm_get_dashboard.
type GetDashboardInput struct{}

func getDashboard(ctx context.Context, c *ltdf.Client, _ GetDashboardInput) (any, error) {
	return c.GetDashboard(ctx)
}

// GetPageHTMLInput is the typed input for ltdynamicfarm_get_page_html.
type GetPageHTMLInput struct {
	Page string `json:"page" jsonschema:"description=GenericPage query value e.g. 'dashboard.aspx' 'farmlist.aspx' 'saleshistory.aspx' 'probatelist' 'ecampaignhistory',required"`
}

func getPageHTML(ctx context.Context, c *ltdf.Client, in GetPageHTMLInput) (any, error) {
	html, err := c.GetPageHTML(ctx, in.Page)
	if err != nil {
		return nil, err
	}
	return map[string]any{"page": in.Page, "html": html}, nil
}

var pageTools = []mcptool.Tool{
	mcptool.Define[*ltdf.Client, GetDashboardInput](
		"ltdynamicfarm_get_dashboard",
		"Fetch a normalized snapshot of the Dynamic Farm dashboard (title, greeting, response size).",
		"GetDashboard",
		getDashboard,
	),
	mcptool.Define[*ltdf.Client, GetPageHTMLInput](
		"ltdynamicfarm_get_page_html",
		"Fetch the raw HTML of any GenericPage by name (e.g. farmlist.aspx, saleshistory.aspx, probatelist).",
		"GetPageHTML",
		getPageHTML,
	),
}
