package incidentio

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// WorkflowsService handles communication with the workflows related methods.
type WorkflowsService struct {
	client *Client
}

// WorkflowTrigger identifies what triggers a workflow.
type WorkflowTrigger struct {
	Name  string `json:"name"`
	Label string `json:"label"`
}

// Workflow represents a workflow in Incident.io.
//
// The nested workflow engine structures (condition groups, expressions, steps,
// delay, form fields and annotations) are kept as raw JSON so they can be
// round-tripped without loss.
type Workflow struct {
	ID                        string            `json:"id"`
	Name                      string            `json:"name"`
	Version                   int               `json:"version"`
	State                     string            `json:"state"`
	Trigger                   *WorkflowTrigger  `json:"trigger,omitempty"`
	Folder                    string            `json:"folder,omitempty"`
	Shortform                 string            `json:"shortform,omitempty"`
	OnceFor                   []json.RawMessage `json:"once_for,omitempty"`
	RunsOnIncidents           string            `json:"runs_on_incidents"`
	RunsOnIncidentModes       []string          `json:"runs_on_incident_modes"`
	RunsFrom                  string            `json:"runs_from,omitempty"`
	ContinueOnStepError       bool              `json:"continue_on_step_error"`
	IncludePrivateIncidents   bool              `json:"include_private_incidents"`
	IncludePrivateEscalations bool              `json:"include_private_escalations"`
	PrivateIncidentScope      string            `json:"private_incident_scope,omitempty"`
	AutoRunMode               string            `json:"auto_run_mode,omitempty"`
	OwningTeamIDs             []string          `json:"owning_team_ids,omitempty"`
	ConditionGroups           json.RawMessage   `json:"condition_groups,omitempty"`
	Expressions               json.RawMessage   `json:"expressions,omitempty"`
	Steps                     json.RawMessage   `json:"steps,omitempty"`
	Delay                     json.RawMessage   `json:"delay,omitempty"`
	FormFields                json.RawMessage   `json:"form_fields,omitempty"`
}

// CreateWorkflowOptions represents the options for creating a workflow.
type CreateWorkflowOptions struct {
	Name                      string          `json:"name"`
	Trigger                   string          `json:"trigger"`
	OnceFor                   []string        `json:"once_for"`
	ConditionGroups           json.RawMessage `json:"condition_groups"`
	Steps                     json.RawMessage `json:"steps"`
	Expressions               json.RawMessage `json:"expressions"`
	RunsOnIncidentModes       []string        `json:"runs_on_incident_modes"`
	RunsOnIncidents           string          `json:"runs_on_incidents"`
	ContinueOnStepError       bool            `json:"continue_on_step_error"`
	State                     string          `json:"state,omitempty"`
	Folder                    string          `json:"folder,omitempty"`
	Shortform                 string          `json:"shortform,omitempty"`
	AutoRunMode               string          `json:"auto_run_mode,omitempty"`
	PrivateIncidentScope      string          `json:"private_incident_scope,omitempty"`
	IncludePrivateIncidents   bool            `json:"include_private_incidents,omitempty"`
	IncludePrivateEscalations bool            `json:"include_private_escalations,omitempty"`
	OwningTeamIDs             []string        `json:"owning_team_ids,omitempty"`
	Delay                     json.RawMessage `json:"delay,omitempty"`
	FormFields                json.RawMessage `json:"form_fields,omitempty"`
	Annotations               json.RawMessage `json:"annotations,omitempty"`
}

// UpdateWorkflowOptions represents the options for updating a workflow.
type UpdateWorkflowOptions struct {
	Name                      string          `json:"name"`
	OnceFor                   []string        `json:"once_for"`
	ConditionGroups           json.RawMessage `json:"condition_groups"`
	Steps                     json.RawMessage `json:"steps"`
	Expressions               json.RawMessage `json:"expressions"`
	RunsOnIncidentModes       []string        `json:"runs_on_incident_modes"`
	RunsOnIncidents           string          `json:"runs_on_incidents"`
	ContinueOnStepError       bool            `json:"continue_on_step_error"`
	SkipStepUpgrades          bool            `json:"skip_step_upgrades,omitempty"`
	State                     string          `json:"state,omitempty"`
	Folder                    string          `json:"folder,omitempty"`
	Shortform                 string          `json:"shortform,omitempty"`
	AutoRunMode               string          `json:"auto_run_mode,omitempty"`
	PrivateIncidentScope      string          `json:"private_incident_scope,omitempty"`
	IncludePrivateIncidents   bool            `json:"include_private_incidents,omitempty"`
	IncludePrivateEscalations bool            `json:"include_private_escalations,omitempty"`
	OwningTeamIDs             []string        `json:"owning_team_ids,omitempty"`
	Delay                     json.RawMessage `json:"delay,omitempty"`
	FormFields                json.RawMessage `json:"form_fields,omitempty"`
	Annotations               json.RawMessage `json:"annotations,omitempty"`
}

// List returns a list of workflows.
func (s *WorkflowsService) List(ctx context.Context) ([]*Workflow, *http.Response, error) {
	return getKey[[]*Workflow](ctx, s.client, "GET", "v2/workflows", nil, nil, "workflows")
}

// Get returns a single workflow. When skipStepUpgrades is true, steps are
// returned as stored rather than upgraded to their latest version.
func (s *WorkflowsService) Get(ctx context.Context, id string, skipStepUpgrades bool) (*Workflow, *http.Response, error) {
	q := url.Values{}
	if skipStepUpgrades {
		q.Set("skip_step_upgrades", "true")
	}

	return getKey[*Workflow](ctx, s.client, "GET", fmt.Sprintf("v2/workflows/%s", id), q, nil, "workflow")
}

// Create creates a new workflow.
func (s *WorkflowsService) Create(ctx context.Context, opts *CreateWorkflowOptions) (*Workflow, *http.Response, error) {
	return getKey[*Workflow](ctx, s.client, "POST", "v2/workflows", nil, opts, "workflow")
}

// Update updates a workflow.
func (s *WorkflowsService) Update(ctx context.Context, id string, opts *UpdateWorkflowOptions) (*Workflow, *http.Response, error) {
	return getKey[*Workflow](ctx, s.client, "PUT", fmt.Sprintf("v2/workflows/%s", id), nil, opts, "workflow")
}

// Delete deletes a workflow.
func (s *WorkflowsService) Delete(ctx context.Context, id string) (*http.Response, error) {
	return s.client.send(ctx, "DELETE", fmt.Sprintf("v2/workflows/%s", id), nil, nil, nil)
}
