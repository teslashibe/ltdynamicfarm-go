package ltdynamicfarm

// Auth carries the credentials needed to talk to ltdynamicfarm.com.
//
// Either (Email + Password) or SessionID must be supplied. The first
// authenticated request will lazily POST /Login.aspx if no SessionID is
// cached yet, then store the ASP.NET_SessionId cookie returned by the
// site.
type Auth struct {
	Email    string
	Password string

	// SessionID is the value of the ASP.NET_SessionId cookie issued by
	// ltdynamicfarm.com on a successful login. Populated automatically
	// after Login(); supply it manually to skip the bootstrap.
	SessionID string

	// Prefix overrides the hidden "prefix" form value sent with the
	// login POST. Defaults to "LDF" (Lawyers Title Dynamic Farm) which
	// is what the website itself sends.
	Prefix string
}

// User is the bare-minimum profile surfaced by the dashboard. Dynamic
// Farm has no /me JSON endpoint, so GetMe parses the HTML.
type User struct {
	// FirstName is the name displayed in the dashboard header
	// (typically the agent's first name).
	FirstName string `json:"first_name"`

	// LoggedIn is true if the most recent dashboard fetch produced a
	// page containing the "Sign out" link.
	LoggedIn bool `json:"logged_in"`
}

// FarmListItem is one entry from the user's saved farms list.
// Populated lazily as endpoints are reverse-engineered.
type FarmListItem struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Count int    `json:"count,omitempty"`
}

// DashboardPage is the raw HTML body of /genericpage.aspx?page=dashboard.aspx
// together with a few normalized fields.
type DashboardPage struct {
	URL          string `json:"url"`
	StatusCode   int    `json:"status_code"`
	ContentBytes int    `json:"content_bytes"`
	Title        string `json:"title,omitempty"`
	Greeting     string `json:"greeting,omitempty"`
}
