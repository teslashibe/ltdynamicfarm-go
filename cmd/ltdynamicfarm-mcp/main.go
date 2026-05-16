// Command ltdynamicfarm-mcp is a stdio MCP server exposing
// ltdynamicfarm-go to any MCP host (Cursor, Claude Desktop, etc.).
//
// Auth is loaded from ~/.ltdynamicfarm-mcp/config.json. Env vars take
// precedence when set:
//
//	LTDYNAMICFARM_EMAIL       — account email
//	LTDYNAMICFARM_PASSWORD    — account password
//	LTDYNAMICFARM_SESSION_ID  — cached ASP.NET_SessionId (optional, skips Login)
//
// Config file: ~/.ltdynamicfarm-mcp/config.json
//
//	{
//	  "email":      "you@example.com",
//	  "password":   "...",
//	  "session_id": ""
//	}
//
// Register with Cursor (`~/.cursor/mcp.json`):
//
//	{"mcpServers":{"ltdynamicfarm":{"command":"/Users/you/bin/ltdynamicfarm-mcp"}}}
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	ltdf "github.com/teslashibe/ltdynamicfarm-go"
	ltmcp "github.com/teslashibe/ltdynamicfarm-go/mcp"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type configFile struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	SessionID string `json:"session_id,omitempty"`
}

func defaultConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".ltdynamicfarm-mcp", "config.json")
}

func loadAuth() (ltdf.Auth, error) {
	var cfg configFile

	data, err := os.ReadFile(defaultConfigPath())
	if err != nil && !os.IsNotExist(err) {
		return ltdf.Auth{}, fmt.Errorf("read config: %w", err)
	}
	if data != nil {
		if err := json.Unmarshal(data, &cfg); err != nil {
			return ltdf.Auth{}, fmt.Errorf("parse config: %w", err)
		}
	}
	if v := os.Getenv("LTDYNAMICFARM_EMAIL"); v != "" {
		cfg.Email = v
	}
	if v := os.Getenv("LTDYNAMICFARM_PASSWORD"); v != "" {
		cfg.Password = v
	}
	if v := os.Getenv("LTDYNAMICFARM_SESSION_ID"); v != "" {
		cfg.SessionID = v
	}
	if cfg.SessionID == "" && (cfg.Email == "" || cfg.Password == "") {
		return ltdf.Auth{}, fmt.Errorf(
			"ltdynamicfarm credentials not found.\n"+
				"Set LTDYNAMICFARM_EMAIL+LTDYNAMICFARM_PASSWORD (or LTDYNAMICFARM_SESSION_ID), "+
				"or fill %s.", defaultConfigPath())
	}
	return ltdf.Auth{
		Email:     cfg.Email,
		Password:  cfg.Password,
		SessionID: cfg.SessionID,
	}, nil
}

func main() {
	log.SetOutput(os.Stderr)
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "ltdynamicfarm-mcp:", err)
		os.Exit(1)
	}
}

func run() error {
	auth, err := loadAuth()
	if err != nil {
		return err
	}
	client, err := ltdf.New(auth)
	if err != nil {
		return fmt.Errorf("init client: %w", err)
	}
	s := server.NewMCPServer("ltdynamicfarm-mcp", "0.1.0", server.WithToolCapabilities(true))
	for _, t := range (ltmcp.Provider{}).Tools() {
		t := t
		rawSchema, err := json.Marshal(t.InputSchema)
		if err != nil {
			return fmt.Errorf("marshal schema for %s: %w", t.Name, err)
		}
		tool := mcp.NewToolWithRawSchema(t.Name, t.Description, rawSchema)
		s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			raw, err := json.Marshal(req.Params.Arguments)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			result, invokeErr := t.Invoke(ctx, client, raw)
			if invokeErr != nil {
				return mcp.NewToolResultError(invokeErr.Error()), nil
			}
			out, err := json.Marshal(result)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultText(string(out)), nil
		})
	}
	return server.ServeStdio(s)
}
