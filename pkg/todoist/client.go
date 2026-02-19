package todoist

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	// TodoistAPIBaseURL is the base URL for Todoist API
	TodoistAPIBaseURL = "https://api.todoist.com/api/v1"
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

// ClientOption is a function that configures a Client
type ClientOption func(*Client)

// WithTimeout sets the timeout for the HTTP client
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

// WithBaseURL sets the base URL for the Todoist API
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// WithLogger sets the logger for the client
func WithLogger(logger *logrus.Logger) ClientOption {
	return func(c *Client) {
		c.logger = logger
	}
}

// NewClient creates a new Todoist API client
func NewClient(token string, options ...ClientOption) *Client {
	if token == "" {
		token = os.Getenv("TODOIST_API_TOKEN")
	}

	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	client := &Client{
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 100,
				IdleConnTimeout:     90 * time.Second,
			},
		},
		token:   token,
		baseURL: TodoistAPIBaseURL,
		logger:  logger,
	}

	// Apply options
	for _, option := range options {
		option(client)
	}

	return client
}

// doRequest performs an HTTP request and returns the response
func (c *Client) doRequest(ctx context.Context, method, endpoint string, body io.Reader) (*http.Response, error) {
	// Create request
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authorization header
	req.Header.Add("Authorization", "Bearer "+c.token)

	// Add content type for requests with body
	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	// Execute the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	return resp, nil
}

// processResponse processes the HTTP response and handles errors
func (c *Client) processResponse(resp *http.Response, expectedStatus int) ([]byte, error) {
	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.logger.WithError(err).Error("Error closing response body")
		}
	}()

	// Check response status
	if resp.StatusCode != expectedStatus {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// For status codes that don't return content (like 204 No Content)
	if expectedStatus == http.StatusNoContent {
		return nil, nil
	}

	// Read response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return bodyBytes, nil
}
