// Package mcp exposes ltdynamicfarm-go as a set of MCP tools.
package mcp

import "github.com/teslashibe/mcptool"

// Provider implements [mcptool.Provider] for ltdynamicfarm-go.
type Provider struct{}

// Platform returns "ltdynamicfarm".
func (Provider) Platform() string { return "ltdynamicfarm" }

// Tools returns every MCP tool, in registration order.
func (Provider) Tools() []mcptool.Tool {
	out := make([]mcptool.Tool, 0, len(authTools)+len(pageTools))
	out = append(out, authTools...)
	out = append(out, pageTools...)
	return out
}
