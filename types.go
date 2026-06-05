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

// FarmListItem is one homeowner/property record from a farm area, as
// rendered client-side by the farmlist.aspx DataTable.
//
// NOTE: field names are the intended shape; exact JSON tags must be
// reconciled against a live DataTables capture (see structured.go's
// TODO(live-capture) markers) before the parser can be trusted.
type FarmListItem struct {
	ID           string `json:"id,omitempty"`
	OwnerName    string `json:"owner_name,omitempty"`
	Address      string `json:"address,omitempty"`
	City         string `json:"city,omitempty"`
	State        string `json:"state,omitempty"`
	Zip          string `json:"zip,omitempty"`
	PropertyType string `json:"property_type,omitempty"`
	FarmArea     string `json:"farm_area,omitempty"`
	APN          string `json:"apn,omitempty"`
}

// ProbateLeadItem is one probate lead from the probatelist DataTable.
//
// NOTE: field names are the intended shape pending a live capture.
type ProbateLeadItem struct {
	ID           string `json:"id,omitempty"`
	DecedentName string `json:"decedent_name,omitempty"`
	Address      string `json:"address,omitempty"`
	City         string `json:"city,omitempty"`
	State        string `json:"state,omitempty"`
	Zip          string `json:"zip,omitempty"`
	County       string `json:"county,omitempty"`
	FilingDate   string `json:"filing_date,omitempty"`
	CaseNumber   string `json:"case_number,omitempty"`
}

// ECampaignRecord is one e-campaign send + engagement stats from the
// ecampaignhistory DataTable.
//
// NOTE: field names are the intended shape pending a live capture.
type ECampaignRecord struct {
	ID             string  `json:"id,omitempty"`
	CampaignName   string  `json:"campaign_name,omitempty"`
	SendDate       string  `json:"send_date,omitempty"`
	RecipientCount int     `json:"recipient_count,omitempty"`
	Opens          int     `json:"opens,omitempty"`
	Clicks         int     `json:"clicks,omitempty"`
	OpenRate       float64 `json:"open_rate,omitempty"`
	ClickRate      float64 `json:"click_rate,omitempty"`
}

// ListParams carries optional pagination/filter inputs for the
// structured DataTables-backed list methods. Fields are best-effort;
// the exact server-side param names are TBD (see structured.go).
type ListParams struct {
	// FarmArea optionally filters farm-list results to one farm area.
	FarmArea string
	// County optionally filters probate results to one county.
	County string
	// Start is the zero-based row offset for server-side pagination.
	Start int
	// Length is the max number of rows to return (0 = server default).
	Length int
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
