package todoist

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/mark3labs/mcp-go/mcp"
)

// MockHTTPClient is a mock HTTP client for testing
type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

// RoundTrip implements the http.RoundTripper interface
func (m *MockHTTPClient) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

// Do implements the http.Client interface
func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.RoundTrip(req)
}

// MockResponse creates a mock HTTP response with the given status code and body
func MockResponse(statusCode int, body interface{}) *http.Response {
	var bodyBytes []byte
	var err error

	if body != nil {
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			panic(err)
		}
	}

	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(bytes.NewReader(bodyBytes)),
		Header:     make(http.Header),
	}
}

// MockServer creates a test server that returns the given response for all requests
func MockServer(handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(handler)
}

// NewMockClient creates a new Client with a mock HTTP client
func NewMockClient(doFunc func(req *http.Request) (*http.Response, error)) *Client {
	mockClient := &MockHTTPClient{
		DoFunc: doFunc,
	}

	client := NewClient("test-token", nil)
	client.httpClient = &http.Client{
		Transport: mockClient,
	}
	return client
}

// MockTask returns a mock Task for testing
func MockTask() *Task {
	return &Task{
		ID:          "123456789",
		Content:     "Test Task",
		Description: "This is a test task",
		ProjectID:   "987654321",
		SectionID:   "123123123",
		ParentID:    "456456456",
		Labels:      []string{"test", "mock"},
		Priority:    4,
		Due: &Due{
			Date:        "2023-12-31",
			IsRecurring: false,
			Datetime:    "2023-12-31T23:59:59Z",
			String:      "Dec 31",
			Timezone:    "UTC",
		},
		URL: "https://todoist.com/showTask?id=123456789",
	}
}

// MockProject returns a mock Project for testing
func MockProject() *Project {
	return &Project{
		ID:             "987654321",
		Name:           "Test Project",
		CommentCount:   0,
		Order:          1,
		Color:          "red",
		IsShared:       false,
		IsFavorite:     true,
		IsInboxProject: false,
		IsTeamInbox:    false,
		ViewStyle:      "list",
		URL:            "https://todoist.com/showProject?id=987654321",
		ParentID:       "",
	}
}

// MockToolProviderWithHandlers is a mock ToolProvider with custom handlers for testing
type MockToolProviderWithHandlers struct {
	*ToolProvider
}

// NewMockToolProviderWithHandlers creates a new MockToolProviderWithHandlers
func NewMockToolProviderWithHandlers() *MockToolProviderWithHandlers {
	// Create a mock client that returns successful responses for all requests
	mockClient := NewMockClient(func(req *http.Request) (*http.Response, error) {
		var responseData interface{}
		
		// Determine the response data based on the request path
		switch req.URL.Path {
		case "/tasks":
			responseData = map[string]interface{}{
				"tasks": []interface{}{MockTask()},
			}
		case "/tasks/123456789":
			responseData = MockTask()
		case "/projects":
			responseData = map[string]interface{}{
				"projects": []interface{}{MockProject()},
			}
		case "/projects/987654321":
			responseData = MockProject()
		default:
			// Default response for other requests
			responseData = map[string]interface{}{
				"success": true,
			}
		}
		
		return MockResponse(200, responseData), nil
	})
	
	// Create a tool provider with the mock client
	provider := &ToolProvider{
		client: mockClient,
		logger: nil,
	}
	
	return &MockToolProviderWithHandlers{
		ToolProvider: provider,
	}
}

// HandleToolCall handles a tool call for testing
func (m *MockToolProviderWithHandlers) HandleToolCall(ctx context.Context, toolName string, params map[string]interface{}) (*mcp.CallToolResult, error) {
	// Create a mock result based on the tool name
	var content string
	
	switch toolName {
	case "todoist_get_tasks":
		content = `{"tasks":[{"id":"123456789","content":"Test Task"}]}`
	case "todoist_get_task":
		content = `{"id":"123456789","content":"Test Task"}`
	case "todoist_create_task":
		content = `{"id":"123456789","content":"Test Task"}`
	case "todoist_update_task":
		content = `{"id":"123456789","content":"Updated Task"}`
	case "todoist_close_task":
		content = `{"success":true}`
	case "todoist_delete_task":
		content = `{"success":true}`
	case "todoist_get_projects":
		content = `{"projects":[{"id":"987654321","name":"Test Project"}]}`
	case "todoist_get_project":
		content = `{"id":"987654321","name":"Test Project"}`
	case "mock_tool":
		// For the MockCallToolRequest test
		contentBytes, _ := json.Marshal(params)
		content = string(contentBytes)
	default:
		return nil, errors.New("unknown tool: " + toolName)
	}
	
	// Create the result
	result := &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(content),
		},
		IsError: false,
	}
	
	return result, nil
}
