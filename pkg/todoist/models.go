package todoist

import "context"

// TodoistClient defines the interface for Todoist API operations
type TodoistClient interface {
	GetTasks(ctx context.Context, projectID, filter string) ([]Task, error)
	GetTask(ctx context.Context, id string) (*Task, error)
	GetProjects(ctx context.Context) ([]Project, error)
	GetProject(ctx context.Context, id string) (*Project, error)
	CreateTask(ctx context.Context, req CreateTaskRequest) (*Task, error)
	UpdateTask(ctx context.Context, id string, req UpdateTaskRequest) (*Task, error)
	CloseTask(ctx context.Context, id string) error
	ReopenTask(ctx context.Context, id string) error
	DeleteTask(ctx context.Context, id string) error
}

// Task represents a Todoist task
type Task struct {
	ID          string `json:"id"`
	Content     string `json:"content"`
	Description string `json:"description"`
	ProjectID   string `json:"project_id"`
	ParentID    string `json:"parent_id,omitempty"`
	Priority    int    `json:"priority"`
	Due         *Due   `json:"due,omitempty"`
	URL         string `json:"url"`
}

// Due represents a due date for a task
type Due struct {
	Date        string `json:"date"`
	IsRecurring bool   `json:"is_recurring"`
	Datetime    string `json:"datetime,omitempty"`
	String      string `json:"string,omitempty"`
	Timezone    string `json:"timezone,omitempty"`
}

// Project represents a Todoist project
type Project struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	CommentCount   int    `json:"comment_count"`
	Order          int    `json:"order"`
	Color          string `json:"color"`
	IsShared       bool   `json:"is_shared"`
	IsFavorite     bool   `json:"is_favorite"`
	IsInboxProject bool   `json:"is_inbox_project"`
	IsTeamInbox    bool   `json:"is_team_inbox"`
	ViewStyle      string `json:"view_style"`
	URL            string `json:"url"`
	ParentID       string `json:"parent_id,omitempty"`
}

// CreateTaskRequest represents the request to create a task
type CreateTaskRequest struct {
	Content     string `json:"content"`
	Description string `json:"description,omitempty"`
	ProjectID   string `json:"project_id,omitempty"`
	ParentID    string `json:"parent_id,omitempty"`
	Order       int    `json:"order,omitempty"`
	Priority    int    `json:"priority,omitempty"`
	DueString   string `json:"due_string,omitempty"`
	DueDate     string `json:"due_date,omitempty"`
	DueDatetime string `json:"due_datetime,omitempty"`
}

// UpdateTaskRequest represents the request to update a task
type UpdateTaskRequest struct {
	Content     string `json:"content,omitempty"`
	Description string `json:"description,omitempty"`
	Priority    int    `json:"priority,omitempty"`
	DueString   string `json:"due_string,omitempty"`
	DueDate     string `json:"due_date,omitempty"`
	DueDatetime string `json:"due_datetime,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}
