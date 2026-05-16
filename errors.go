package ltdynamicfarm

import (
	"errors"
	"fmt"
)

// Sentinel errors.
var (
	ErrInvalidAuth   = errors.New("ltdynamicfarm: missing or invalid auth credentials")
	ErrUnauthorized  = errors.New("ltdynamicfarm: unauthorized (session expired)")
	ErrForbidden     = errors.New("ltdynamicfarm: forbidden")
	ErrNotFound      = errors.New("ltdynamicfarm: not found")
	ErrRateLimited   = errors.New("ltdynamicfarm: rate limited")
	ErrInvalidParams = errors.New("ltdynamicfarm: invalid parameters")
	ErrRequestFailed = errors.New("ltdynamicfarm: request failed")
	ErrLoginFailed   = errors.New("ltdynamicfarm: login failed (bad email/password, or 2FA required)")
)

// HTTPError is returned for unexpected non-2xx HTTP responses.
type HTTPError struct {
	StatusCode int
	Body       string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("ltdynamicfarm: HTTP %d: %s", e.StatusCode, e.Body)
}
