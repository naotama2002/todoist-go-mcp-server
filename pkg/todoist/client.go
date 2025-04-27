package todoist

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	// TodoistAPIBaseURL is the base URL for Todoist API
	TodoistAPIBaseURL = "https://api.todoist.com/rest/v2"
	// DefaultTimeout is the default timeout for HTTP requests
	DefaultTimeout = 10 * time.Second
)

// Client represents a Todoist API client
type Client struct {
	httpClient *http.Client
	token      string
	baseURL    string
	logger     *logrus.Logger
}

// NewClient creates a new Todoist API client
func NewClient(token string, logger *logrus.Logger) *Client {
	if token == "" {
		token = os.Getenv("TODOIST_API_TOKEN")
	}

	if logger == nil {
		logger = logrus.New()
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	return &Client{
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		token:   token,
		baseURL: TodoistAPIBaseURL,
		logger:  logger,
	}
}

// Task represents a Todoist task
type Task struct {
	ID          string   `json:"id"`
	Content     string   `json:"content"`
	Description string   `json:"description"`
	ProjectID   string   `json:"project_id"`
	SectionID   string   `json:"section_id,omitempty"`
	ParentID    string   `json:"parent_id,omitempty"`
	Labels      []string `json:"labels"`
	Priority    int      `json:"priority"`
	Due         *Due     `json:"due,omitempty"`
	URL         string   `json:"url"`
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
	ID            string `json:"id"`
	Name          string `json:"name"`
	CommentCount  int    `json:"comment_count"`
	Order         int    `json:"order"`
	Color         string `json:"color"`
	IsShared      bool   `json:"is_shared"`
	IsFavorite    bool   `json:"is_favorite"`
	IsInboxProject bool  `json:"is_inbox_project"`
	IsTeamInbox   bool   `json:"is_team_inbox"`
	ViewStyle     string `json:"view_style"`
	URL           string `json:"url"`
	ParentID      string `json:"parent_id,omitempty"`
}

// GetTasks retrieves all active tasks
func (c *Client) GetTasks(projectID, sectionID, label, filter, lang string, ids []string) ([]Task, error) {
	endpoint := "/tasks"

	req, err := http.NewRequest("GET", c.baseURL+endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add query parameters if provided
	q := req.URL.Query()
	if projectID != "" {
		q.Add("project_id", projectID)
	}
	if sectionID != "" {
		q.Add("section_id", sectionID)
	}
	if label != "" {
		q.Add("label", label)
	}
	if filter != "" {
		q.Add("filter", filter)
	}
	if lang != "" {
		q.Add("lang", lang)
	}
	for _, id := range ids {
		q.Add("ids", id)
	}
	req.URL.RawQuery = q.Encode()

	// Add authorization header
	req.Header.Add("Authorization", "Bearer "+c.token)

	// Execute the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Parse response
	var tasks []Task
	if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return tasks, nil
}

// GetTask retrieves a specific task by ID
func (c *Client) GetTask(id string) (*Task, error) {
	endpoint := fmt.Sprintf("/tasks/%s", id)

	req, err := http.NewRequest("GET", c.baseURL+endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authorization header
	req.Header.Add("Authorization", "Bearer "+c.token)

	// Execute the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Parse response
	var task Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &task, nil
}

// GetProjects retrieves all projects
func (c *Client) GetProjects() ([]Project, error) {
	endpoint := "/projects"

	req, err := http.NewRequest("GET", c.baseURL+endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authorization header
	req.Header.Add("Authorization", "Bearer "+c.token)

	// Execute the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Parse response
	var projects []Project
	if err := json.NewDecoder(resp.Body).Decode(&projects); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return projects, nil
}

// GetProject retrieves a specific project by ID
func (c *Client) GetProject(id string) (*Project, error) {
	endpoint := fmt.Sprintf("/projects/%s", id)

	req, err := http.NewRequest("GET", c.baseURL+endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authorization header
	req.Header.Add("Authorization", "Bearer "+c.token)

	// Execute the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Parse response
	var project Project
	if err := json.NewDecoder(resp.Body).Decode(&project); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &project, nil
}
