package incidentio

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// SchedulesService handles communication with the schedules related methods.
type SchedulesService struct {
	client *Client
}

// Schedule represents an on-call schedule in Incident.io.
type Schedule struct {
	ID            string          `json:"id"`
	Name          string          `json:"name"`
	Timezone      string          `json:"timezone"`
	Rotations     []Rotation      `json:"rotations"`
	CurrentShifts []CurrentShift  `json:"current_shifts,omitempty"`
	Config        *ScheduleConfig `json:"config"`
	CreatedAt     Timestamp       `json:"created_at"`
	UpdatedAt     Timestamp       `json:"updated_at"`
}

// Rotation represents a rotation within a schedule.
type Rotation struct {
	ID              string            `json:"id"`
	Name            string            `json:"name"`
	Layers          []Layer           `json:"layers"`
	EffectiveFrom   *Timestamp        `json:"effective_from,omitempty"`
	HandoverStartAt Timestamp         `json:"handover_start_at"`
	HandoversAt     Timestamp         `json:"handovers_at"`
	WorkingInterval []WorkingInterval `json:"working_interval,omitempty"`
}

// Layer represents a layer within a rotation.
type Layer struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// WorkingInterval represents working hours configuration.
type WorkingInterval struct {
	StartTime string   `json:"start_time"`
	EndTime   string   `json:"end_time"`
	Weekdays  []string `json:"weekdays"`
}

// CurrentShift represents the current active shift.
type CurrentShift struct {
	UserID     string    `json:"user_id"`
	User       *User     `json:"user,omitempty"`
	RotationID string    `json:"rotation_id"`
	LayerID    string    `json:"layer_id"`
	StartAt    Timestamp `json:"start_at"`
	EndAt      Timestamp `json:"end_at"`
	FinalShift bool      `json:"final_shift"`
}

// ScheduleConfig represents schedule configuration.
type ScheduleConfig struct {
	Rotations []RotationConfig `json:"rotations"`
}

// RotationConfig represents rotation configuration within a schedule.
type RotationConfig struct {
	ID              string                  `json:"id"`
	Name            string                  `json:"name"`
	Layers          []LayerConfig           `json:"layers"`
	EffectiveFrom   *Timestamp              `json:"effective_from,omitempty"`
	HandoverStartAt Timestamp               `json:"handover_start_at"`
	HandoversAt     Timestamp               `json:"handovers_at"`
	WorkingInterval []WorkingIntervalConfig `json:"working_interval,omitempty"`
}

// LayerConfig represents layer configuration within a rotation.
type LayerConfig struct {
	ID    string      `json:"id"`
	Name  string      `json:"name"`
	Users []LayerUser `json:"users"`
}

// LayerUser represents a user in a layer.
type LayerUser struct {
	UserID string `json:"user_id"`
}

// WorkingIntervalConfig represents working interval configuration.
type WorkingIntervalConfig struct {
	StartTime string   `json:"start_time"`
	EndTime   string   `json:"end_time"`
	Weekdays  []string `json:"weekdays"`
}

// ScheduleEntry represents an entry in a schedule.
type ScheduleEntry struct {
	EntryID     string    `json:"entry_id"`
	Fingerprint string    `json:"fingerprint"`
	RotationID  string    `json:"rotation_id"`
	StartAt     Timestamp `json:"start_at"`
	EndAt       Timestamp `json:"end_at"`
	User        *User     `json:"user,omitempty"`
}

// ScheduleEntries groups the entries of a schedule for a window of time.
// Final is the effective schedule once overrides have been applied.
type ScheduleEntries struct {
	Scheduled []*ScheduleEntry `json:"scheduled"`
	Overrides []*ScheduleEntry `json:"overrides"`
	Final     []*ScheduleEntry `json:"final"`
}

// Override represents a schedule override.
type Override struct {
	ID         string    `json:"id"`
	ScheduleID string    `json:"schedule_id"`
	RotationID string    `json:"rotation_id"`
	LayerID    string    `json:"layer_id"`
	User       *User     `json:"user,omitempty"`
	StartAt    Timestamp `json:"start_at"`
	EndAt      Timestamp `json:"end_at"`
	CreatedAt  Timestamp `json:"created_at"`
	UpdatedAt  Timestamp `json:"updated_at"`
}

// UserReference identifies a user by ID, email or Slack user ID.
type UserReference struct {
	ID          string `json:"id,omitempty"`
	Email       string `json:"email,omitempty"`
	SlackUserID string `json:"slack_user_id,omitempty"`
}

// ScheduleListOptions represents options for listing schedules.
type ScheduleListOptions struct {
	ListOptions
}

// ScheduleEntriesOptions represents options for listing schedule entries.
type ScheduleEntriesOptions struct {
	EntryWindow *TimeWindow
}

// TimeWindow represents a time window for schedule entries.
type TimeWindow struct {
	StartAt time.Time
	EndAt   time.Time
}

// ListOverridesOptions represents options for listing schedule overrides.
type ListOverridesOptions struct {
	ListOptions
	RotationID string
	LayerID    string
}

// CreateScheduleOptions represents options for creating a schedule.
type CreateScheduleOptions struct {
	Name     string          `json:"name"`
	Timezone string          `json:"timezone"`
	Config   *ScheduleConfig `json:"config"`
}

// UpdateScheduleOptions represents options for updating a schedule.
type UpdateScheduleOptions struct {
	Name     *string         `json:"name,omitempty"`
	Timezone *string         `json:"timezone,omitempty"`
	Config   *ScheduleConfig `json:"config,omitempty"`
}

// CreateOverrideOptions represents options for creating an override.
type CreateOverrideOptions struct {
	ScheduleID string        `json:"schedule_id"`
	RotationID string        `json:"rotation_id"`
	LayerID    string        `json:"layer_id"`
	User       UserReference `json:"user"`
	StartAt    Timestamp     `json:"start_at"`
	EndAt      Timestamp     `json:"end_at"`
}

// UpdateOverrideOptions represents options for updating an override.
type UpdateOverrideOptions struct {
	RotationID string        `json:"rotation_id"`
	LayerID    string        `json:"layer_id"`
	User       UserReference `json:"user"`
	StartAt    Timestamp     `json:"start_at"`
	EndAt      Timestamp     `json:"end_at"`
}

// List returns a list of schedules.
func (s *SchedulesService) List(ctx context.Context, opts *ScheduleListOptions) ([]*Schedule, *http.Response, error) {
	q := url.Values{}
	if opts != nil {
		opts.apply(q)
	}

	return getKey[[]*Schedule](ctx, s.client, "GET", "v2/schedules", q, nil, "schedules")
}

// Get returns a single schedule.
func (s *SchedulesService) Get(ctx context.Context, id string) (*Schedule, *http.Response, error) {
	return getKey[*Schedule](ctx, s.client, "GET", fmt.Sprintf("v2/schedules/%s", id), nil, nil, "schedule")
}

// Create creates a new schedule.
func (s *SchedulesService) Create(ctx context.Context, opts *CreateScheduleOptions) (*Schedule, *http.Response, error) {
	return getKey[*Schedule](ctx, s.client, "POST", "v2/schedules", nil, opts, "schedule")
}

// Update updates a schedule.
func (s *SchedulesService) Update(ctx context.Context, id string, opts *UpdateScheduleOptions) (*Schedule, *http.Response, error) {
	return getKey[*Schedule](ctx, s.client, "PUT", fmt.Sprintf("v2/schedules/%s", id), nil, opts, "schedule")
}

// Delete deletes a schedule.
func (s *SchedulesService) Delete(ctx context.Context, id string) (*http.Response, error) {
	return s.client.send(ctx, "DELETE", fmt.Sprintf("v2/schedules/%s", id), nil, nil, nil)
}

// ListEntries returns the entries of a schedule, grouped as scheduled, overrides and final.
func (s *SchedulesService) ListEntries(ctx context.Context, scheduleID string, opts *ScheduleEntriesOptions) (*ScheduleEntries, *http.Response, error) {
	q := url.Values{"schedule_id": {scheduleID}}
	if opts != nil && opts.EntryWindow != nil {
		q.Set("entry_window_start", opts.EntryWindow.StartAt.Format(time.RFC3339))
		q.Set("entry_window_end", opts.EntryWindow.EndAt.Format(time.RFC3339))
	}

	return getKey[*ScheduleEntries](ctx, s.client, "GET", "v2/schedule_entries", q, nil, "schedule_entries")
}

// ListOverrides returns overrides for a schedule.
func (s *SchedulesService) ListOverrides(ctx context.Context, scheduleID string, opts *ListOverridesOptions) ([]*Override, *http.Response, error) {
	q := url.Values{"schedule_id": {scheduleID}}
	if opts != nil {
		opts.apply(q)
		if opts.RotationID != "" {
			q.Set("rotation_id", opts.RotationID)
		}
		if opts.LayerID != "" {
			q.Set("layer_id", opts.LayerID)
		}
	}

	return getKey[[]*Override](ctx, s.client, "GET", "v2/schedule_overrides", q, nil, "overrides")
}

// GetOverride returns a single override.
func (s *SchedulesService) GetOverride(ctx context.Context, overrideID string) (*Override, *http.Response, error) {
	return getKey[*Override](ctx, s.client, "GET", fmt.Sprintf("v2/schedule_overrides/%s", overrideID), nil, nil, "override")
}

// CreateOverride creates a new override for a schedule.
func (s *SchedulesService) CreateOverride(ctx context.Context, opts *CreateOverrideOptions) (*Override, *http.Response, error) {
	return getKey[*Override](ctx, s.client, "POST", "v2/schedule_overrides", nil, opts, "override")
}

// UpdateOverride updates an override.
func (s *SchedulesService) UpdateOverride(ctx context.Context, overrideID string, opts *UpdateOverrideOptions) (*Override, *http.Response, error) {
	return getKey[*Override](ctx, s.client, "PUT", fmt.Sprintf("v2/schedule_overrides/%s", overrideID), nil, opts, "override")
}

// DeleteOverride deletes an override.
func (s *SchedulesService) DeleteOverride(ctx context.Context, overrideID string) (*http.Response, error) {
	return s.client.send(ctx, "DELETE", fmt.Sprintf("v2/schedule_overrides/%s", overrideID), nil, nil, nil)
}
