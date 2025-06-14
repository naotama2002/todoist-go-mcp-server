package todoist

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

// GetTasksParams represents the parameters for the todoist_get_tasks tool
type GetTasksParams struct {
	ProjectID string `json:"projectId,omitempty"`
	Filter    string `json:"filter,omitempty"`
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
	Content     string `json:"content"`
	Description string `json:"description,omitempty"`
	ProjectID   string `json:"projectId,omitempty"`
	ParentID    string `json:"parentId,omitempty"`
	Order       int    `json:"order,omitempty"`
	Priority    int    `json:"priority,omitempty"`
	DueString   string `json:"dueString,omitempty"`
	DueDate     string `json:"dueDate,omitempty"`
	DueDatetime string `json:"dueDatetime,omitempty"`
}

// CreateTaskResponse represents the response from the todoist_create_task tool
type CreateTaskResponse struct {
	Task Task `json:"task"`
}

// UpdateTaskParams represents the parameters for the todoist_update_task tool
type UpdateTaskParams struct {
	ID          string `json:"id"`
	Content     string `json:"content,omitempty"`
	Description string `json:"description,omitempty"`
	Priority    int    `json:"priority,omitempty"`
	DueString   string `json:"dueString,omitempty"`
	DueDate     string `json:"dueDate,omitempty"`
	DueDatetime string `json:"dueDatetime,omitempty"`
}

// UpdateTaskResponse represents the response from the todoist_update_task tool
type UpdateTaskResponse struct {
	Task Task `json:"task"`
}

// CloseTaskParams represents the parameters for the todoist_close_task tool
type CloseTaskParams struct {
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
				"description": "Filter tasks by project ID. Retrieves only tasks belonging to the specified project.",
			},
			"filter": map[string]interface{}{
				"type":        "string",
				"description": "Todoist filter query using the Todoist filter syntax. Examples: 'today', 'tomorrow', 'next week', 'overdue', 'priority 1', 'search: meeting', 'date: 2023-12-31', 'no date'. For comprehensive filter rules and examples, use the todoist_get_task_filter_rules tool to get detailed information about available filter syntax.",
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
	tool.Annotations.ReadOnlyHint = mcp.ToBoolPtr(true)

	return tool
}

// HandleGetTasks handles the todoist_get_tasks tool request
func (tp *ToolProvider) HandleGetTasks(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Parse parameters
	projectID, _ := OptionalParam[string](request, "projectId")
	filter, _ := OptionalParam[string](request, "filter")

	// Log the request
	tp.logger.WithFields(map[string]interface{}{
		"projectId": projectID,
		"filter":    filter,
	}).Info("Getting tasks")

	// Call the Todoist API
	tasks, err := tp.client.GetTasks(ctx, projectID, filter)
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
		"type":     "object",
		"required": []string{"id"},
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "string",
				"description": "The unique identifier of the task to retrieve. Specify the numeric Todoist task ID (e.g., '2995104339').",
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
	tool.Annotations.ReadOnlyHint = mcp.ToBoolPtr(true)

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
	task, err := tp.client.GetTask(ctx, id)
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
		"type":     "object",
		"required": []string{"content"},
		"properties": map[string]interface{}{
			"content": map[string]interface{}{
				"type":        "string",
				"description": "The content of the task (required). Supports text formatting using Markdown syntax. See https://todoist.com/help/articles/format-text-in-a-todoist-task for formatting options.",
			},
			"description": map[string]interface{}{
				"type":        "string",
				"description": "Detailed description or notes for the task. Supports Markdown formatting for rich text.",
			},
			"projectId": map[string]interface{}{
				"type":        "string",
				"description": "Project ID to assign the task to. If not specified, the task will be added to the Inbox project.",
			},
			"parentId": map[string]interface{}{
				"type":        "string",
				"description": "Parent task ID for creating subtasks. The task will be created as a child of this task.",
			},
			"order": map[string]interface{}{
				"type":        "integer",
				"description": "Order value for positioning the task within its parent or project. Tasks are sorted by this value in ascending order.",
			},
			"priority": map[string]interface{}{
				"type":        "integer",
				"description": "Task priority: 4 (normal, default), 3 (medium), 2 (high), 1 (urgent). Note that 1 is the highest priority, 4 is the lowest.",
				"minimum":     1,
				"maximum":     4,
			},
			"dueString": map[string]interface{}{
				"type":        "string",
				"description": "Due date in natural language, e.g., 'today', 'tomorrow', 'next Monday', 'Jan 15'. Only one of dueString, dueDate, or dueDatetime should be used.",
			},
			"dueDate": map[string]interface{}{
				"type":        "string",
				"description": "Due date in YYYY-MM-DD format, e.g., '2023-12-31'. Only one of dueString, dueDate, or dueDatetime should be used.",
			},
			"dueDatetime": map[string]interface{}{
				"type":        "string",
				"description": "Due date and time in RFC3339 format, e.g., '2023-12-31T10:00:00Z'. Only one of dueString, dueDate, or dueDatetime should be used.",
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
	parentID, _ := OptionalParam[string](request, "parentId")
	order, _ := OptionalParam[int](request, "order")
	priority, _ := OptionalParam[int](request, "priority")
	dueString, _ := OptionalParam[string](request, "dueString")
	dueDate, _ := OptionalParam[string](request, "dueDate")
	dueDatetime, _ := OptionalParam[string](request, "dueDatetime")

	// Log the request
	tp.logger.WithFields(map[string]interface{}{
		"content":     content,
		"description": description,
		"projectId":   projectID,
		"parentId":    parentID,
		"priority":    priority,
	}).Info("Creating task")

	// Create request
	createReq := CreateTaskRequest{
		Content:     content,
		Description: description,
		ProjectID:   projectID,
		ParentID:    parentID,
		Order:       order,
		Priority:    priority,
		DueString:   dueString,
		DueDate:     dueDate,
		DueDatetime: dueDatetime,
	}

	// Call the Todoist API
	task, err := tp.client.CreateTask(ctx, createReq)
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
		"type":     "object",
		"required": []string{"id"},
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "string",
				"description": "The unique identifier of the task to update (required). Specify the numeric Todoist task ID (e.g., '2995104339').",
			},
			"content": map[string]interface{}{
				"type":        "string",
				"description": "The new content of the task. Supports text formatting using Markdown syntax. See https://todoist.com/help/articles/format-text-in-a-todoist-task for formatting options.",
			},
			"description": map[string]interface{}{
				"type":        "string",
				"description": "Detailed description or notes for the task. Supports Markdown formatting for rich text.",
			},
			"priority": map[string]interface{}{
				"type":        "integer",
				"description": "Task priority: 4 (normal, default), 3 (medium), 2 (high), 1 (urgent). Note that 1 is the highest priority, 4 is the lowest.",
				"minimum":     1,
				"maximum":     4,
			},
			"dueString": map[string]interface{}{
				"type":        "string",
				"description": "Due date in natural language, e.g., 'today', 'tomorrow', 'next Monday', 'Jan 15'. Only one of dueString, dueDate, or dueDatetime should be used.",
			},
			"dueDate": map[string]interface{}{
				"type":        "string",
				"description": "Due date in YYYY-MM-DD format, e.g., '2023-12-31'. Only one of dueString, dueDate, or dueDatetime should be used.",
			},
			"dueDatetime": map[string]interface{}{
				"type":        "string",
				"description": "Due date and time in RFC3339 format, e.g., '2023-12-31T10:00:00Z'. Only one of dueString, dueDate, or dueDatetime should be used.",
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
	priority, _ := OptionalParam[int](request, "priority")
	dueString, _ := OptionalParam[string](request, "dueString")
	dueDate, _ := OptionalParam[string](request, "dueDate")
	dueDatetime, _ := OptionalParam[string](request, "dueDatetime")

	// Log the request
	tp.logger.WithFields(map[string]interface{}{
		"id":          id,
		"content":     content,
		"description": description,
		"priority":    priority,
	}).Info("Updating task")

	// Create request
	updateReq := UpdateTaskRequest{
		Content:     content,
		Description: description,
		Priority:    priority,
		DueString:   dueString,
		DueDate:     dueDate,
		DueDatetime: dueDatetime,
	}

	// Call the Todoist API
	task, err := tp.client.UpdateTask(ctx, id, updateReq)
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
		"type":     "object",
		"required": []string{"id"},
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "string",
				"description": "The unique identifier of the task to mark as completed (required). Specify the numeric Todoist task ID (e.g., '2995104339').",
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
	err = tp.client.CloseTask(ctx, id)
	if err != nil {
		tp.logger.WithError(err).Error("Failed to close task")
		return mcp.NewToolResultErrorFromErr("Failed to close task", err), nil
	}

	// Return success response
	return mcp.NewToolResultText(`{"success": true}`), nil
}

// DeleteTask returns the todoist_delete_task tool
func (tp *ToolProvider) DeleteTask() mcp.Tool {
	// Define the input schema for the tool
	inputSchema := map[string]interface{}{
		"type":     "object",
		"required": []string{"id"},
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "string",
				"description": "The unique identifier of the task to delete (required). Specify the numeric Todoist task ID (e.g., '2995104339'). Warning: This action is permanent and cannot be undone.",
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
	err = tp.client.DeleteTask(ctx, id)
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
	args := r.GetArguments()

	// Check if the parameter is present in the request
	if _, ok := args[p]; !ok {
		return zero, nil
	}

	// Check if the parameter is of the expected type
	if _, ok := args[p].(T); !ok {
		return zero, fmt.Errorf("parameter %s is not of the expected type", p)
	}

	return args[p].(T), nil
}

// OptionalStringArrayParam is a helper function that can be used to fetch a requested parameter from the request.
// It does the following checks:
// 1. Checks if the parameter is present in the request, if not, it returns nil
// 2. If it is present, it checks if the parameter is an array
// 3. If it is an array, it checks each element is a string
func OptionalStringArrayParam(r mcp.CallToolRequest, p string) ([]string, error) {
	args := r.GetArguments()
	// Check if the parameter is present in the request
	if _, ok := args[p]; !ok {
		return []string{}, nil
	}

	switch v := args[p].(type) {
	case nil:
		return []string{}, nil
	case []string:
		return v, nil
	case []interface{}:
		strSlice := make([]string, len(v))
		for i, val := range v {
			s, ok := val.(string)
			if !ok {
				return nil, fmt.Errorf("parameter %s contains non-string elements", p)
			}
			strSlice[i] = s
		}
		return strSlice, nil
	default:
		return nil, fmt.Errorf("parameter %s is not an array", p)
	}
}

// RequiredParam is a helper function that can be used to fetch a required parameter from the request.
// It returns an error if the parameter is not present or not of the expected type.
func RequiredParam[T any](r mcp.CallToolRequest, p string) (T, error) {
	var zero T
	args := r.GetArguments()

	// Check if the parameter is present in the request
	if _, ok := args[p]; !ok {
		return zero, fmt.Errorf("parameter %s is required", p)
	}

	// Check if the parameter is of the expected type
	if _, ok := args[p].(T); !ok {
		return zero, fmt.Errorf("parameter %s is not of the expected type", p)
	}

	return args[p].(T), nil
}
