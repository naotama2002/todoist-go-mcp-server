# Todoist MCP Server

A Model Context Protocol (MCP) server that provides Todoist API integration for AI assistants.

## Overview

The Todoist MCP Server allows AI assistants to interact with Todoist, enabling them to manage tasks and projects on behalf of users. This server implements the [Model Context Protocol (MCP)](https://github.com/mcp-sh/mcp), providing a standardized interface for AI systems to access Todoist functionality.

## Features

- **Task Management**
  - Get filter rules and examples for task filtering
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

#### `todoist_get_task_filter_rules`

Get filter rules and examples for Todoist task filters. This tool helps translate natural language queries into Todoist filter syntax for the `todoist_get_tasks` tool.

Parameters: None

Example response (Markdown format):
```markdown
# Introduction to Filters

Filters in Todoist are custom views that display tasks based on specific criteria. You can filter tasks by name, date, project, label, priority, creation date, and more.

## Creating Filters

1. Select "Filters & Labels" in the sidebar
2. Click the add icon next to Filters
3. Enter a name for your filter and change its color (optional)
4. Enter your filter query
5. Click "Add" to save your filter

## Filter Symbols

| Symbol | Meaning | Example |
|--------|---------|--------|
| \| | OR | today \| overdue |
| & | AND | today & p1 |
| ! | NOT | !subtask |
| () | Priority processing | (today \| overdue) & #work |

## Advanced Queries

### Keyword-Based Filters

| Description | Query |
|-------------|-------|
| Tasks containing "meeting" | search: meeting |
| Tasks containing "meeting" scheduled for today | search: meeting & today |

### Date-Based Filters

| Description | Query |
|-------------|-------|
| Tasks for a specific date | date: jan 3 |
| Tasks before a specific date | before: may 5 |
| Tasks with no date | no date |

## Useful Filter Examples

| Description | Query |
|-------------|-------|
| Overdue or today's tasks in "Work" project | (today \| overdue) & #work |
| Tasks with no date | no date |
| Tasks with @waiting label in the next 7 days | 7 days & @waiting |
```
```

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

### Download from GitHub Releases

The easiest way to get started is to download a pre-built binary from the [GitHub Releases page](https://github.com/naotama2002/todoist-go-mcp-server/releases):

1. Go to the [Releases page](https://github.com/naotama2002/todoist-go-mcp-server/releases) and download the latest release for your platform:
   - Linux: `todoist-mcp-server_Linux_x86_64.tar.gz` or `todoist-mcp-server_Linux_arm64.tar.gz`
   - macOS: `todoist-mcp-server_Darwin_x86_64.tar.gz` or `todoist-mcp-server_Darwin_arm64.tar.gz`
   - Windows: `todoist-mcp-server_Windows_x86_64.zip`

2. Extract the archive to get the `todoist-mcp-server` binary:
   ```bash
   # For Linux/macOS
   tar -xzf todoist-mcp-server_*_*.tar.gz
   
   # For Windows
   # Extract the zip file using Windows Explorer or a tool like 7-Zip
   ```

3. Make the binary executable (Linux/macOS only):
   ```bash
   chmod +x todoist-mcp-server
   ```

4. Optionally, move the binary to a location in your PATH:
   ```bash
   # Linux/macOS
   sudo mv todoist-mcp-server /usr/local/bin/
   
   # Windows
   # Move the .exe file to a location in your PATH
   ```

### Build from source

Alternatively, you can build the binary from source:

```bash
# Build the binary
make build
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
      "name": "todoist_get_task_filter_rules",
      "description": "Get filter rules and examples for Todoist task filters"
    },
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

Claude Desktop / Cline, Roo Cline / Windsurf

```json
{
  "mcpServers": {
    "github": {
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

VSCode
```
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

1. In Claude Desktop, go to Settings > Tools > Add Tool, and paste the appropriate JSON configuration.

2. Save the configuration and restart Claude Desktop if necessary.

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
make test
```
