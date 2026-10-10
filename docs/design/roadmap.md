# Roadmap: remaining API coverage and improvements

Status as of v0.1.0. Specs live at `https://docs.incident.io/openapi/tags/<tag>.json`
(index: `https://docs.incident.io/llms.txt`).

## Conventions for new services

- Build on `Client.send` and `getKey[T]` in `incidentio/request.go`.
- Wire the service in the `Client` struct and `NewClient` (`incidentio/incidentio.go`).
- Add every method to the `apiCalls` map in `incidentio/client_test.go` (error-path tests).
  Methods that decode no body (`Delete`, `BulkUpdate`) are skipped in the invalid-JSON test.
- Use fixture-based tests with `setup()`, `testMethod`, `testQuery`, `testBody`.
- Expose deeply nested, rarely used structures as `json.RawMessage` (see `alert_sources.go`).
- Update the README "Available Services" and "API Coverage" lists and `CHANGELOG.md`.
- One PR per area; run `gofmt`, `go vet`, `go test -race -cover ./...`, `golangci-lint run`.

## Implemented (v0.1.0)

incidents-v2, incident-updates-v2, incident-timeline-items-v2, incident-participants-v2,
incident-roles-v2, incident-types-v1, severities-v1, custom-fields-v2, actions-v3,
workflows-v2, schedules-v2, schedule-overrides-v2, schedule-entries-v2 (via schedules),
users-v2 (List, Get), alerts-v2, alert-sources-v2, catalog-types-v3, catalog-entries-v3.

## Next up (suggested order)

1. **Escalations and on-call**: escalations-v2, escalation-paths-v2, escalation-path-templates-v2.
2. **Incident statuses**: incident-statuses-v1 CRUD (`UpdateIncidentOptions` already uses status IDs).
3. **Pagination cursor**: list methods return only items and `*http.Response`, so callers cannot
   read `pagination_meta.after`. Options: a result wrapper, `ListAll` iterators, or a cursor accessor.
   Needs a decision before changing the public API.
4. **Teams, follow-ups, timestamps**: teams-v3, follow-ups-v3 (v2 is legacy), incident-timestamps-v2.

## Backlog by area

### Incidents
- incident-attachments-v1, incident-relationships-v1, incident-memberships-v1,
  incident-team-memberships-v1, incident-alerts-v2, incident-activity-log-entries-v2,
  incident-forms-v3, incident-templates-v1, incident-participant-workloads-v2,
  postmortemdocuments-v1 (beyond import), custom-field-options-v1.

### Alerts
- Alert sources: `actions/validate` endpoint.
- alert-attributes-v2, alert-source-attributes-v3, alert-events-v2, alert-notes-v1,
  alert-tags-v2, alert-routes-v2/v3, announcement-rules-v2, announcement-templates-v2.
- Alerts list: the `attributes` filter is not implemented.

### On-call and schedules
- schedule-replicas-v2, schedule-sync-rules-v2, schedule-sync-targets-v2.
- Schedules: `preview_entries`.
- notification-methods-v2, notification-rules-v2, on-call-notification-pauses-v2.
- Pay: pay-configs-v2, pay-config-one-off-rules-v2, pay-config-weekly-rules-v2, pay-reports-v2.
- Users: `paging_provider` endpoints (GET/POST `/v2/users/{id}/paging_provider`).

### Catalog
- catalog-resources-v3.
- Entries: `expand` query on `Get`; return the accompanying catalog type from `Get`/`Update`
  (currently dropped).

### Calls
- call-routes-v2, call-route-options-v2, call-route-allowed-callers-v2, call-sessions-v2,
  call-transcript-entries-v2.

### Status pages
- status-pages-v1/v2, status-page-components-v2, status-page-component-availability-v2,
  status-page-structures-v2, status-page-incidents-v2, status-page-incident-updates-v2,
  status-page-maintenances-v2, status-page-maintenance-updates-v2,
  status-page-response-incidents-v1.

### Platform and admin
- api-keys-v1, ipallowlists-v1, maintenancewindows-v1, policies-v2, policy-findings-v2,
  secrets-v2, heartbeat-v2, telemetry-v2, utilities-v1, workflow-runs-v2.
- Legacy versions to skip unless requested: actions-v2, follow-ups-v2.

## Client and quality improvements

- **Webhooks**: `webhooks.go` is a placeholder (the spec only documents webhook events).
  Decide between removing it, or adding signature verification and typed event payloads.
- **Retries and rate limiting**: honour `429` and `Retry-After`; optional retry policy.
- **Pagination helpers**: see item 3 above.
- **Examples**: runnable `example_test.go` snippets for common flows.
- **Release process**: releases use signed tags via gitsign; `tag.gpgsign` is off, so use
  `git tag -s`. Consider automating release notes from `CHANGELOG.md`.
