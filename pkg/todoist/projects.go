package todoist

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// GetProjectsResponse represents the response from the todoist_get_projects tool
type GetProjectsResponse struct {
	Projects []Project `json:"projects"`
}

// GetProjectParams represents the parameters for the todoist_get_project tool
type GetProjectParams struct {
	ID string `json:"id"`
}

// GetProjectResponse represents the response from the todoist_get_project tool
type GetProjectResponse struct {
	Project Project `json:"project"`
}

// GetProjects returns the todoist_get_projects tool
func (tp *ToolProvider) GetProjects() mcp.Tool {
	// Define the input schema for the tool
	inputSchema := map[string]interface{}{
		"type":       "object",
		"properties": map[string]interface{}{},
	}

	// Convert the input schema to JSON
	inputSchemaJSON, err := json.Marshal(inputSchema)
	if err != nil {
		tp.logger.WithError(err).Error("Failed to marshal input schema")
		return mcp.Tool{}
	}

	return mcp.Tool{
		Name:        "todoist_get_projects",
		Description: "Get a list of projects.",
		InputSchema: json.RawMessage(inputSchemaJSON),
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}
}

// HandleGetProjects handles the todoist_get_projects tool request
func (tp *ToolProvider) HandleGetProjects(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Log the request
	tp.logger.Info("Getting projects")

	// Call the Todoist API
	projects, err := tp.client.GetProjects(ctx)
	if err != nil {
		tp.logger.WithError(err).Error("Failed to get projects")
		return newToolResultError("Failed to get projects", err), nil
	}

	// Convert projects to JSON
	response := GetProjectsResponse{
		Projects: projects,
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		return newToolResultError("Failed to marshal response", err), nil
	}

	// Return the response
	return newToolResultText(string(responseJSON)), nil
}

// GetProject returns the todoist_get_project tool
func (tp *ToolProvider) GetProject() mcp.Tool {
	// Define the input schema for the tool
	inputSchema := map[string]interface{}{
		"type":     "object",
		"required": []string{"id"},
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "string",
				"description": "The unique identifier of the project to retrieve. Specify the numeric Todoist project ID (e.g., '2203306141').",
			},
		},
	}

	// Convert the input schema to JSON
	inputSchemaJSON, err := json.Marshal(inputSchema)
	if err != nil {
		tp.logger.WithError(err).Error("Failed to marshal input schema")
		return mcp.Tool{}
	}

	return mcp.Tool{
		Name:        "todoist_get_project",
		Description: "Get a project by ID.",
		InputSchema: json.RawMessage(inputSchemaJSON),
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}
}

// HandleGetProject handles the todoist_get_project tool request
func (tp *ToolProvider) HandleGetProject(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Parse parameters
	id, err := RequiredParam[string](request, "id")
	if err != nil {
		return newToolResultError(fmt.Sprintf("Invalid parameter: %s", err.Error()), err), nil
	}

	// Log the request
	tp.logger.WithField("id", id).Info("Getting project")

	// Call the Todoist API
	project, err := tp.client.GetProject(ctx, id)
	if err != nil {
		tp.logger.WithError(err).Error("Failed to get project")
		return newToolResultError("Failed to get project", err), nil
	}

	// Convert project to JSON
	response := GetProjectResponse{
		Project: *project,
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		return newToolResultError("Failed to marshal response", err), nil
	}

	// Return the response
	return newToolResultText(string(responseJSON)), nil
}
