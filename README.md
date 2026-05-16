# ltdynamicfarm-go

Private Go client + MCP server for [`ltdynamicfarm.com`](https://ltdynamicfarm.com)
(Lawyers Title "Dynamic Farm"). There is no public API — this is
reverse-engineered from the website.

## Authentication

`POST /Login.aspx` (`application/x-www-form-urlencoded`) with
`prefix=LDF&Email=...&PassWord=...`. On success the server responds 302
→ `/genericpage.aspx?page=dashboard.aspx` and sets an
`ASP.NET_SessionId` cookie. Subsequent requests carry that cookie.

## Quick start (library)

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    ltdf "github.com/teslashibe/ltdynamicfarm-go"
)

func main() {
    c, err := ltdf.New(ltdf.Auth{
        Email:    os.Getenv("LTDYNAMICFARM_EMAIL"),
        Password: os.Getenv("LTDYNAMICFARM_PASSWORD"),
    })
    if err != nil { log.Fatal(err) }

    me, err := c.Login(context.Background())
    if err != nil { log.Fatal(err) }
    fmt.Printf("logged in as %s\n", me.FirstName)

    page, err := c.GetDashboard(context.Background())
    if err != nil { log.Fatal(err) }
    fmt.Printf("dashboard title=%q greeting=%q bytes=%d\n",
        page.Title, page.Greeting, page.ContentBytes)
}
```

## Supported operations (v0.1)

| Area      | Methods                                       |
|-----------|-----------------------------------------------|
| Auth      | `Login`, `GetMe`, `AuthSnapshot`              |
| Pages     | `GetDashboard`, `GetPageHTML`                 |

`GetPageHTML(ctx, page)` is the generic reader. Useful `page` values:

- `dashboard.aspx`
- `farmlist.aspx`
- `farmsprofiles.aspx`
- `saleshistory.aspx`
- `probatelist`
- `ecampaignhistory`
- `ecampaign`
- `marketing`
- `facebookad`
- `reverseafarm.aspx`
- `previousreverses`
- `reversephoneemail`
- `prop19ttf`
- `ticklers`
- `webinars` / `pastwebinars`
- `chatgptinterface`
- `subjectlinegrader` / `subjectlinerater`
- `analyticareareportinterface`

## TODO — endpoints to model with typed parsers

- [ ] `farmlist.aspx` → typed `[]FarmListItem`
- [ ] `saleshistory.aspx` → typed `[]Sale`
- [ ] `probatelist` → typed `[]ProbateRecord`
- [ ] `ecampaignhistory` → typed `[]Campaign`
- [ ] `PropertyData.aspx` AJAX → typed `Property`
- [ ] `getSaved.aspx` AJAX → typed saved-search list
- [ ] `getCounties.aspx` AJAX → typed county lookup
- [ ] `Heatmap.aspx?groupmode=Group` → typed heatmap data
- [ ] `Settings.aspx` → typed user profile

## MCP server

Build and install:

```bash
go install github.com/teslashibe/ltdynamicfarm-go/cmd/ltdynamicfarm-mcp@latest
```

Configure credentials at `~/.ltdynamicfarm-mcp/config.json`:

```json
{
  "email":    "you@example.com",
  "password": "..."
}
```

Or use env vars: `LTDYNAMICFARM_EMAIL`, `LTDYNAMICFARM_PASSWORD`,
`LTDYNAMICFARM_SESSION_ID`.

Register with Cursor (`~/.cursor/mcp.json`):

```json
{
  "mcpServers": {
    "ltdynamicfarm": { "command": "/Users/you/go/bin/ltdynamicfarm-mcp" }
  }
}
```

## License

Private. Internal use only.
