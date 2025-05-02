# Todoist MCP Server

A Model Context Protocol (MCP) server that provides Todoist API integration for AI assistants.

## Overview

The Todoist MCP Server allows AI assistants to interact with Todoist, enabling them to manage tasks and projects on behalf of users. This server implements the [Model Context Protocol (MCP)](https://github.com/mcp-sh/mcp), providing a standardized interface for AI systems to access Todoist functionality.

## Features

- **Task Management**
  - Get tasks with filtering options
  - Get task details
  - Create new tasks
  - Update existing tasks
  - Mark tasks as completed
  - Delete tasks

- **Project Management**
  - Get all projects
  - Get project details

## Installation

### Prerequisites

- Go 1.21 or later
- Todoist API token

### Setup

1. Clone the repository:

```bash
git clone https://github.com/naotama2002/todoist-go-mcp-server.git
cd todoist-go-mcp-server
```

2. Install dependencies:

```bash
go mod download
```

3. Set up your Todoist API token:

Create a `.envrc` file in the root directory with the following content:

```bash
export TODOIST_API_TOKEN=your_todoist_api_token
```

If you're using [direnv](https://direnv.net/), run:

```bash
direnv allow
```

Alternatively, you can set the environment variable directly:

```bash
export TODOIST_API_TOKEN=your_todoist_api_token
```

## Usage

### Running the Server

#### HTTP Mode

Run the server in HTTP mode:

```bash
go run cmd/todoist-mcp-server/main.go --mode http --addr :8080
```

#### Standard I/O Mode

Run the server in stdio mode for integration with MCP clients:

```bash
go run cmd/todoist-mcp-server/main.go --mode stdio
```

### Command Line Options

The server supports the following command line options:

- `--mode <mode>`: Server mode, either 'http' or 'stdio' (default: "http")
  - `http`: Run as an HTTP server
  - `stdio`: Run using standard input/output for MCP communication
- `--addr <address>`: Address to listen on in HTTP mode (default: ":8080")
- `--token <token>`: Todoist API token (can also be set via TODOIST_API_TOKEN environment variable)

Examples:

```bash
# Run HTTP server on port 3000
go run cmd/todoist-mcp-server/main.go --mode http --addr :3000

# Run with a specific Todoist API token
go run cmd/todoist-mcp-server/main.go --mode http --token your_todoist_api_token

# Run in stdio mode
go run cmd/todoist-mcp-server/main.go --mode stdio
```

### Testing with the MCP Client

You can test the server using the included test client:

```bash
go run cmd/test-mcp-client/main.go
```

## Available Tools

### Task Management

#### `todoist_get_tasks`

Get a list of tasks with filtering options.

Parameters:
- `projectId` (string, optional): Filter tasks by project ID
- `filter` (string, optional): Todoist filter query using the Todoist filter syntax

Example:
```json
{
  "projectId": "2203306141",
  "filter": "today"
}
```

#### `todoist_get_task`

Get details of a specific task.

Parameters:
- `id` (string, required): The unique identifier of the task

Example:
```json
{
  "id": "2995104339"
}
```

#### `todoist_create_task`

Create a new task.

Parameters:
- `content` (string, required): The content of the task
- `description` (string, optional): Detailed description or notes for the task
- `projectId` (string, optional): Project ID to assign the task to
- `parentId` (string, optional): Parent task ID for creating subtasks
- `order` (integer, optional): Order value for positioning the task
- `priority` (integer, optional): Task priority: 4 (normal), 3 (medium), 2 (high), 1 (urgent)
- `dueString` (string, optional): Due date in natural language, e.g., 'today', 'tomorrow'
- `dueDate` (string, optional): Due date in YYYY-MM-DD format
- `dueDatetime` (string, optional): Due date and time in RFC3339 format

Example:
```json
{
  "content": "Buy groceries",
  "description": "Need to buy milk, eggs, and bread",
  "projectId": "2203306141",
  "priority": 2,
  "dueString": "tomorrow at 10am"
}
```

#### `todoist_update_task`

Update an existing task.

Parameters:
- `id` (string, required): The unique identifier of the task to update
- `content` (string, optional): The new content of the task
- `description` (string, optional): Detailed description or notes for the task
- `priority` (integer, optional): Task priority: 4 (normal), 3 (medium), 2 (high), 1 (urgent)
- `dueString` (string, optional): Due date in natural language
- `dueDate` (string, optional): Due date in YYYY-MM-DD format
- `dueDatetime` (string, optional): Due date and time in RFC3339 format

Example:
```json
{
  "id": "2995104339",
  "content": "Buy groceries and household items",
  "priority": 1
}
```

#### `todoist_close_task`

Mark a task as completed.

Parameters:
- `id` (string, required): The unique identifier of the task to mark as completed

Example:
```json
{
  "id": "2995104339"
}
```

#### `todoist_delete_task`

Delete a task.

Parameters:
- `id` (string, required): The unique identifier of the task to delete

Example:
```json
{
  "id": "2995104339"
}
```

### Project Management

#### `todoist_get_projects`

Get a list of all projects.

Parameters: None

#### `todoist_get_project`

Get details of a specific project.

Parameters:
- `id` (string, required): The unique identifier of the project

Example:
```json
{
  "id": "2203306141"
}
```

## Integration with Claude Desktop

To use the Todoist MCP Server with Claude Desktop, you need to add it to your Claude Desktop configuration.

### Build from source

You can use `go build` to build the binary:

```bash
# Build the binary
go build -o todoist-mcp-server ./cmd/todoist-mcp-server

# Move the binary to a location in your PATH (optional)
sudo mv todoist-mcp-server /usr/local/bin/
```

### Configuration

Add the Todoist MCP Server to your Claude Desktop configuration:

#### HTTP Mode

1. Run the Todoist MCP Server in HTTP mode:

```bash
todoist-mcp-server --mode http --addr :8080
```

2. Add the following configuration to your Claude Desktop settings:

```json
{
  "name": "Todoist",
  "description": "Manage Todoist tasks and projects",
  "endpoint": "http://localhost:8080",
  "tools": [
    {
      "name": "todoist_get_tasks",
      "description": "Get a list of tasks with filtering options"
    },
    {
      "name": "todoist_get_task",
      "description": "Get details of a specific task"
    },
    {
      "name": "todoist_create_task",
      "description": "Create a new task"
    },
    {
      "name": "todoist_update_task",
      "description": "Update an existing task"
    },
    {
      "name": "todoist_close_task",
      "description": "Mark a task as completed"
    },
    {
      "name": "todoist_delete_task",
      "description": "Delete a task"
    },
    {
      "name": "todoist_get_projects",
      "description": "Get a list of all projects"
    },
    {
      "name": "todoist_get_project",
      "description": "Get details of a specific project"
    }
  ]
}
```

#### Stdio Mode

Alternatively, you can configure Claude Desktop to directly execute the binary in stdio mode:

```json
{
  "mcp": {
    "servers": {
      "todoist": {
        "command": "/path/to/todoist-mcp-server",
        "args": ["--mode", "stdio"],
        "env": {
          "TODOIST_API_TOKEN": "<YOUR_TODOIST_API_TOKEN>"
        }
      }
    }
  }
}
```

Replace `/path/to/todoist-mcp-server` with the actual path to the binary and `<YOUR_TODOIST_API_TOKEN>` with your Todoist API token.

3. In Claude Desktop, go to Settings > Tools > Add Tool, and paste the appropriate JSON configuration.

4. Save the configuration and restart Claude Desktop if necessary.

Now you can ask Claude to manage your Todoist tasks and projects using natural language.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Todoist API](https://developer.todoist.com/rest/v2/)
- [Model Context Protocol](https://github.com/mcp-sh/mcp)

## Development

### Project Structure

```
todoist-go-mcp-server/
├── cmd/
│   ├── test-client/        # Test client for the Todoist API
│   ├── test-mcp-client/    # Test client for the MCP server
│   └── todoist-mcp-server/ # Main MCP server application
├── docs/                   # Documentation
├── pkg/
│   ├── log/                # Logging utilities
│   ├── todoist/            # Todoist API client and tools
│   └── toolsets/           # MCP toolset definitions
└── todo/                   # Implementation plans and notes
```

### Running Tests

```bash
go test ./...
```
