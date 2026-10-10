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

    incidentio "github.com/cpanato/go-incident-io/incidentio"
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
    Mode:           "standard", // or "test", "retrospective", "tutorial"
    Visibility:     "public",   // or "private"
    // IdempotencyKey is generated for you when empty. Set it to make retries safe.
}

incident, _, err := client.Incidents.Create(ctx, opts)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Created incident: %s\n", incident.ID)
```

### Updating an Incident

```go
statusID := "01FCNDV6P870EA6S7TK1DSYDG0" // from the incident statuses of your organisation
summary := "Issue has been resolved by restarting the database connection pool"
updateOpts := &incidentio.UpdateIncidentOptions{
    IncidentStatusID:      &statusID,
    Summary:               &summary,
    NotifyIncidentChannel: true,
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
            Assignee:       incidentio.UserReference{ID: users[0].ID},
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
    CustomFieldEntries: []incidentio.CustomFieldEntryPayload{
        {
            CustomFieldID: customFields[0].ID,
            Values:        []incidentio.CustomFieldValuePayload{{ValueText: "high-priority"}},
        },
    },
}

// Manage custom fields
field, _, _ := client.CustomFields.Create(ctx, &incidentio.CreateCustomFieldOptions{
    Name:        "Customer impact",
    Description: "How customers are affected",
    FieldType:   "single_select",
})
```

### Filtering Incidents

```go
incidents, _, err := client.Incidents.List(ctx, &incidentio.IncidentListOptions{
    ListOptions:    incidentio.ListOptions{PageSize: 50},
    SortBy:         "created_at_oldest_first",
    StatusCategory: incidentio.Filter{"one_of": {"active"}},
    CreatedAt:      incidentio.Filter{"gte": {"2024-01-01"}},
})
```

### Schedules and Overrides

```go
entries, _, _ := client.Schedules.ListEntries(ctx, "schedule-id", &incidentio.ScheduleEntriesOptions{
    EntryWindow: &incidentio.TimeWindow{StartAt: start, EndAt: end},
})
fmt.Println(entries.Final) // effective schedule with overrides applied

override, _, _ := client.Schedules.CreateOverride(ctx, &incidentio.CreateOverrideOptions{
    ScheduleID: "schedule-id",
    RotationID: "rotation-id",
    LayerID:    "layer-id",
    User:       incidentio.UserReference{Email: "jane@example.com"},
    StartAt:    incidentio.Timestamp{Time: start},
    EndAt:      incidentio.Timestamp{Time: end},
})
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

List options embed `ListOptions` (`PageSize` and `After`). The cursor is the ID of
the last item you received:

```go
opts := &incidentio.IncidentListOptions{
    ListOptions: incidentio.ListOptions{
        PageSize: 25,
        After:    "01FCNDV6P870EA6S7TK1DSYDG0", // ID of the last incident of the previous page
    },
}

incidents, _, err := client.Incidents.List(ctx, opts)
```

## Available Services

The client provides access to the following Incident.io API resources:

- **Incidents** - Create, list (with filters), get, edit and import postmortem documents
- **Severities** - List available severity levels
- **IncidentTypes** - List available incident types
- **IncidentRoles** - List available incident roles
- **CustomFields** - Create, list, get, update and delete custom fields
- **Users** - List and get users in your organization
- **IncidentUpdates** - List and post incident updates
- **IncidentTimelineItems** - List and create timeline items
- **IncidentParticipants** - List incident participants
- **Actions** - Create, list, get, update and delete incident actions
- **Workflows** - Create, list, get, update and delete workflows
- **Schedules** - Manage on-call schedules, entries and overrides
- **Alerts** - List (with filters), get, resolve and manage tags
- **AlertSources** - Create, list, get, update and delete alert sources
- **Webhooks** - Placeholder (the API only documents webhook events)

## API Coverage

This client currently implements the core functionality of the Incident.io API. The following endpoints are fully supported:

- ✅ Incidents (Create, List, Get, Edit, Import postmortem document)
- ✅ Severities (List)
- ✅ Incident Types (List)
- ✅ Incident Roles (List)
- ✅ Custom Fields (Create, List, Get, Update, Delete)
- ✅ Users (List, Get)
- ✅ Incident Updates (List, Create)
- ✅ Incident Timeline Items (List, Create)
- ✅ Incident Participants (List)
- ✅ Actions v3 (Create, List, Get, Update, Delete)
- ✅ Workflows (Create, List, Get, Update, Delete)
- ✅ Schedules (CRUD, entries, overrides)
- ✅ Alerts (List, Get, Resolve, AddTags, RemoveTags, SetTags)
- ✅ Alert Sources (Create, List, Get, Update, Delete)
- 🚧 Webhooks (no management API)

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
