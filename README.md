# Incident.io Go Client

A Go client library for the [Incident.io API](https://api-docs.incident.io/).

## Installation

```bash
go get github.com/cpanato/go-incident-io
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    incidentio "github.com/cpanato/go-incident-io"
)

func main() {
    // Create a new client with your API key
    client := incidentio.NewClient("YOUR-API-KEY-HERE")

    ctx := context.Background()

    // List all incidents
    incidents, _, err := client.Incidents.List(ctx, nil)
    if err != nil {
        log.Fatal(err)
    }

    for _, incident := range incidents {
        fmt.Printf("Incident: %s - %s\n", incident.ID, incident.Name)
    }
}
```

## Authentication

The client requires an API key for authentication. You can generate an API key from your [Incident.io dashboard](https://app.incident.io/).

```go
client := incidentio.NewClient("YOUR-API-KEY-HERE")
```

## Configuration Options

The client supports several configuration options:

```go
// Use a custom HTTP client
httpClient := &http.Client{
    Timeout: 60 * time.Second,
}

client := incidentio.NewClient("YOUR-API-KEY-HERE",
    incidentio.WithHTTPClient(httpClient))

// Use a custom base URL (e.g., for testing)
client := incidentio.NewClient("YOUR-API-KEY-HERE",
    incidentio.WithBaseURL("https://api.staging.incident.io"))
```

## Examples

### Creating an Incident

```go
ctx := context.Background()

// First, get the required IDs
incidentTypes, _, _ := client.IncidentTypes.List(ctx)
severities, _, _ := client.Severities.List(ctx)

// Create the incident
opts := &incidentio.CreateIncidentOptions{
    Name:           "Database Connection Issues",
    Summary:        "Users are experiencing intermittent connection timeouts",
    IncidentTypeID: incidentTypes[0].ID,
    SeverityID:     severities[0].ID,
    Mode:           "real", // or "test"
    Visibility:     "public", // or "private"
}

incident, _, err := client.Incidents.Create(ctx, opts)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Created incident: %s\n", incident.ID)
```

### Updating an Incident

```go
status := "resolved"
updateOpts := &incidentio.UpdateIncidentOptions{
    Status: &status,
    Summary: &"Issue has been resolved by restarting the database connection pool",
}

incident, _, err := client.Incidents.Update(ctx, "incident-id", updateOpts)
if err != nil {
    log.Fatal(err)
}
```

### Assigning Roles

```go
// List available roles
roles, _, _ := client.IncidentRoles.List(ctx)
users, _, _ := client.Users.List(ctx, nil)

// Assign roles during incident creation
opts := &incidentio.CreateIncidentOptions{
    Name:           "Service Outage",
    IncidentTypeID: "type-id",
    IncidentRoleAssignments: []incidentio.CreateRoleAssignment{
        {
            IncidentRoleID: roles[0].ID, // e.g., "Incident Lead"
            UserID:         users[0].ID,
        },
    },
}
```

### Working with Custom Fields

```go
// List custom fields to see what's available
customFields, _, _ := client.CustomFields.List(ctx)

// Set custom field values when creating an incident
opts := &incidentio.CreateIncidentOptions{
    Name:           "Performance Degradation",
    IncidentTypeID: "type-id",
    CustomFieldValues: map[string]interface{}{
        customFields[0].ID: "high-priority",
        customFields[1].ID: "customer-facing",
    },
}
```

### Error Handling

The client provides detailed error information:

```go
incident, resp, err := client.Incidents.Get(ctx, "invalid-id")
if err != nil {
    if errResp, ok := err.(*incidentio.ErrorResponse); ok {
        fmt.Printf("API Error: %s (Status: %d)\n", errResp.Detail, errResp.Status)
        fmt.Printf("Request: %s %s\n", errResp.Response.Request.Method,
            errResp.Response.Request.URL)
    }
    return
}
```

### Pagination

For endpoints that support pagination:

```go
opts := &incidentio.ListOptions{
    PageSize: 25,
    After:    "01FCNDV6P870EA6S7TK1DSYDG0", // cursor from previous response
}

incidents, resp, err := client.Incidents.List(ctx, opts)
// Check resp.Header for pagination info
```

## Available Services

The client provides access to the following Incident.io API resources:

- **Incidents** - Create, read, update, and delete incidents
- **Severities** - List available severity levels
- **IncidentTypes** - List available incident types
- **IncidentRoles** - List available incident roles
- **CustomFields** - List custom fields configured for your organization
- **Users** - List users in your organization
- **Actions** - Manage incident actions (coming soon)
- **Workflows** - Manage workflows (coming soon)
- **Schedules** - Manage on-call schedules (coming soon)
- **Webhooks** - Manage webhook endpoints (coming soon)

## API Coverage

This client currently implements the core functionality of the Incident.io API. The following endpoints are fully supported:

- ✅ Incidents (Create, List, Get, Update, Delete)
- ✅ Severities (List)
- ✅ Incident Types (List)
- ✅ Incident Roles (List)
- ✅ Custom Fields (List)
- ✅ Users (List)
- 🚧 Actions (Coming soon)
- 🚧 Workflows (Coming soon)
- 🚧 Schedules (Coming soon)
- 🚧 Webhooks (Coming soon)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Running Tests

```bash
go test -v ./...
```

## License

This library is distributed under the Apache License 2.0. See the [LICENSE](LICENSE) file for details.

## Support

For issues related to this client library, please open an issue on GitHub.

For questions about the Incident.io API itself, refer to the [official API documentation](https://api-docs.incident.io/) or contact Incident.io support.
