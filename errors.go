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

	// ErrPageUnavailable is returned when a GenericPage exists in our
	// documented list but the server cannot render it for the current
	// site variant — typically a 5xx caused by a missing server-side
	// .html template (e.g. saleshistory.aspx on homeprofile-us). It lets
	// callers distinguish "page unavailable" from auth/network errors.
	ErrPageUnavailable = errors.New("ltdynamicfarm: page unavailable for this site variant")

	// ErrEndpointNotConfigured is returned by the structured data methods
	// (GetFarmList/GetProbateList/GetECampaignHistory) while their
	// DataTables AJAX endpoint paths/params are still placeholders awaiting
	// a live DevTools capture. It fails loudly rather than guessing a URL.
	ErrEndpointNotConfigured = errors.New("ltdynamicfarm: data endpoint not yet configured (awaiting live capture)")
)

// HTTPError is returned for unexpected non-2xx HTTP responses.
type HTTPError struct {
	StatusCode int
	Body       string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("ltdynamicfarm: HTTP %d: %s", e.StatusCode, e.Body)
}
