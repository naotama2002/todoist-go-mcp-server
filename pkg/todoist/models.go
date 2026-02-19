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

// PaginatedResponse is a generic paginated response from the Todoist API v1
type PaginatedResponse[T any] struct {
	Results    []T     `json:"results"`
	NextCursor *string `json:"next_cursor"`
}

// Task represents a Todoist task (API v1)
type Task struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	Content        string    `json:"content"`
	Description    string    `json:"description"`
	ProjectID      string    `json:"project_id"`
	SectionID      *string   `json:"section_id"`
	ParentID       *string   `json:"parent_id"`
	AddedByUID     *string   `json:"added_by_uid"`
	AssignedByUID  *string   `json:"assigned_by_uid"`
	ResponsibleUID *string   `json:"responsible_uid"`
	Labels         []string  `json:"labels"`
	Deadline       *Deadline `json:"deadline"`
	Duration       *Duration `json:"duration"`
	Checked        bool      `json:"checked"`
	IsDeleted      bool      `json:"is_deleted"`
	AddedAt        *string   `json:"added_at"`
	CompletedAt    *string   `json:"completed_at"`
	UpdatedAt      *string   `json:"updated_at"`
	Due            *Due      `json:"due"`
	Priority       int       `json:"priority"`
	ChildOrder     int       `json:"child_order"`
	NoteCount      int       `json:"note_count"`
	DayOrder       int       `json:"day_order"`
	IsCollapsed    bool      `json:"is_collapsed"`
}

// Due represents a due date for a task
type Due struct {
	Date        string `json:"date"`
	IsRecurring bool   `json:"is_recurring"`
	Datetime    string `json:"datetime,omitempty"`
	String      string `json:"string,omitempty"`
	Timezone    string `json:"timezone,omitempty"`
	Lang        string `json:"lang,omitempty"`
}

// Deadline represents a task deadline
type Deadline struct {
	Date string `json:"date"`
	Lang string `json:"lang"`
}

// Duration represents a task duration
type Duration struct {
	Amount int    `json:"amount"`
	Unit   string `json:"unit"`
}

// Project represents a Todoist project (API v1)
type Project struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	CanAssignTasks bool    `json:"can_assign_tasks"`
	ChildOrder     int     `json:"child_order"`
	Color          string  `json:"color"`
	CreatedAt      *string `json:"created_at"`
	IsArchived     bool    `json:"is_archived"`
	IsDeleted      bool    `json:"is_deleted"`
	IsFavorite     bool    `json:"is_favorite"`
	IsFrozen       bool    `json:"is_frozen"`
	UpdatedAt      *string `json:"updated_at"`
	ViewStyle      string  `json:"view_style"`
	DefaultOrder   int     `json:"default_order"`
	Description    string  `json:"description"`
	IsCollapsed    bool    `json:"is_collapsed"`
	IsShared       bool    `json:"is_shared"`
	ParentID       *string `json:"parent_id"`
	InboxProject   bool    `json:"inbox_project"`
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
