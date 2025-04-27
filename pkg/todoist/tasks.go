package todoist

import (
	"context"
	"encoding/json"
	"fmt"

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

// GetTaskParams represents the parameters for the todoist_get_task tool
type GetTaskParams struct {
	ID string `json:"id"`
}

// GetTaskResponse represents the response from the todoist_get_task tool
type GetTaskResponse struct {
	Task Task `json:"task"`
}

// CreateTaskParams represents the parameters for the todoist_create_task tool
type CreateTaskParams struct {
	Content     string   `json:"content"`
	Description string   `json:"description,omitempty"`
	ProjectID   string   `json:"projectId,omitempty"`
	SectionID   string   `json:"sectionId,omitempty"`
	ParentID    string   `json:"parentId,omitempty"`
	Order       int      `json:"order,omitempty"`
	Labels      []string `json:"labels,omitempty"`
	Priority    int      `json:"priority,omitempty"`
	DueString   string   `json:"dueString,omitempty"`
	DueDate     string   `json:"dueDate,omitempty"`
	DueDatetime string   `json:"dueDatetime,omitempty"`
	DueLang     string   `json:"dueLang,omitempty"`
}

// CreateTaskResponse represents the response from the todoist_create_task tool
type CreateTaskResponse struct {
	Task Task `json:"task"`
}

// UpdateTaskParams represents the parameters for the todoist_update_task tool
type UpdateTaskParams struct {
	ID          string   `json:"id"`
	Content     string   `json:"content,omitempty"`
	Description string   `json:"description,omitempty"`
	Labels      []string `json:"labels,omitempty"`
	Priority    int      `json:"priority,omitempty"`
	DueString   string   `json:"dueString,omitempty"`
	DueDate     string   `json:"dueDate,omitempty"`
	DueDatetime string   `json:"dueDatetime,omitempty"`
	DueLang     string   `json:"dueLang,omitempty"`
}

// UpdateTaskResponse represents the response from the todoist_update_task tool
type UpdateTaskResponse struct {
	Task Task `json:"task"`
}

// CloseTaskParams represents the parameters for the todoist_close_task tool
type CloseTaskParams struct {
	ID string `json:"id"`
}

// ReopenTaskParams represents the parameters for the todoist_reopen_task tool
type ReopenTaskParams struct {
	ID string `json:"id"`
}

// DeleteTaskParams represents the parameters for the todoist_delete_task tool
type DeleteTaskParams struct {
	ID string `json:"id"`
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

// GetTask returns the todoist_get_task tool
func (tp *ToolProvider) GetTask() mcp.Tool {
	// Define the input schema for the tool
	inputSchema := map[string]interface{}{
		"type": "object",
		"required": []string{"id"},
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "string",
				"description": "Task ID",
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
		"todoist_get_task",
		"Get a specific task by ID.",
		inputSchemaJSON,
	)
	
	// Mark as read-only
	tool.Annotations.ReadOnlyHint = true
	
	return tool
}

// HandleGetTask handles the todoist_get_task tool request
func (tp *ToolProvider) HandleGetTask(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Parse parameters
	id, err := RequiredParam[string](request, "id")
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Missing required parameter: id", err), nil
	}

	// Log the request
	tp.logger.WithFields(map[string]interface{}{
		"id": id,
	}).Info("Getting task")

	// Call the Todoist API
	task, err := tp.client.GetTask(id)
	if err != nil {
		tp.logger.WithError(err).Error("Failed to get task")
		return mcp.NewToolResultErrorFromErr("Failed to get task", err), nil
	}

	// Convert task to JSON
	response := GetTaskResponse{
		Task: *task,
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to marshal response", err), nil
	}

	// Return the response
	return mcp.NewToolResultText(string(responseJSON)), nil
}

// CreateTask returns the todoist_create_task tool
func (tp *ToolProvider) CreateTask() mcp.Tool {
	// Define the input schema for the tool
	inputSchema := map[string]interface{}{
		"type": "object",
		"required": []string{"content"},
		"properties": map[string]interface{}{
			"content": map[string]interface{}{
				"type":        "string",
				"description": "Task content",
			},
			"description": map[string]interface{}{
				"type":        "string",
				"description": "Task description",
			},
			"projectId": map[string]interface{}{
				"type":        "string",
				"description": "Project ID",
			},
			"sectionId": map[string]interface{}{
				"type":        "string",
				"description": "Section ID",
			},
			"parentId": map[string]interface{}{
				"type":        "string",
				"description": "Parent task ID",
			},
			"order": map[string]interface{}{
				"type":        "integer",
				"description": "Task order",
			},
			"labels": map[string]interface{}{
				"type":        "array",
				"description": "Task labels",
				"items": map[string]interface{}{
					"type": "string",
				},
			},
			"priority": map[string]interface{}{
				"type":        "integer",
				"description": "Task priority (1-4)",
				"minimum":     1,
				"maximum":     4,
			},
			"dueString": map[string]interface{}{
				"type":        "string",
				"description": "Due date in natural language",
			},
			"dueDate": map[string]interface{}{
				"type":        "string",
				"description": "Due date in YYYY-MM-DD format",
			},
			"dueDatetime": map[string]interface{}{
				"type":        "string",
				"description": "Due date and time in RFC3339 format",
			},
			"dueLang": map[string]interface{}{
				"type":        "string",
				"description": "Language for parsing due_string",
			},
		},
	}

	// Convert the input schema to JSON
	inputSchemaJSON, err := json.Marshal(inputSchema)
	if err != nil {
		tp.logger.WithError(err).Error("Failed to marshal input schema")
		return mcp.Tool{}
	}

	// Create the tool
	tool := mcp.NewToolWithRawSchema(
		"todoist_create_task",
		"Create a new task.",
		inputSchemaJSON,
	)
	
	return tool
}

// HandleCreateTask handles the todoist_create_task tool request
func (tp *ToolProvider) HandleCreateTask(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Parse parameters
	content, err := RequiredParam[string](request, "content")
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Missing required parameter: content", err), nil
	}
	
	description, _ := OptionalParam[string](request, "description")
	projectID, _ := OptionalParam[string](request, "projectId")
	sectionID, _ := OptionalParam[string](request, "sectionId")
	parentID, _ := OptionalParam[string](request, "parentId")
	order, _ := OptionalParam[int](request, "order")
	labels, _ := OptionalStringArrayParam(request, "labels")
	priority, _ := OptionalParam[int](request, "priority")
	dueString, _ := OptionalParam[string](request, "dueString")
	dueDate, _ := OptionalParam[string](request, "dueDate")
	dueDatetime, _ := OptionalParam[string](request, "dueDatetime")
	dueLang, _ := OptionalParam[string](request, "dueLang")

	// Log the request
	tp.logger.WithFields(map[string]interface{}{
		"content":     content,
		"description": description,
		"projectId":   projectID,
		"sectionId":   sectionID,
		"parentId":    parentID,
		"labels":      labels,
		"priority":    priority,
	}).Info("Creating task")

	// Create request
	createReq := CreateTaskRequest{
		Content:     content,
		Description: description,
		ProjectID:   projectID,
		SectionID:   sectionID,
		ParentID:    parentID,
		Order:       order,
		Labels:      labels,
		Priority:    priority,
		DueString:   dueString,
		DueDate:     dueDate,
		DueDatetime: dueDatetime,
		DueLang:     dueLang,
	}

	// Call the Todoist API
	task, err := tp.client.CreateTask(createReq)
	if err != nil {
		tp.logger.WithError(err).Error("Failed to create task")
		return mcp.NewToolResultErrorFromErr("Failed to create task", err), nil
	}

	// Convert task to JSON
	response := CreateTaskResponse{
		Task: *task,
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to marshal response", err), nil
	}

	// Return the response
	return mcp.NewToolResultText(string(responseJSON)), nil
}

// UpdateTask returns the todoist_update_task tool
func (tp *ToolProvider) UpdateTask() mcp.Tool {
	// Define the input schema for the tool
	inputSchema := map[string]interface{}{
		"type": "object",
		"required": []string{"id"},
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "string",
				"description": "Task ID",
			},
			"content": map[string]interface{}{
				"type":        "string",
				"description": "Task content",
			},
			"description": map[string]interface{}{
				"type":        "string",
				"description": "Task description",
			},
			"labels": map[string]interface{}{
				"type":        "array",
				"description": "Task labels",
				"items": map[string]interface{}{
					"type": "string",
				},
			},
			"priority": map[string]interface{}{
				"type":        "integer",
				"description": "Task priority (1-4)",
				"minimum":     1,
				"maximum":     4,
			},
			"dueString": map[string]interface{}{
				"type":        "string",
				"description": "Due date in natural language",
			},
			"dueDate": map[string]interface{}{
				"type":        "string",
				"description": "Due date in YYYY-MM-DD format",
			},
			"dueDatetime": map[string]interface{}{
				"type":        "string",
				"description": "Due date and time in RFC3339 format",
			},
			"dueLang": map[string]interface{}{
				"type":        "string",
				"description": "Language for parsing due_string",
			},
		},
	}

	// Convert the input schema to JSON
	inputSchemaJSON, err := json.Marshal(inputSchema)
	if err != nil {
		tp.logger.WithError(err).Error("Failed to marshal input schema")
		return mcp.Tool{}
	}

	// Create the tool
	tool := mcp.NewToolWithRawSchema(
		"todoist_update_task",
		"Update an existing task.",
		inputSchemaJSON,
	)
	
	return tool
}

// HandleUpdateTask handles the todoist_update_task tool request
func (tp *ToolProvider) HandleUpdateTask(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Parse parameters
	id, err := RequiredParam[string](request, "id")
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Missing required parameter: id", err), nil
	}
	
	content, _ := OptionalParam[string](request, "content")
	description, _ := OptionalParam[string](request, "description")
	labels, _ := OptionalStringArrayParam(request, "labels")
	priority, _ := OptionalParam[int](request, "priority")
	dueString, _ := OptionalParam[string](request, "dueString")
	dueDate, _ := OptionalParam[string](request, "dueDate")
	dueDatetime, _ := OptionalParam[string](request, "dueDatetime")
	dueLang, _ := OptionalParam[string](request, "dueLang")

	// Log the request
	tp.logger.WithFields(map[string]interface{}{
		"id":          id,
		"content":     content,
		"description": description,
		"labels":      labels,
		"priority":    priority,
	}).Info("Updating task")

	// Create request
	updateReq := UpdateTaskRequest{
		Content:     content,
		Description: description,
		Labels:      labels,
		Priority:    priority,
		DueString:   dueString,
		DueDate:     dueDate,
		DueDatetime: dueDatetime,
		DueLang:     dueLang,
	}

	// Call the Todoist API
	task, err := tp.client.UpdateTask(id, updateReq)
	if err != nil {
		tp.logger.WithError(err).Error("Failed to update task")
		return mcp.NewToolResultErrorFromErr("Failed to update task", err), nil
	}

	// Convert task to JSON
	response := UpdateTaskResponse{
		Task: *task,
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to marshal response", err), nil
	}

	// Return the response
	return mcp.NewToolResultText(string(responseJSON)), nil
}

// CloseTask returns the todoist_close_task tool
func (tp *ToolProvider) CloseTask() mcp.Tool {
	// Define the input schema for the tool
	inputSchema := map[string]interface{}{
		"type": "object",
		"required": []string{"id"},
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "string",
				"description": "Task ID",
			},
		},
	}

	// Convert the input schema to JSON
	inputSchemaJSON, err := json.Marshal(inputSchema)
	if err != nil {
		tp.logger.WithError(err).Error("Failed to marshal input schema")
		return mcp.Tool{}
	}

	// Create the tool
	tool := mcp.NewToolWithRawSchema(
		"todoist_close_task",
		"Mark a task as completed.",
		inputSchemaJSON,
	)
	
	return tool
}

// HandleCloseTask handles the todoist_close_task tool request
func (tp *ToolProvider) HandleCloseTask(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Parse parameters
	id, err := RequiredParam[string](request, "id")
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Missing required parameter: id", err), nil
	}

	// Log the request
	tp.logger.WithFields(map[string]interface{}{
		"id": id,
	}).Info("Closing task")

	// Call the Todoist API
	err = tp.client.CloseTask(id)
	if err != nil {
		tp.logger.WithError(err).Error("Failed to close task")
		return mcp.NewToolResultErrorFromErr("Failed to close task", err), nil
	}

	// Return success response
	return mcp.NewToolResultText(`{"success": true}`), nil
}

// ReopenTask returns the todoist_reopen_task tool
func (tp *ToolProvider) ReopenTask() mcp.Tool {
	// Define the input schema for the tool
	inputSchema := map[string]interface{}{
		"type": "object",
		"required": []string{"id"},
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "string",
				"description": "Task ID",
			},
		},
	}

	// Convert the input schema to JSON
	inputSchemaJSON, err := json.Marshal(inputSchema)
	if err != nil {
		tp.logger.WithError(err).Error("Failed to marshal input schema")
		return mcp.Tool{}
	}

	// Create the tool
	tool := mcp.NewToolWithRawSchema(
		"todoist_reopen_task",
		"Mark a task as not completed.",
		inputSchemaJSON,
	)
	
	return tool
}

// HandleReopenTask handles the todoist_reopen_task tool request
func (tp *ToolProvider) HandleReopenTask(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Parse parameters
	id, err := RequiredParam[string](request, "id")
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Missing required parameter: id", err), nil
	}

	// Log the request
	tp.logger.WithFields(map[string]interface{}{
		"id": id,
	}).Info("Reopening task")

	// Call the Todoist API
	err = tp.client.ReopenTask(id)
	if err != nil {
		tp.logger.WithError(err).Error("Failed to reopen task")
		return mcp.NewToolResultErrorFromErr("Failed to reopen task", err), nil
	}

	// Return success response
	return mcp.NewToolResultText(`{"success": true}`), nil
}

// DeleteTask returns the todoist_delete_task tool
func (tp *ToolProvider) DeleteTask() mcp.Tool {
	// Define the input schema for the tool
	inputSchema := map[string]interface{}{
		"type": "object",
		"required": []string{"id"},
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "string",
				"description": "Task ID",
			},
		},
	}

	// Convert the input schema to JSON
	inputSchemaJSON, err := json.Marshal(inputSchema)
	if err != nil {
		tp.logger.WithError(err).Error("Failed to marshal input schema")
		return mcp.Tool{}
	}

	// Create the tool
	tool := mcp.NewToolWithRawSchema(
		"todoist_delete_task",
		"Delete a task.",
		inputSchemaJSON,
	)
	
	return tool
}

// HandleDeleteTask handles the todoist_delete_task tool request
func (tp *ToolProvider) HandleDeleteTask(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Parse parameters
	id, err := RequiredParam[string](request, "id")
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Missing required parameter: id", err), nil
	}

	// Log the request
	tp.logger.WithFields(map[string]interface{}{
		"id": id,
	}).Info("Deleting task")

	// Call the Todoist API
	err = tp.client.DeleteTask(id)
	if err != nil {
		tp.logger.WithError(err).Error("Failed to delete task")
		return mcp.NewToolResultErrorFromErr("Failed to delete task", err), nil
	}

	// Return success response
	return mcp.NewToolResultText(`{"success": true}`), nil
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

// RequiredParam is a helper function that can be used to fetch a required parameter from the request.
// It returns an error if the parameter is not present or not of the expected type.
func RequiredParam[T any](r mcp.CallToolRequest, p string) (T, error) {
	var zero T

	// Check if the parameter is present in the request
	if _, ok := r.Params.Arguments[p]; !ok {
		return zero, fmt.Errorf("parameter %s is required", p)
	}

	// Check if the parameter is of the expected type
	if _, ok := r.Params.Arguments[p].(T); !ok {
		return zero, fmt.Errorf("parameter %s is not of the expected type", p)
	}

	return r.Params.Arguments[p].(T), nil
}
