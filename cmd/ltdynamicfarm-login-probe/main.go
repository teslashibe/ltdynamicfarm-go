// Command ltdynamicfarm-login-probe is a tiny smoke test that runs
// Login() against the configured credentials and prints the result.
package main

import (
	"context"
	"fmt"
	"os"

	ltdf "github.com/teslashibe/ltdynamicfarm-go"
)

func main() {
	auth := ltdf.Auth{
		Email:    os.Getenv("LTDYNAMICFARM_EMAIL"),
		Password: os.Getenv("LTDYNAMICFARM_PASSWORD"),
	}
	c, err := ltdf.New(auth)
	if err != nil {
		fmt.Fprintln(os.Stderr, "init:", err)
		os.Exit(1)
	}
	user, err := c.Login(context.Background())
	if err != nil {
		fmt.Fprintln(os.Stderr, "login:", err)
		os.Exit(1)
	}
	snap := c.AuthSnapshot()
	fmt.Printf("logged in as %s (logged_in=%v)\n", user.FirstName, user.LoggedIn)
	fmt.Printf("session_id=%s\n", snap.SessionID)
}
