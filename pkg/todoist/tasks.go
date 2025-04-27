package todoist

import (
	"context"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
)

// GetTasksParams represents the parameters for the todoist_get_tasks tool
type GetTasksParams struct {
	ProjectID string   `json:"projectId,omitempty"`
	SectionID string   `json:"sectionId,omitempty"`
	Label     string   `json:"label,omitempty"`
	Filter    string   `json:"filter,omitempty"`
	Lang      string   `json:"lang,omitempty"`
	IDs       []string `json:"ids,omitempty"`
}

// GetTasksResponse represents the response from the todoist_get_tasks tool
type GetTasksResponse struct {
	Tasks []Task `json:"tasks"`
}

// GetTasks returns the todoist_get_tasks tool
func (tp *ToolProvider) GetTasks() mcp.Tool {
	// Define the input schema for the tool
	inputSchema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"projectId": map[string]interface{}{
				"type":        "string",
				"description": "Project ID",
			},
			"sectionId": map[string]interface{}{
				"type":        "string",
				"description": "Section ID",
			},
			"label": map[string]interface{}{
				"type":        "string",
				"description": "Label name",
			},
			"filter": map[string]interface{}{
				"type":        "string",
				"description": "Filter string",
			},
			"lang": map[string]interface{}{
				"type":        "string",
				"description": "Language code",
			},
			"ids": map[string]interface{}{
				"type":        "array",
				"description": "List of task IDs",
				"items": map[string]interface{}{
					"type": "string",
				},
			},
		},
	}

	// Convert the input schema to JSON
	inputSchemaJSON, err := json.Marshal(inputSchema)
	if err != nil {
		tp.logger.WithError(err).Error("Failed to marshal input schema")
		return mcp.Tool{}
	}

	// Create the tool with read-only annotation
	tool := mcp.NewToolWithRawSchema(
		"todoist_get_tasks",
		"Get a list of tasks.",
		inputSchemaJSON,
	)
	
	// Mark as read-only
	tool.Annotations.ReadOnlyHint = true
	
	return tool
}

// HandleGetTasks handles the todoist_get_tasks tool request
func (tp *ToolProvider) HandleGetTasks(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Parse parameters
	projectID, _ := OptionalParam[string](request, "projectId")
	sectionID, _ := OptionalParam[string](request, "sectionId")
	label, _ := OptionalParam[string](request, "label")
	filter, _ := OptionalParam[string](request, "filter")
	lang, _ := OptionalParam[string](request, "lang")
	
	// Parse IDs array
	ids, _ := OptionalStringArrayParam(request, "ids")

	// Log the request
	tp.logger.WithFields(map[string]interface{}{
		"projectId": projectID,
		"sectionId": sectionID,
		"label":     label,
		"filter":    filter,
		"lang":      lang,
		"ids":       ids,
	}).Info("Getting tasks")

	// Call the Todoist API
	tasks, err := tp.client.GetTasks(projectID, sectionID, label, filter, lang, ids)
	if err != nil {
		tp.logger.WithError(err).Error("Failed to get tasks")
		return mcp.NewToolResultErrorFromErr("Failed to get tasks", err), nil
	}

	// Convert tasks to JSON
	response := GetTasksResponse{
		Tasks: tasks,
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to marshal response", err), nil
	}

	// Return the response
	return mcp.NewToolResultText(string(responseJSON)), nil
}

// OptionalParam is a helper function that can be used to fetch a requested parameter from the request.
// It does the following checks:
// 1. Checks if the parameter is present in the request, if not, it returns its zero-value
// 2. If it is present, it checks if the parameter is of the expected type and returns it
func OptionalParam[T any](r mcp.CallToolRequest, p string) (T, error) {
	var zero T

	// Check if the parameter is present in the request
	if _, ok := r.Params.Arguments[p]; !ok {
		return zero, nil
	}

	// Check if the parameter is of the expected type
	if _, ok := r.Params.Arguments[p].(T); !ok {
		return zero, nil
	}

	return r.Params.Arguments[p].(T), nil
}

// OptionalStringArrayParam is a helper function that can be used to fetch a requested parameter from the request.
// It does the following checks:
// 1. Checks if the parameter is present in the request, if not, it returns its zero-value
// 2. If it is present, iterates the elements and checks each is a string
func OptionalStringArrayParam(r mcp.CallToolRequest, p string) ([]string, error) {
	// Check if the parameter is present in the request
	if _, ok := r.Params.Arguments[p]; !ok {
		return []string{}, nil
	}

	switch v := r.Params.Arguments[p].(type) {
	case nil:
		return []string{}, nil
	case []string:
		return v, nil
	case []interface{}:
		strSlice := make([]string, len(v))
		for i, v := range v {
			s, ok := v.(string)
			if !ok {
				return []string{}, nil
			}
			strSlice[i] = s
		}
		return strSlice, nil
	default:
		return []string{}, nil
	}
}
