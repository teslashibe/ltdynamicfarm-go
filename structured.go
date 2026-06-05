package ltdynamicfarm

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// Structured (DataTables-backed) data access.
//
// The list views on ltdynamicfarm.com (farmlist.aspx, probatelist,
// ecampaignhistory) render their rows client-side via DataTables AJAX
// calls; the GenericPage HTML returned by GetPageHTML is only the page
// shell and contains no row data. The methods below are the structured
// data path: they hit the underlying AJAX endpoints directly with the
// active session cookie and unmarshal the DataTables JSON into typed
// structs.
//
// TODO(live-capture): the exact AJAX endpoint paths and POST/GET params
// are NOT yet known — they must be captured from a live, authenticated
// browser session (DevTools → Network) before these methods can return
// real data. The constants below are intentionally empty placeholders.
// While they are empty, every method returns ErrEndpointNotConfigured
// so we fail loudly instead of silently hitting a wrong/guessed URL.
//
// Open questions to resolve from the capture (see teslashibe/smore#140):
//   - exact endpoint path for each list view
//   - GET query vs form POST, and the param names
//   - server-side pagination params (start/length) vs all-rows
//   - whether a CSRF/anti-forgery token must be extracted first
//   - the JSON envelope shape (DataTables default is {data:[...]}, but
//     ASP.NET handlers often wrap in {d:"<json-string>"})
const (
	// farmListEndpoint is the DataTables AJAX path behind farmlist.aspx.
	// TODO(live-capture): e.g. "/PropertyData.aspx" or a dedicated handler.
	farmListEndpoint = ""
	// probateListEndpoint is the DataTables AJAX path behind probatelist.
	// TODO(live-capture).
	probateListEndpoint = ""
	// ecampaignHistoryEndpoint is the DataTables AJAX path behind
	// ecampaignhistory. TODO(live-capture).
	ecampaignHistoryEndpoint = ""
)

// dataTablesEnvelope models the common DataTables server-side response
// shape. Real LTDF responses may differ (e.g. an ASP.NET {"d":"..."}
// wrapper); reconcile against the live capture.
type dataTablesEnvelope struct {
	Data            json.RawMessage `json:"data"`
	RecordsTotal    int             `json:"recordsTotal,omitempty"`
	RecordsFiltered int             `json:"recordsFiltered,omitempty"`
}

// fetchDataTable performs the shared GET + DataTables-envelope decode for
// the structured list methods. It returns the raw `data` array bytes so
// each caller can unmarshal into its own typed slice.
func (c *Client) fetchDataTable(ctx context.Context, endpoint string, params ListParams) (json.RawMessage, error) {
	if endpoint == "" {
		// Placeholder not yet filled in from a live DevTools capture.
		return nil, ErrEndpointNotConfigured
	}

	q := url.Values{}
	if params.FarmArea != "" {
		q.Set("farmArea", params.FarmArea) // TODO(live-capture): real param name
	}
	if params.County != "" {
		q.Set("county", params.County) // TODO(live-capture): real param name
	}
	if params.Start > 0 {
		q.Set("start", fmt.Sprintf("%d", params.Start)) // TODO(live-capture)
	}
	if params.Length > 0 {
		q.Set("length", fmt.Sprintf("%d", params.Length)) // TODO(live-capture)
	}

	body, status, err := c.getBytes(ctx, endpoint, q)
	if err != nil {
		return nil, err
	}
	if status >= 500 {
		return nil, ErrPageUnavailable
	}

	var env dataTablesEnvelope
	if err := json.Unmarshal(body, &env); err != nil {
		return nil, fmt.Errorf("%w: decoding DataTables envelope: %v", ErrRequestFailed, err)
	}
	return env.Data, nil
}

// GetFarmList returns the homeowner/property records for the user's farm
// area(s) as typed structs.
//
// TODO(live-capture): returns ErrEndpointNotConfigured until
// farmListEndpoint is filled in from a live capture.
func (c *Client) GetFarmList(ctx context.Context, params ListParams) ([]FarmListItem, error) {
	data, err := c.fetchDataTable(ctx, farmListEndpoint, params)
	if err != nil {
		return nil, err
	}
	var out []FarmListItem
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("%w: decoding farm list rows: %v", ErrRequestFailed, err)
	}
	return out, nil
}

// GetProbateList returns probate leads as typed structs.
//
// TODO(live-capture): returns ErrEndpointNotConfigured until
// probateListEndpoint is filled in from a live capture.
func (c *Client) GetProbateList(ctx context.Context, params ListParams) ([]ProbateLeadItem, error) {
	data, err := c.fetchDataTable(ctx, probateListEndpoint, params)
	if err != nil {
		return nil, err
	}
	var out []ProbateLeadItem
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("%w: decoding probate rows: %v", ErrRequestFailed, err)
	}
	return out, nil
}

// GetECampaignHistory returns e-campaign sends + engagement stats as
// typed structs.
//
// TODO(live-capture): returns ErrEndpointNotConfigured until
// ecampaignHistoryEndpoint is filled in from a live capture.
func (c *Client) GetECampaignHistory(ctx context.Context, params ListParams) ([]ECampaignRecord, error) {
	data, err := c.fetchDataTable(ctx, ecampaignHistoryEndpoint, params)
	if err != nil {
		return nil, err
	}
	var out []ECampaignRecord
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("%w: decoding e-campaign rows: %v", ErrRequestFailed, err)
	}
	return out, nil
}
