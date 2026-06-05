package mcp

import (
	"context"

	ltdf "github.com/teslashibe/ltdynamicfarm-go"
	"github.com/teslashibe/mcptool"
)

// GetFarmListInput is the typed input for ltdynamicfarm_get_farm_list.
type GetFarmListInput struct {
	FarmArea string `json:"farm_area,omitempty" jsonschema:"description=Optional farm-area filter"`
	Start    int    `json:"start,omitempty" jsonschema:"description=Zero-based row offset for pagination"`
	Length   int    `json:"length,omitempty" jsonschema:"description=Max rows to return (0 = server default)"`
}

func getFarmList(ctx context.Context, c *ltdf.Client, in GetFarmListInput) (any, error) {
	return c.GetFarmList(ctx, ltdf.ListParams{FarmArea: in.FarmArea, Start: in.Start, Length: in.Length})
}

// GetProbateListInput is the typed input for ltdynamicfarm_get_probate_list.
type GetProbateListInput struct {
	County string `json:"county,omitempty" jsonschema:"description=Optional county filter"`
	Start  int    `json:"start,omitempty" jsonschema:"description=Zero-based row offset for pagination"`
	Length int    `json:"length,omitempty" jsonschema:"description=Max rows to return (0 = server default)"`
}

func getProbateList(ctx context.Context, c *ltdf.Client, in GetProbateListInput) (any, error) {
	return c.GetProbateList(ctx, ltdf.ListParams{County: in.County, Start: in.Start, Length: in.Length})
}

// GetECampaignHistoryInput is the typed input for ltdynamicfarm_get_ecampaign_history.
type GetECampaignHistoryInput struct {
	Start  int `json:"start,omitempty" jsonschema:"description=Zero-based row offset for pagination"`
	Length int `json:"length,omitempty" jsonschema:"description=Max rows to return (0 = server default)"`
}

func getECampaignHistory(ctx context.Context, c *ltdf.Client, in GetECampaignHistoryInput) (any, error) {
	return c.GetECampaignHistory(ctx, ltdf.ListParams{Start: in.Start, Length: in.Length})
}

var structuredTools = []mcptool.Tool{
	mcptool.Define[*ltdf.Client, GetFarmListInput](
		"ltdynamicfarm_get_farm_list",
		"Return farm-area homeowner/property records as a typed JSON array (not raw HTML).",
		"GetFarmList",
		getFarmList,
	),
	mcptool.Define[*ltdf.Client, GetProbateListInput](
		"ltdynamicfarm_get_probate_list",
		"Return probate leads (decedent, property, county, filing date) as a typed JSON array.",
		"GetProbateList",
		getProbateList,
	),
	mcptool.Define[*ltdf.Client, GetECampaignHistoryInput](
		"ltdynamicfarm_get_ecampaign_history",
		"Return e-campaign sends + open/click stats as a typed JSON array.",
		"GetECampaignHistory",
		getECampaignHistory,
	),
}
