# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2026-10-10

First tagged release.

### Added

- Client for the Incident.io API with functional options, context support and
  typed `ErrorResponse` errors.
- Incidents: create, list (with filters and pagination), get, edit and import
  postmortem documents (#16).
- Custom fields, actions (v3) and workflows: create, list, get, update and
  delete (#16).
- Schedules: CRUD, entries and overrides (#16).
- Severities, incident types, incident roles and users (list) (#16).
- Incident updates (list, create), incident timeline items (list, create),
  incident participants (list) and `Users.Get` (#19).
- Alerts (list with filters, get, resolve, add/remove/set tags) and alert
  sources (CRUD) (#20).
- Catalog types (CRUD and update schema) and catalog entries (CRUD and bulk
  update) (#21).
- Idempotency keys are generated automatically for create calls that require
  them.

### Changed

- Share request and decode helpers across all services, with no change to the
  public API (#18).
- Require Go 1.27.2 and update golangci-lint to v2.14 (#15).
- Raise test coverage to over 97% (#17).

[0.1.0]: https://github.com/cpanato/go-incident-io/releases/tag/v0.1.0
