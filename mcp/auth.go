package mcp

import (
	"context"

	ltdf "github.com/teslashibe/ltdynamicfarm-go"
	"github.com/teslashibe/mcptool"
)

// LoginInput is the typed input for ltdynamicfarm_login.
type LoginInput struct{}

func login(ctx context.Context, c *ltdf.Client, _ LoginInput) (any, error) {
	return c.Login(ctx)
}

// GetMeInput is the typed input for ltdynamicfarm_get_me.
type GetMeInput struct{}

func getMe(ctx context.Context, c *ltdf.Client, _ GetMeInput) (any, error) {
	return c.GetMe(ctx)
}

var authTools = []mcptool.Tool{
	mcptool.Define[*ltdf.Client, LoginInput](
		"ltdynamicfarm_login",
		"Authenticate against ltdynamicfarm.com using the configured email/password and cache the session.",
		"Login",
		login,
	),
	mcptool.Define[*ltdf.Client, GetMeInput](
		"ltdynamicfarm_get_me",
		"Confirm the cached session is alive and return the displayed first name from the dashboard header.",
		"GetMe",
		getMe,
	),
}
