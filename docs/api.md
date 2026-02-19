# Todoist MCP Server API Documentation

## Overview

This document describes the API tools provided by the Todoist MCP Server. The server implements the Model Context Protocol (MCP) to enable AI assistants to interact with Todoist's task and project management features.

## Authentication

To use the Todoist API, you need a personal API token from Todoist. This token should be set as the environment variable `TODOIST_API_TOKEN` or specified in a configuration file.

```go
// Example of retrieving the authentication token
token := os.Getenv("TODOIST_API_TOKEN")
if token == "" {
    // Implement alternative methods, such as reading from a configuration file
}
```

## MCP Tools

### Task Management

#### 1. `todoist_get_tasks`

Retrieves a list of tasks with optional filtering.

**Parameters:**
- `projectId` (string, optional): Filter tasks by project ID
- `filter` (string, optional): Todoist filter query using the Todoist filter syntax

**Response:**
```json
{
  "tasks": [
    {
      "id": "2995104339",
      "content": "Buy Milk",
      "description": "",
      "project_id": "2203306141",
      "parent_id": null,
      "priority": 1,
      "due": {
        "date": "2016-09-01",
        "is_recurring": false,
        "datetime": "2016-09-01T12:00:00.000000Z",
        "string": "tomorrow at 12",
        "timezone": "Europe/Moscow"
      }
    }
  ]
}
```

**Usage Example:**
```javascript
// Example MCP client call to get today's tasks
const result = await mcpClient.callTool("todoist_get_tasks", {
  filter: "today"
});
console.log(result.tasks);

// Example to get tasks from a specific project
const projectTasks = await mcpClient.callTool("todoist_get_tasks", {
  projectId: "2203306141"
});
console.log(projectTasks.tasks);
```

#### 2. `todoist_get_task`

Retrieves a specific task by its ID.

**Parameters:**
- `id` (string, required): The unique identifier of the task to retrieve

**Response:**
```json
{
  "task": {
    "id": "2995104339",
    "content": "Buy Milk",
    "description": "",
    "project_id": "2203306141",
    "parent_id": null,
    "priority": 1,
    "due": {
      "date": "2016-09-01",
      "is_recurring": false,
      "datetime": "2016-09-01T12:00:00.000000Z",
      "string": "tomorrow at 12",
      "timezone": "Europe/Moscow"
    }
  }
}
```

**Usage Example:**
```javascript
// Example MCP client call to get a specific task
const result = await mcpClient.callTool("todoist_get_task", {
  id: "2995104339"
});
console.log(result.task);
```

#### 3. `todoist_create_task`

Creates a new task.

**Parameters:**
- `content` (string, required): The content of the task
- `description` (string, optional): Detailed description or notes for the task
- `projectId` (string, optional): Project ID to assign the task to
- `parentId` (string, optional): Parent task ID for creating subtasks
- `order` (integer, optional): Order value for positioning the task
- `priority` (integer, optional): Task priority: 1 (normal), 2 (medium), 3 (high), 4 (urgent)
- `dueString` (string, optional): Due date in natural language, e.g., 'today', 'tomorrow'
- `dueDate` (string, optional): Due date in YYYY-MM-DD format
- `dueDatetime` (string, optional): Due date and time in RFC3339 format

**Response:**
```json
{
  "task": {
    "id": "2995104339",
    "content": "Buy Milk",
    "description": "",
    "project_id": "2203306141",
    "parent_id": null,
    "priority": 1,
    "due": {
      "date": "2016-09-01",
      "is_recurring": false,
      "datetime": "2016-09-01T12:00:00.000000Z",
      "string": "tomorrow at 12",
      "timezone": "Europe/Moscow"
    }
  }
}
```

**Usage Example:**
```javascript
// Example MCP client call to create a new task
const result = await mcpClient.callTool("todoist_create_task", {
  content: "Buy groceries",
  description: "Need to buy milk, eggs, and bread",
  projectId: "2203306141",
  priority: 2,
  dueString: "tomorrow at 10am"
});
console.log(result.task);
```

#### 4. `todoist_update_task`

Updates an existing task.

**Parameters:**
- `id` (string, required): The unique identifier of the task to update
- `content` (string, optional): The new content of the task
- `description` (string, optional): Detailed description or notes for the task
- `priority` (integer, optional): Task priority: 1 (normal), 2 (medium), 3 (high), 4 (urgent)
- `dueString` (string, optional): Due date in natural language
- `dueDate` (string, optional): Due date in YYYY-MM-DD format
- `dueDatetime` (string, optional): Due date and time in RFC3339 format

**Response:**
```json
{
  "task": {
    "id": "2995104339",
    "content": "Buy Milk and Bread",
    "description": "From the organic store",
    "project_id": "2203306141",
    "parent_id": null,
    "priority": 2,
    "due": {
      "date": "2016-09-02",
      "is_recurring": false,
      "datetime": "2016-09-02T10:00:00.000000Z",
      "string": "tomorrow at 10",
      "timezone": "Europe/Moscow"
    }
  }
}
```

**Usage Example:**
```javascript
// Example MCP client call to update a task
const result = await mcpClient.callTool("todoist_update_task", {
  id: "2995104339",
  content: "Buy Milk and Bread",
  description: "From the organic store",
  priority: 2,
  dueString: "tomorrow at 10am"
});
console.log(result.task);
```

#### 5. `todoist_close_task`

Marks a task as completed.

**Parameters:**
- `id` (string, required): The unique identifier of the task to mark as completed

**Response:**
```json
{
  "success": true
}
```

**Usage Example:**
```javascript
// Example MCP client call to mark a task as completed
const result = await mcpClient.callTool("todoist_close_task", {
  id: "2995104339"
});
console.log(result.success);
```

#### 6. `todoist_delete_task`

Deletes a task.

**Parameters:**
- `id` (string, required): The unique identifier of the task to delete

**Response:**
```json
{
  "success": true
}
```

**Usage Example:**
```javascript
// Example MCP client call to delete a task
const result = await mcpClient.callTool("todoist_delete_task", {
  id: "2995104339"
});
console.log(result.success);
```

### Project Management

#### 1. `todoist_get_projects`

Retrieves a list of all projects.

**Parameters:** None

**Response:**
```json
{
  "projects": [
    {
      "id": "2203306141",
      "name": "Inbox",
      "color": "grey",
      "parent_id": null,
      "child_order": 0,
      "is_shared": false,
      "is_favorite": false,
      "inbox_project": true,
      "view_style": "list"
    },
    {
      "id": "2203306142",
      "name": "Personal",
      "color": "blue",
      "parent_id": null,
      "child_order": 1,
      "is_shared": false,
      "is_favorite": true,
      "inbox_project": false,
      "view_style": "list"
    }
  ]
}
```

**Usage Example:**
```javascript
// Example MCP client call to get all projects
const result = await mcpClient.callTool("todoist_get_projects", {});
console.log(result.projects);
```

#### 2. `todoist_get_project`

Retrieves a specific project by its ID.

**Parameters:**
- `id` (string, required): The unique identifier of the project to retrieve

**Response:**
```json
{
  "project": {
    "id": "2203306142",
    "name": "Personal",
    "color": "blue",
    "parent_id": null,
    "child_order": 1,
    "is_shared": false,
    "is_favorite": true,
    "inbox_project": false,
    "view_style": "list"
  }
}
```

**Usage Example:**
```javascript
// Example MCP client call to get a specific project
const result = await mcpClient.callTool("todoist_get_project", {
  id: "2203306142"
});
console.log(result.project);
```

## Error Handling

All tools return appropriate error messages when operations fail. Error responses follow this format:

```json
{
  "error": {
    "message": "Error message describing what went wrong",
    "code": "ERROR_CODE"
  }
}
```

Common error scenarios include:
- Invalid parameters
- Resource not found
- Authentication failures
- API rate limiting
- Network errors

## Rate Limiting

The Todoist API has rate limits. The server implements appropriate handling to respect these limits and provides meaningful error messages when limits are exceeded.

## References

- [Todoist API v1 Documentation](https://developer.todoist.com/api/v1/)
- [Model Context Protocol Specification](https://github.com/mcp-sh/mcp)
