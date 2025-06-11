package incidentio

import (
	"context"
	"fmt"
	"net/http"
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
	UserID         string    `json:"user_id"`
	User           *User     `json:"user,omitempty"`
	ScheduleID     string    `json:"schedule_id"`
	RotationID     string    `json:"rotation_id"`
	LayerID        string    `json:"layer_id"`
	Interval       *Interval `json:"interval"`
	IsOverride     bool      `json:"is_override"`
	FinalShift     bool      `json:"final_shift"`
	OverriddenByID string    `json:"overridden_by_id,omitempty"`
}

// Interval represents a time interval.
type Interval struct {
	StartAt Timestamp `json:"start_at"`
	EndAt   Timestamp `json:"end_at"`
}

// Override represents a schedule override.
type Override struct {
	ID         string    `json:"id"`
	ScheduleID string    `json:"schedule_id"`
	UserID     string    `json:"user_id"`
	User       *User     `json:"user,omitempty"`
	StartAt    Timestamp `json:"start_at"`
	EndAt      Timestamp `json:"end_at"`
	CreatedAt  Timestamp `json:"created_at"`
	UpdatedAt  Timestamp `json:"updated_at"`
}

// ScheduleListOptions represents options for listing schedules.
type ScheduleListOptions struct {
	ListOptions
}

// ScheduleEntriesOptions represents options for listing schedule entries.
type ScheduleEntriesOptions struct {
	EntryWindow *TimeWindow `url:"-"`
	ListOptions
}

// TimeWindow represents a time window for schedule entries.
type TimeWindow struct {
	StartAt time.Time
	EndAt   time.Time
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
	UserID  string    `json:"user_id"`
	StartAt Timestamp `json:"start_at"`
	EndAt   Timestamp `json:"end_at"`
}

// UpdateOverrideOptions represents options for updating an override.
type UpdateOverrideOptions struct {
	UserID  *string    `json:"user_id,omitempty"`
	StartAt *Timestamp `json:"start_at,omitempty"`
	EndAt   *Timestamp `json:"end_at,omitempty"`
}

// List returns a list of schedules.
func (s *SchedulesService) List(ctx context.Context, _ *ScheduleListOptions) ([]*Schedule, *http.Response, error) {
	u := "v2/schedules"

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var result struct {
		Schedules []*Schedule `json:"schedules"`
	}
	resp, err := s.client.Do(ctx, req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result.Schedules, resp, nil
}

// Get returns a single schedule.
func (s *SchedulesService) Get(ctx context.Context, id string) (*Schedule, *http.Response, error) {
	u := fmt.Sprintf("v2/schedules/%s", id)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var result struct {
		Schedule *Schedule `json:"schedule"`
	}
	resp, err := s.client.Do(ctx, req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result.Schedule, resp, nil
}

// Create creates a new schedule.
func (s *SchedulesService) Create(ctx context.Context, opts *CreateScheduleOptions) (*Schedule, *http.Response, error) {
	u := "v2/schedules"

	req, err := s.client.NewRequest("POST", u, opts)
	if err != nil {
		return nil, nil, err
	}

	var result struct {
		Schedule *Schedule `json:"schedule"`
	}
	resp, err := s.client.Do(ctx, req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result.Schedule, resp, nil
}

// Update updates a schedule.
func (s *SchedulesService) Update(ctx context.Context, id string, opts *UpdateScheduleOptions) (*Schedule, *http.Response, error) {
	u := fmt.Sprintf("v2/schedules/%s", id)

	req, err := s.client.NewRequest("PUT", u, opts)
	if err != nil {
		return nil, nil, err
	}

	var result struct {
		Schedule *Schedule `json:"schedule"`
	}
	resp, err := s.client.Do(ctx, req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result.Schedule, resp, nil
}

// Delete deletes a schedule.
func (s *SchedulesService) Delete(ctx context.Context, id string) (*http.Response, error) {
	u := fmt.Sprintf("v2/schedules/%s", id)

	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// ListEntries returns entries for a schedule.
func (s *SchedulesService) ListEntries(ctx context.Context, scheduleID string, opts *ScheduleEntriesOptions) ([]*ScheduleEntry, *http.Response, error) {
	u := fmt.Sprintf("v2/schedules/%s/entries", scheduleID)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	// Add time window parameters if provided
	if opts != nil && opts.EntryWindow != nil {
		q := req.URL.Query()
		q.Add("entry_window[start_at]", opts.EntryWindow.StartAt.Format(time.RFC3339))
		q.Add("entry_window[end_at]", opts.EntryWindow.EndAt.Format(time.RFC3339))
		req.URL.RawQuery = q.Encode()
	}

	var result struct {
		ScheduleEntries []*ScheduleEntry `json:"schedule_entries"`
	}
	resp, err := s.client.Do(ctx, req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result.ScheduleEntries, resp, nil
}

// ListOverrides returns overrides for a schedule.
func (s *SchedulesService) ListOverrides(ctx context.Context, scheduleID string, _ *ListOptions) ([]*Override, *http.Response, error) {
	u := fmt.Sprintf("v2/schedules/%s/overrides", scheduleID)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var result struct {
		Overrides []*Override `json:"overrides"`
	}
	resp, err := s.client.Do(ctx, req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result.Overrides, resp, nil
}

// GetOverride returns a single override.
func (s *SchedulesService) GetOverride(ctx context.Context, scheduleID, overrideID string) (*Override, *http.Response, error) {
	u := fmt.Sprintf("v2/schedules/%s/overrides/%s", scheduleID, overrideID)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var result struct {
		Override *Override `json:"override"`
	}
	resp, err := s.client.Do(ctx, req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result.Override, resp, nil
}

// CreateOverride creates a new override for a schedule.
func (s *SchedulesService) CreateOverride(ctx context.Context, scheduleID string, opts *CreateOverrideOptions) (*Override, *http.Response, error) {
	u := fmt.Sprintf("v2/schedules/%s/overrides", scheduleID)

	req, err := s.client.NewRequest("POST", u, opts)
	if err != nil {
		return nil, nil, err
	}

	var result struct {
		Override *Override `json:"override"`
	}
	resp, err := s.client.Do(ctx, req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result.Override, resp, nil
}

// UpdateOverride updates an override.
func (s *SchedulesService) UpdateOverride(ctx context.Context, scheduleID, overrideID string, opts *UpdateOverrideOptions) (*Override, *http.Response, error) {
	u := fmt.Sprintf("v2/schedules/%s/overrides/%s", scheduleID, overrideID)

	req, err := s.client.NewRequest("PUT", u, opts)
	if err != nil {
		return nil, nil, err
	}

	var result struct {
		Override *Override `json:"override"`
	}
	resp, err := s.client.Do(ctx, req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result.Override, resp, nil
}

// DeleteOverride deletes an override.
func (s *SchedulesService) DeleteOverride(ctx context.Context, scheduleID, overrideID string) (*http.Response, error) {
	u := fmt.Sprintf("v2/schedules/%s/overrides/%s", scheduleID, overrideID)

	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
