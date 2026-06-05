// Package mcp exposes ltdynamicfarm-go as a set of MCP tools.
package mcp

import "github.com/teslashibe/mcptool"

// Provider implements [mcptool.Provider] for ltdynamicfarm-go.
type Provider struct{}

// Platform returns "ltdynamicfarm".
func (Provider) Platform() string { return "ltdynamicfarm" }

// Tools returns every MCP tool, in registration order.
//
// NOTE (#140): structuredTools (get_farm_list / get_probate_list /
// get_ecampaign_history) are intentionally NOT exposed yet. Live probing
// showed the three list views use heterogeneous, undocumented rendering
// (server-rendered table after a search POST, JS-injected table, and an
// /el/AutoLogNET.aspx iframe respectively) rather than a single DataTables
// JSON endpoint, so the client methods currently return
// ErrEndpointNotConfigured. Exposing always-erroring tools to the agent
// would be worse than not having them. The scaffolding (structured.go +
// types) is kept so the follow-up can fill in the endpoints and re-enable
// here once the per-view flows are mapped. See teslashibe/smore#140.
func (Provider) Tools() []mcptool.Tool {
	out := make([]mcptool.Tool, 0, len(authTools)+len(pageTools))
	out = append(out, authTools...)
	out = append(out, pageTools...)
	return out
}

// structuredToolsDisabled keeps a reference to the not-yet-exposed #140 tools
// so the var (and structured.go) remain compiled and ready to re-enable.
var _ = structuredTools
