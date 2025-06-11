package incidentio

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestSchedulesService_List(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/schedules", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Authorization", "Bearer test-key")

		response := `{
			"schedules": [
				{
					"id": "01G0J1EXE7AXZ2C93K61WBPYEH",
					"name": "Engineering On-Call",
					"timezone": "America/New_York",
					"rotations": [
						{
							"id": "01G0J1EXE7AXZ2C93K61WBPYEH",
							"name": "Weekly Rotation",
							"layers": [
								{
									"id": "01G0J1EXE7AXZ2C93K61WBPYEH",
									"name": "Primary On-Call"
								}
							],
							"handover_start_at": "2021-08-17T13:28:57.801578Z",
							"handovers_at": "2021-08-24T13:28:57.801578Z"
						}
					],
					"current_shifts": [
						{
							"user_id": "01FCNDV6P870EA6S7TK1DSYDG0",
							"rotation_id": "01G0J1EXE7AXZ2C93K61WBPYEH",
							"layer_id": "01G0J1EXE7AXZ2C93K61WBPYEH",
							"start_at": "2021-08-17T13:28:57.801578Z",
							"end_at": "2021-08-24T13:28:57.801578Z",
							"final_shift": false
						}
					],
					"created_at": "2021-08-17T13:28:57.801578Z",
					"updated_at": "2021-08-17T13:28:57.801578Z"
				}
			]
		}`

		_, _ = fmt.Fprint(w, response)
	})

	ctx := context.Background()
	schedules, _, err := client.Schedules.List(ctx, nil)
	if err != nil {
		t.Errorf("Schedules.List returned error: %v", err)
	}

	expected := []*Schedule{
		{
			ID:       "01G0J1EXE7AXZ2C93K61WBPYEH",
			Name:     "Engineering On-Call",
			Timezone: "America/New_York",
			Rotations: []Rotation{
				{
					ID:   "01G0J1EXE7AXZ2C93K61WBPYEH",
					Name: "Weekly Rotation",
					Layers: []Layer{
						{
							ID:   "01G0J1EXE7AXZ2C93K61WBPYEH",
							Name: "Primary On-Call",
						},
					},
					HandoverStartAt: Timestamp{parseTime("2021-08-17T13:28:57.801578Z")},
					HandoversAt:     Timestamp{parseTime("2021-08-24T13:28:57.801578Z")},
				},
			},
			CurrentShifts: []CurrentShift{
				{
					UserID:     "01FCNDV6P870EA6S7TK1DSYDG0",
					RotationID: "01G0J1EXE7AXZ2C93K61WBPYEH",
					LayerID:    "01G0J1EXE7AXZ2C93K61WBPYEH",
					StartAt:    Timestamp{parseTime("2021-08-17T13:28:57.801578Z")},
					EndAt:      Timestamp{parseTime("2021-08-24T13:28:57.801578Z")},
					FinalShift: false,
				},
			},
			CreatedAt: Timestamp{parseTime("2021-08-17T13:28:57.801578Z")},
			UpdatedAt: Timestamp{parseTime("2021-08-17T13:28:57.801578Z")},
		},
	}

	if !reflect.DeepEqual(schedules, expected) {
		t.Errorf("Schedules.List returned %+v, want %+v", schedules, expected)
	}
}

func TestSchedulesService_Get(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	scheduleID := "01G0J1EXE7AXZ2C93K61WBPYEH"

	mux.HandleFunc(fmt.Sprintf("/v2/schedules/%s", scheduleID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		response := `{
			"schedule": {
				"id": "01G0J1EXE7AXZ2C93K61WBPYEH",
				"name": "Engineering On-Call",
				"timezone": "America/New_York",
				"config": {
					"rotations": [
						{
							"id": "01G0J1EXE7AXZ2C93K61WBPYEH",
							"name": "Weekly Rotation",
							"layers": [
								{
									"id": "01G0J1EXE7AXZ2C93K61WBPYEH",
									"name": "Primary On-Call",
									"users": [
										{"user_id": "01FCNDV6P870EA6S7TK1DSYDG0"},
										{"user_id": "01FCQSP07Z74QMMYPDDGQB9FTG"}
									]
								}
							],
							"handover_start_at": "2021-08-17T13:28:57.801578Z",
							"handovers_at": "2021-08-24T13:28:57.801578Z",
							"working_interval": [
								{
									"start_time": "09:00",
									"end_time": "17:00",
									"weekdays": ["monday", "tuesday", "wednesday", "thursday", "friday"]
								}
							]
						}
					]
				},
				"created_at": "2021-08-17T13:28:57.801578Z",
				"updated_at": "2021-08-17T13:28:57.801578Z"
			}
		}`

		_, _ = fmt.Fprint(w, response)
	})

	ctx := context.Background()
	schedule, _, err := client.Schedules.Get(ctx, scheduleID)
	if err != nil {
		t.Errorf("Schedules.Get returned error: %v", err)
	}

	if schedule.ID != scheduleID {
		t.Errorf("Schedules.Get returned ID %s, want %s", schedule.ID, scheduleID)
	}

	if schedule.Config == nil {
		t.Error("Schedules.Get returned nil Config")
	} else if len(schedule.Config.Rotations) != 1 {
		t.Errorf("Schedules.Get returned %d rotations, want 1", len(schedule.Config.Rotations))
	}
}

func TestSchedulesService_Create(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &CreateScheduleOptions{
		Name:     "New On-Call Schedule",
		Timezone: "Europe/London",
		Config: &ScheduleConfig{
			Rotations: []RotationConfig{
				{
					ID:   "weekly-rotation",
					Name: "Weekly Rotation",
					Layers: []LayerConfig{
						{
							ID:   "layer-1",
							Name: "Primary",
							Users: []LayerUser{
								{UserID: "01FCNDV6P870EA6S7TK1DSYDG0"},
							},
						},
					},
					HandoverStartAt: Timestamp{parseTime("2021-08-17T09:00:00Z")},
					HandoversAt:     Timestamp{parseTime("2021-08-24T09:00:00Z")},
				},
			},
		},
	}

	mux.HandleFunc("/v2/schedules", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Content-Type", "application/json")

		var received CreateScheduleOptions
		err := json.NewDecoder(r.Body).Decode(&received)
		if err != nil {
			t.Errorf("error decoding request body: %v", err)
		}

		// Compare JSON representation to handle time formatting
		expectedJSON, _ := json.Marshal(input)
		receivedJSON, _ := json.Marshal(received)
		if string(expectedJSON) != string(receivedJSON) {
			t.Errorf("Request body = %s, want %s", string(receivedJSON), string(expectedJSON))
		}

		response := `{
			"schedule": {
				"id": "01G0J1EXE7AXZ2C93K61WBPYEH",
				"name": "New On-Call Schedule",
				"timezone": "Europe/London",
				"created_at": "2021-08-17T13:28:57.801578Z",
				"updated_at": "2021-08-17T13:28:57.801578Z"
			}
		}`

		w.WriteHeader(http.StatusCreated)
		_, _ = fmt.Fprint(w, response)
	})

	ctx := context.Background()
	schedule, resp, err := client.Schedules.Create(ctx, input)
	if err != nil {
		t.Errorf("Schedules.Create returned error: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Schedules.Create returned status %d, want %d", resp.StatusCode, http.StatusCreated)
	}

	if schedule.Name != "New On-Call Schedule" {
		t.Errorf("Schedules.Create returned name %s, want %s", schedule.Name, "New On-Call Schedule")
	}
}

func TestSchedulesService_Update(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	scheduleID := "01G0J1EXE7AXZ2C93K61WBPYEH"
	name := "Updated Schedule Name"
	timezone := "Asia/Tokyo"

	input := &UpdateScheduleOptions{
		Name:     &name,
		Timezone: &timezone,
	}

	mux.HandleFunc(fmt.Sprintf("/v2/schedules/%s", scheduleID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")

		var received UpdateScheduleOptions
		err := json.NewDecoder(r.Body).Decode(&received)
		if err != nil {
			t.Errorf("error decoding request body: %v", err)
		}

		if !reflect.DeepEqual(received, *input) {
			t.Errorf("Request body = %+v, want %+v", received, *input)
		}

		response := `{
			"schedule": {
				"id": "01G0J1EXE7AXZ2C93K61WBPYEH",
				"name": "Updated Schedule Name",
				"timezone": "Asia/Tokyo",
				"created_at": "2021-08-17T13:28:57.801578Z",
				"updated_at": "2021-08-17T14:28:57.801578Z"
			}
		}`

		_, _ = fmt.Fprint(w, response)
	})

	ctx := context.Background()
	schedule, _, err := client.Schedules.Update(ctx, scheduleID, input)
	if err != nil {
		t.Errorf("Schedules.Update returned error: %v", err)
	}

	if schedule.Name != "Updated Schedule Name" {
		t.Errorf("Schedules.Update returned name %s, want %s", schedule.Name, "Updated Schedule Name")
	}

	if schedule.Timezone != "Asia/Tokyo" {
		t.Errorf("Schedules.Update returned timezone %s, want %s", schedule.Timezone, "Asia/Tokyo")
	}
}

func TestSchedulesService_Delete(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	scheduleID := "01G0J1EXE7AXZ2C93K61WBPYEH"

	mux.HandleFunc(fmt.Sprintf("/v2/schedules/%s", scheduleID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	resp, err := client.Schedules.Delete(ctx, scheduleID)
	if err != nil {
		t.Errorf("Schedules.Delete returned error: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Schedules.Delete returned status %d, want %d", resp.StatusCode, http.StatusNoContent)
	}
}

func TestSchedulesService_ListEntries(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	scheduleID := "01G0J1EXE7AXZ2C93K61WBPYEH"
	startAt := time.Date(2021, 8, 17, 0, 0, 0, 0, time.UTC)
	endAt := time.Date(2021, 8, 24, 0, 0, 0, 0, time.UTC)

	mux.HandleFunc(fmt.Sprintf("/v2/schedules/%s/entries", scheduleID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		// Check query parameters
		query := r.URL.Query()
		if got := query.Get("entry_window[start_at]"); got != startAt.Format(time.RFC3339) {
			t.Errorf("start_at query param = %s, want %s", got, startAt.Format(time.RFC3339))
		}
		if got := query.Get("entry_window[end_at]"); got != endAt.Format(time.RFC3339) {
			t.Errorf("end_at query param = %s, want %s", got, endAt.Format(time.RFC3339))
		}

		response := `{
			"schedule_entries": [
				{
					"user_id": "01FCNDV6P870EA6S7TK1DSYDG0",
					"schedule_id": "01G0J1EXE7AXZ2C93K61WBPYEH",
					"rotation_id": "01G0J1EXE7AXZ2C93K61WBPYEH",
					"layer_id": "01G0J1EXE7AXZ2C93K61WBPYEH",
					"interval": {
						"start_at": "2021-08-17T13:00:00Z",
						"end_at": "2021-08-24T13:00:00Z"
					},
					"is_override": false,
					"final_shift": false
				}
			]
		}`

		_, _ = fmt.Fprint(w, response)
	})

	ctx := context.Background()
	opts := &ScheduleEntriesOptions{
		EntryWindow: &TimeWindow{
			StartAt: startAt,
			EndAt:   endAt,
		},
	}

	entries, _, err := client.Schedules.ListEntries(ctx, scheduleID, opts)
	if err != nil {
		t.Errorf("Schedules.ListEntries returned error: %v", err)
	}

	if len(entries) != 1 {
		t.Errorf("Schedules.ListEntries returned %d entries, want 1", len(entries))
	}

	if entries[0].UserID != "01FCNDV6P870EA6S7TK1DSYDG0" {
		t.Errorf("Entry UserID = %s, want %s", entries[0].UserID, "01FCNDV6P870EA6S7TK1DSYDG0")
	}
}

func TestSchedulesService_ListOverrides(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	scheduleID := "01G0J1EXE7AXZ2C93K61WBPYEH"

	mux.HandleFunc(fmt.Sprintf("/v2/schedules/%s/overrides", scheduleID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		response := `{
			"overrides": [
				{
					"id": "01G0J1EXE7AXZ2C93K61WBPYEH",
					"schedule_id": "01G0J1EXE7AXZ2C93K61WBPYEH",
					"user_id": "01FCNDV6P870EA6S7TK1DSYDG0",
					"start_at": "2021-08-20T09:00:00Z",
					"end_at": "2021-08-20T17:00:00Z",
					"created_at": "2021-08-17T13:28:57.801578Z",
					"updated_at": "2021-08-17T13:28:57.801578Z"
				}
			]
		}`

		_, _ = fmt.Fprint(w, response)
	})

	ctx := context.Background()
	overrides, _, err := client.Schedules.ListOverrides(ctx, scheduleID, nil)
	if err != nil {
		t.Errorf("Schedules.ListOverrides returned error: %v", err)
	}

	expected := []*Override{
		{
			ID:         "01G0J1EXE7AXZ2C93K61WBPYEH",
			ScheduleID: "01G0J1EXE7AXZ2C93K61WBPYEH",
			UserID:     "01FCNDV6P870EA6S7TK1DSYDG0",
			StartAt:    Timestamp{parseTime("2021-08-20T09:00:00Z")},
			EndAt:      Timestamp{parseTime("2021-08-20T17:00:00Z")},
			CreatedAt:  Timestamp{parseTime("2021-08-17T13:28:57.801578Z")},
			UpdatedAt:  Timestamp{parseTime("2021-08-17T13:28:57.801578Z")},
		},
	}

	if !reflect.DeepEqual(overrides, expected) {
		t.Errorf("Schedules.ListOverrides returned %+v, want %+v", overrides, expected)
	}
}

func TestSchedulesService_GetOverride(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	scheduleID := "01G0J1EXE7AXZ2C93K61WBPYEH"
	overrideID := "01G0J1EXE7AXZ2C93K61WBPYEH"

	mux.HandleFunc(fmt.Sprintf("/v2/schedules/%s/overrides/%s", scheduleID, overrideID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		response := `{
			"override": {
				"id": "01G0J1EXE7AXZ2C93K61WBPYEH",
				"schedule_id": "01G0J1EXE7AXZ2C93K61WBPYEH",
				"user_id": "01FCNDV6P870EA6S7TK1DSYDG0",
				"user": {
					"id": "01FCNDV6P870EA6S7TK1DSYDG0",
					"name": "John Doe",
					"email": "john@example.com",
					"role": "user"
				},
				"start_at": "2021-08-20T09:00:00Z",
				"end_at": "2021-08-20T17:00:00Z",
				"created_at": "2021-08-17T13:28:57.801578Z",
				"updated_at": "2021-08-17T13:28:57.801578Z"
			}
		}`

		_, _ = fmt.Fprint(w, response)
	})

	ctx := context.Background()
	override, _, err := client.Schedules.GetOverride(ctx, scheduleID, overrideID)
	if err != nil {
		t.Errorf("Schedules.GetOverride returned error: %v", err)
	}

	if override.ID != overrideID {
		t.Errorf("Schedules.GetOverride returned ID %s, want %s", override.ID, overrideID)
	}

	if override.User == nil {
		t.Error("Schedules.GetOverride returned nil User")
	} else if override.User.Email != "john@example.com" {
		t.Errorf("Override User.Email = %s, want %s", override.User.Email, "john@example.com")
	}
}

func TestSchedulesService_CreateOverride(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	scheduleID := "01G0J1EXE7AXZ2C93K61WBPYEH"

	input := &CreateOverrideOptions{
		UserID:  "01FCNDV6P870EA6S7TK1DSYDG0",
		StartAt: Timestamp{parseTime("2021-08-20T09:00:00Z")},
		EndAt:   Timestamp{parseTime("2021-08-20T17:00:00Z")},
	}

	mux.HandleFunc(fmt.Sprintf("/v2/schedules/%s/overrides", scheduleID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Content-Type", "application/json")

		var received CreateOverrideOptions
		err := json.NewDecoder(r.Body).Decode(&received)
		if err != nil {
			t.Errorf("error decoding request body: %v", err)
		}

		// Compare JSON representation to handle time formatting
		expectedJSON, _ := json.Marshal(input)
		receivedJSON, _ := json.Marshal(received)
		if string(expectedJSON) != string(receivedJSON) {
			t.Errorf("Request body = %s, want %s", string(receivedJSON), string(expectedJSON))
		}

		response := `{
			"override": {
				"id": "01G0J1EXE7AXZ2C93K61WBPYEH",
				"schedule_id": "01G0J1EXE7AXZ2C93K61WBPYEH",
				"user_id": "01FCNDV6P870EA6S7TK1DSYDG0",
				"start_at": "2021-08-20T09:00:00Z",
				"end_at": "2021-08-20T17:00:00Z",
				"created_at": "2021-08-17T13:28:57.801578Z",
				"updated_at": "2021-08-17T13:28:57.801578Z"
			}
		}`

		w.WriteHeader(http.StatusCreated)
		_, _ = fmt.Fprint(w, response)
	})

	ctx := context.Background()
	override, resp, err := client.Schedules.CreateOverride(ctx, scheduleID, input)
	if err != nil {
		t.Errorf("Schedules.CreateOverride returned error: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Schedules.CreateOverride returned status %d, want %d", resp.StatusCode, http.StatusCreated)
	}

	if override.UserID != "01FCNDV6P870EA6S7TK1DSYDG0" {
		t.Errorf("Override UserID = %s, want %s", override.UserID, "01FCNDV6P870EA6S7TK1DSYDG0")
	}
}

func TestSchedulesService_UpdateOverride(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	scheduleID := "01G0J1EXE7AXZ2C93K61WBPYEH"
	overrideID := "01G0J1EXE7AXZ2C93K61WBPYEH"
	newUserID := "01FCQSP07Z74QMMYPDDGQB9FTG"

	input := &UpdateOverrideOptions{
		UserID: &newUserID,
	}

	mux.HandleFunc(fmt.Sprintf("/v2/schedules/%s/overrides/%s", scheduleID, overrideID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")

		var received UpdateOverrideOptions
		err := json.NewDecoder(r.Body).Decode(&received)
		if err != nil {
			t.Errorf("error decoding request body: %v", err)
		}

		if !reflect.DeepEqual(received, *input) {
			t.Errorf("Request body = %+v, want %+v", received, *input)
		}

		response := `{
			"override": {
				"id": "01G0J1EXE7AXZ2C93K61WBPYEH",
				"schedule_id": "01G0J1EXE7AXZ2C93K61WBPYEH",
				"user_id": "01FCQSP07Z74QMMYPDDGQB9FTG",
				"start_at": "2021-08-20T09:00:00Z",
				"end_at": "2021-08-20T17:00:00Z",
				"created_at": "2021-08-17T13:28:57.801578Z",
				"updated_at": "2021-08-17T14:28:57.801578Z"
			}
		}`

		_, _ = fmt.Fprint(w, response)
	})

	ctx := context.Background()
	override, _, err := client.Schedules.UpdateOverride(ctx, scheduleID, overrideID, input)
	if err != nil {
		t.Errorf("Schedules.UpdateOverride returned error: %v", err)
	}

	if override.UserID != newUserID {
		t.Errorf("Override UserID = %s, want %s", override.UserID, newUserID)
	}
}

func TestSchedulesService_DeleteOverride(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	scheduleID := "01G0J1EXE7AXZ2C93K61WBPYEH"
	overrideID := "01G0J1EXE7AXZ2C93K61WBPYEH"

	mux.HandleFunc(fmt.Sprintf("/v2/schedules/%s/overrides/%s", scheduleID, overrideID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	resp, err := client.Schedules.DeleteOverride(ctx, scheduleID, overrideID)
	if err != nil {
		t.Errorf("Schedules.DeleteOverride returned error: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Schedules.DeleteOverride returned status %d, want %d", resp.StatusCode, http.StatusNoContent)
	}
}

func TestSchedulesService_ComplexRotation(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	// Test creating a schedule with complex rotation including working intervals
	input := &CreateScheduleOptions{
		Name:     "24/7 Support",
		Timezone: "UTC",
		Config: &ScheduleConfig{
			Rotations: []RotationConfig{
				{
					ID:   "business-hours",
					Name: "Business Hours",
					Layers: []LayerConfig{
						{
							ID:   "primary",
							Name: "Primary",
							Users: []LayerUser{
								{UserID: "user1"},
								{UserID: "user2"},
							},
						},
						{
							ID:   "secondary",
							Name: "Secondary",
							Users: []LayerUser{
								{UserID: "user3"},
								{UserID: "user4"},
							},
						},
					},
					HandoverStartAt: Timestamp{parseTime("2021-08-16T09:00:00Z")},
					HandoversAt:     Timestamp{parseTime("2021-08-23T09:00:00Z")},
					WorkingInterval: []WorkingIntervalConfig{
						{
							StartTime: "09:00",
							EndTime:   "17:00",
							Weekdays:  []string{"monday", "tuesday", "wednesday", "thursday", "friday"},
						},
					},
				},
				{
					ID:   "after-hours",
					Name: "After Hours",
					Layers: []LayerConfig{
						{
							ID:   "on-call",
							Name: "On-Call",
							Users: []LayerUser{
								{UserID: "user5"},
								{UserID: "user6"},
							},
						},
					},
					HandoverStartAt: Timestamp{parseTime("2021-08-16T17:00:00Z")},
					HandoversAt:     Timestamp{parseTime("2021-08-23T17:00:00Z")},
					WorkingInterval: []WorkingIntervalConfig{
						{
							StartTime: "17:00",
							EndTime:   "09:00",
							Weekdays:  []string{"monday", "tuesday", "wednesday", "thursday", "friday"},
						},
						{
							StartTime: "00:00",
							EndTime:   "23:59",
							Weekdays:  []string{"saturday", "sunday"},
						},
					},
				},
			},
		},
	}

	mux.HandleFunc("/v2/schedules", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		var received CreateScheduleOptions
		err := json.NewDecoder(r.Body).Decode(&received)
		if err != nil {
			t.Errorf("error decoding request body: %v", err)
		}

		// Verify complex structure
		if len(received.Config.Rotations) != 2 {
			t.Errorf("Expected 2 rotations, got %d", len(received.Config.Rotations))
		}

		if len(received.Config.Rotations[0].Layers) != 2 {
			t.Errorf("Expected 2 layers in first rotation, got %d", len(received.Config.Rotations[0].Layers))
		}

		if len(received.Config.Rotations[0].WorkingInterval) != 1 {
			t.Errorf("Expected 1 working interval in first rotation, got %d", len(received.Config.Rotations[0].WorkingInterval))
		}

		response := `{
			"schedule": {
				"id": "01G0J1EXE7AXZ2C93K61WBPYEH",
				"name": "24/7 Support",
				"timezone": "UTC",
				"created_at": "2021-08-17T13:28:57.801578Z",
				"updated_at": "2021-08-17T13:28:57.801578Z"
			}
		}`

		w.WriteHeader(http.StatusCreated)
		_, _ = fmt.Fprint(w, response)
	})

	ctx := context.Background()
	_, _, err := client.Schedules.Create(ctx, input)
	if err != nil {
		t.Errorf("Schedules.Create with complex rotation returned error: %v", err)
	}
}

func TestSchedulesService_ErrorHandling(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/schedules/not-found", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		response := `{
			"type": "validation_error",
			"status": 404,
			"detail": "Schedule not found",
			"errors": [
				{
					"code": "not_found",
					"detail": "No schedule with ID 'not-found'",
					"source": {
						"pointer": "/id"
					}
				}
			]
		}`
		_, _ = fmt.Fprint(w, response)
	})

	ctx := context.Background()
	_, _, err := client.Schedules.Get(ctx, "not-found")
	if err == nil {
		t.Error("Expected error, got nil")
	}

	errResp := &ErrorResponse{}
	ok := errors.As(err, &errResp)
	if !ok {
		t.Errorf("Error type = %T, want *ErrorResponse", err)
	}

	if errResp.Status != 404 {
		t.Errorf("Error status = %d, want %d", errResp.Status, 404)
	}

	if errResp.Detail != "Schedule not found" {
		t.Errorf("Error detail = %s, want %s", errResp.Detail, "Schedule not found")
	}
}
