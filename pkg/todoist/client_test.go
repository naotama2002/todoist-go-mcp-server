package todoist

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	// Create a client with a valid token
	logger := logrus.New()
	client := NewClient("valid_token", WithLogger(logger))

	// Check that the client was created correctly
	assert.NotNil(t, client)
	assert.Equal(t, "valid_token", client.token)
	assert.Equal(t, "https://api.todoist.com/api/v1", client.baseURL)
	assert.NotNil(t, client.httpClient)
	assert.Equal(t, logger, client.logger)
}

func TestGetTasks(t *testing.T) {
	// モックタスクのデリファレンス
	mockTask := *MockTask()

	tests := []struct {
		name      string
		projectID string
		filter    string
		mockTasks []Task
		mockErr   error
		wantErr   bool
	}{
		{
			name:      "success",
			projectID: "123456789",
			filter:    "",
			mockTasks: []Task{mockTask, mockTask},
			mockErr:   nil,
			wantErr:   false,
		},
		{
			name:      "empty response",
			projectID: "",
			filter:    "",
			mockTasks: []Task{},
			mockErr:   nil,
			wantErr:   false,
		},
		{
			name:      "api error",
			projectID: "",
			filter:    "",
			mockTasks: nil,
			mockErr:   errors.New("api error"),
			wantErr:   true,
		},
		{
			name:      "filter uses filter endpoint",
			projectID: "",
			filter:    "today",
			mockTasks: []Task{mockTask},
			mockErr:   nil,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock client
			client := NewMockClient(func(req *http.Request) (*http.Response, error) {
				if tt.mockErr != nil {
					return nil, tt.mockErr
				}
				// Return paginated response
				paginatedResp := MockPaginatedTasks(tt.mockTasks)
				return MockResponse(200, paginatedResp), nil
			})

			// Call the method
			tasks, err := client.GetTasks(context.Background(), tt.projectID, tt.filter)

			// Check error
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.mockTasks), len(tasks))
				if len(tt.mockTasks) > 0 {
					assert.Equal(t, tt.mockTasks[0].ID, tasks[0].ID)
					assert.Equal(t, tt.mockTasks[0].Content, tasks[0].Content)
				}
			}
		})
	}
}

func TestGetTasksFilterEndpoint(t *testing.T) {
	// Verify that filter requests go to /tasks/filter with query param
	client := NewMockClient(func(req *http.Request) (*http.Response, error) {
		assert.Contains(t, req.URL.Path, "/tasks/filter")
		assert.Equal(t, "today", req.URL.Query().Get("query"))
		paginatedResp := MockPaginatedTasks([]Task{*MockTask()})
		return MockResponse(200, paginatedResp), nil
	})

	tasks, err := client.GetTasks(context.Background(), "", "today")
	assert.NoError(t, err)
	assert.Len(t, tasks, 1)
}

func TestGetTask(t *testing.T) {
	// モックタスクを取得
	mockTask := MockTask()

	tests := []struct {
		name     string
		id       string
		mockTask *Task
		mockErr  error
		wantErr  bool
	}{
		{
			name:     "success",
			id:       "123456789",
			mockTask: mockTask,
			mockErr:  nil,
			wantErr:  false,
		},
		{
			name:     "api error",
			id:       "123456789",
			mockTask: nil,
			mockErr:  errors.New("api error"),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock client
			client := NewMockClient(func(req *http.Request) (*http.Response, error) {
				if tt.mockErr != nil {
					return nil, tt.mockErr
				}
				if tt.mockTask == nil {
					return MockResponse(404, nil), nil
				}
				return MockResponse(200, tt.mockTask), nil
			})

			// Call the method
			task, err := client.GetTask(context.Background(), tt.id)

			// Check error
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, task)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, task)
				assert.Equal(t, tt.mockTask.ID, task.ID)
				assert.Equal(t, tt.mockTask.Content, task.Content)
			}
		})
	}
}

func TestCreateTask(t *testing.T) {
	// モックタスクを取得
	mockTask := MockTask()

	tests := []struct {
		name     string
		req      CreateTaskRequest
		mockTask *Task
		mockErr  error
		wantErr  bool
	}{
		{
			name: "success",
			req: CreateTaskRequest{
				Content:   "Test Task",
				ProjectID: "123456789",
			},
			mockTask: mockTask,
			mockErr:  nil,
			wantErr:  false,
		},
		{
			name: "api error",
			req: CreateTaskRequest{
				Content: "Test Task",
			},
			mockTask: nil,
			mockErr:  errors.New("api error"),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock client
			client := NewMockClient(func(req *http.Request) (*http.Response, error) {
				if tt.mockErr != nil {
					return nil, tt.mockErr
				}
				if tt.mockTask == nil {
					return MockResponse(400, nil), nil
				}
				return MockResponse(200, tt.mockTask), nil
			})

			// Call the method
			task, err := client.CreateTask(context.Background(), tt.req)

			// Check error
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, task)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, task)
				assert.Equal(t, tt.mockTask.ID, task.ID)
				assert.Equal(t, tt.mockTask.Content, task.Content)
			}
		})
	}
}

func TestUpdateTask(t *testing.T) {
	// モックタスクを取得
	mockTask := MockTask()

	tests := []struct {
		name     string
		id       string
		req      UpdateTaskRequest
		mockTask *Task
		mockErr  error
		wantErr  bool
	}{
		{
			name: "success",
			id:   "123456789",
			req: UpdateTaskRequest{
				Content: "Updated Task",
			},
			mockTask: mockTask,
			mockErr:  nil,
			wantErr:  false,
		},
		{
			name: "api error",
			id:   "123456789",
			req: UpdateTaskRequest{
				Content: "Updated Task",
			},
			mockTask: nil,
			mockErr:  errors.New("api error"),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock client
			client := NewMockClient(func(req *http.Request) (*http.Response, error) {
				if tt.mockErr != nil {
					return nil, tt.mockErr
				}
				if tt.mockTask == nil {
					return MockResponse(400, nil), nil
				}
				return MockResponse(200, tt.mockTask), nil
			})

			// Call the method
			task, err := client.UpdateTask(context.Background(), tt.id, tt.req)

			// Check error
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, task)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, task)
				assert.Equal(t, tt.mockTask.ID, task.ID)
				assert.Equal(t, tt.mockTask.Content, task.Content)
			}
		})
	}
}

func TestCloseTask(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		mockErr error
		wantErr bool
	}{
		{
			name:    "success",
			id:      "123456789",
			mockErr: nil,
			wantErr: false,
		},
		{
			name:    "api error",
			id:      "123456789",
			mockErr: errors.New("api error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock client
			client := NewMockClient(func(req *http.Request) (*http.Response, error) {
				if tt.mockErr != nil {
					return nil, tt.mockErr
				}
				return MockResponse(204, nil), nil
			})

			// Call the method
			err := client.CloseTask(context.Background(), tt.id)

			// Check error
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeleteTask(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		mockErr error
		wantErr bool
	}{
		{
			name:    "success",
			id:      "123456789",
			mockErr: nil,
			wantErr: false,
		},
		{
			name:    "api error",
			id:      "123456789",
			mockErr: errors.New("api error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock client
			client := NewMockClient(func(req *http.Request) (*http.Response, error) {
				if tt.mockErr != nil {
					return nil, tt.mockErr
				}
				return MockResponse(204, nil), nil
			})

			// Call the method
			err := client.DeleteTask(context.Background(), tt.id)

			// Check error
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetProjects(t *testing.T) {
	// モックプロジェクトのデリファレンス
	mockProject := *MockProject()

	tests := []struct {
		name         string
		mockProjects []Project
		mockErr      error
		wantErr      bool
	}{
		{
			name:         "success",
			mockProjects: []Project{mockProject, mockProject},
			mockErr:      nil,
			wantErr:      false,
		},
		{
			name:         "empty response",
			mockProjects: []Project{},
			mockErr:      nil,
			wantErr:      false,
		},
		{
			name:         "api error",
			mockProjects: nil,
			mockErr:      errors.New("api error"),
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock client
			client := NewMockClient(func(req *http.Request) (*http.Response, error) {
				if tt.mockErr != nil {
					return nil, tt.mockErr
				}
				// Return paginated response
				paginatedResp := MockPaginatedProjects(tt.mockProjects)
				return MockResponse(200, paginatedResp), nil
			})

			// Call the method
			projects, err := client.GetProjects(context.Background())

			// Check error
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.mockProjects), len(projects))
				if len(tt.mockProjects) > 0 {
					assert.Equal(t, tt.mockProjects[0].ID, projects[0].ID)
					assert.Equal(t, tt.mockProjects[0].Name, projects[0].Name)
				}
			}
		})
	}
}

func TestGetProject(t *testing.T) {
	// モックプロジェクトを取得
	mockProject := MockProject()

	tests := []struct {
		name        string
		id          string
		mockProject *Project
		mockErr     error
		wantErr     bool
	}{
		{
			name:        "success",
			id:          "987654321",
			mockProject: mockProject,
			mockErr:     nil,
			wantErr:     false,
		},
		{
			name:        "api error",
			id:          "987654321",
			mockProject: nil,
			mockErr:     errors.New("api error"),
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock client
			client := NewMockClient(func(req *http.Request) (*http.Response, error) {
				if tt.mockErr != nil {
					return nil, tt.mockErr
				}
				if tt.mockProject == nil {
					return MockResponse(404, nil), nil
				}
				return MockResponse(200, tt.mockProject), nil
			})

			// Call the method
			project, err := client.GetProject(context.Background(), tt.id)

			// Check error
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, project)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, project)
				assert.Equal(t, tt.mockProject.ID, project.ID)
				assert.Equal(t, tt.mockProject.Name, project.Name)
			}
		})
	}
}
