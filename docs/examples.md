# Todoist MCP Server Usage Examples

This document provides practical examples of how to use the Todoist MCP Server with different MCP clients.

## JavaScript Examples

### Setup

```javascript
// Import your preferred MCP client library
const { MCPClient } = require('mcp-client');

// Initialize the client
const mcpClient = new MCPClient({
  serverUrl: 'http://localhost:8080', // For HTTP mode
  // Or for stdio mode, configure accordingly
});
```

### Task Management Examples

#### Getting Today's Tasks

```javascript
async function getTodaysTasks() {
  try {
    const result = await mcpClient.callTool("todoist_get_tasks", {
      filter: "today"
    });
    
    console.log("Today's tasks:");
    result.tasks.forEach(task => {
      console.log(`- ${task.content} (Priority: ${task.priority})`);
    });
    
    return result.tasks;
  } catch (error) {
    console.error("Error fetching today's tasks:", error);
  }
}
```

#### Creating a Task with Due Date

```javascript
async function createTaskWithDueDate(content, dueString, priority = 4) {
  try {
    const result = await mcpClient.callTool("todoist_create_task", {
      content: content,
      dueString: dueString,
      priority: priority
    });
    
    console.log(`Task created: ${result.task.content}`);
    console.log(`Due: ${result.task.due?.string}`);
    
    return result.task;
  } catch (error) {
    console.error("Error creating task:", error);
  }
}

// Example usage
createTaskWithDueDate(
  "Prepare presentation for meeting", 
  "tomorrow at 9am", 
  1 // Urgent priority
);
```

#### Updating a Task

```javascript
async function updateTaskPriority(taskId, newPriority) {
  try {
    const result = await mcpClient.callTool("todoist_update_task", {
      id: taskId,
      priority: newPriority
    });
    
    console.log(`Task updated: ${result.task.content}`);
    console.log(`New priority: ${result.task.priority}`);
    
    return result.task;
  } catch (error) {
    console.error("Error updating task:", error);
  }
}
```

#### Completing a Task

```javascript
async function completeTask(taskId) {
  try {
    const result = await mcpClient.callTool("todoist_close_task", {
      id: taskId
    });
    
    if (result.success) {
      console.log(`Task ${taskId} marked as completed`);
    }
    
    return result.success;
  } catch (error) {
    console.error("Error completing task:", error);
  }
}
```

### Project Management Examples

#### Getting All Projects

```javascript
async function getAllProjects() {
  try {
    const result = await mcpClient.callTool("todoist_get_projects", {});
    
    console.log("Your projects:");
    result.projects.forEach(project => {
      console.log(`- ${project.name} (${project.id})`);
    });
    
    return result.projects;
  } catch (error) {
    console.error("Error fetching projects:", error);
  }
}
```

#### Getting Tasks for a Specific Project

```javascript
async function getProjectTasks(projectId) {
  try {
    const result = await mcpClient.callTool("todoist_get_tasks", {
      projectId: projectId
    });
    
    console.log(`Tasks in project ${projectId}:`);
    result.tasks.forEach(task => {
      console.log(`- ${task.content}`);
    });
    
    return result.tasks;
  } catch (error) {
    console.error("Error fetching project tasks:", error);
  }
}
```

## Python Examples

### Setup

```python
# Import your preferred MCP client library
from mcp_client import MCPClient

# Initialize the client
mcp_client = MCPClient(
    server_url="http://localhost:8080"  # For HTTP mode
    # Or for stdio mode, configure accordingly
)
```

### Task Management Examples

#### Getting Today's Tasks

```python
def get_todays_tasks():
    try:
        result = mcp_client.call_tool("todoist_get_tasks", {
            "filter": "today"
        })
        
        print("Today's tasks:")
        for task in result["tasks"]:
            print(f"- {task['content']} (Priority: {task['priority']})")
        
        return result["tasks"]
    except Exception as e:
        print(f"Error fetching today's tasks: {e}")
```

#### Creating a Task with Due Date

```python
def create_task_with_due_date(content, due_string, priority=4):
    try:
        result = mcp_client.call_tool("todoist_create_task", {
            "content": content,
            "dueString": due_string,
            "priority": priority
        })
        
        task = result["task"]
        print(f"Task created: {task['content']}")
        if "due" in task and task["due"]:
            print(f"Due: {task['due'].get('string', '')}")
        
        return task
    except Exception as e:
        print(f"Error creating task: {e}")

# Example usage
create_task_with_due_date(
    "Prepare presentation for meeting", 
    "tomorrow at 9am", 
    1  # Urgent priority
)
```

#### Completing a Task

```python
def complete_task(task_id):
    try:
        result = mcp_client.call_tool("todoist_close_task", {
            "id": task_id
        })
        
        if result.get("success"):
            print(f"Task {task_id} marked as completed")
        
        return result.get("success")
    except Exception as e:
        print(f"Error completing task: {e}")
```

## Go Examples

### Setup

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	// Initialize the MCP client
	client, err := mcp.NewClient("http://localhost:8080") // For HTTP mode
	if err != nil {
		log.Fatalf("Failed to create MCP client: %v", err)
	}

	// Example function calls
	getTodaysTasks(client)
}
```

### Task Management Examples

#### Getting Today's Tasks

```go
func getTodaysTasks(client *mcp.Client) {
	ctx := context.Background()
	
	// Prepare parameters
	params := map[string]interface{}{
		"filter": "today",
	}
	
	// Call the tool
	result, err := client.CallTool(ctx, "todoist_get_tasks", params)
	if err != nil {
		log.Fatalf("Error calling todoist_get_tasks: %v", err)
	}
	
	// Parse the response
	var response struct {
		Tasks []struct {
			ID      string `json:"id"`
			Content string `json:"content"`
			Priority int    `json:"priority"`
		} `json:"tasks"`
	}
	
	if err := json.Unmarshal([]byte(result.Text), &response); err != nil {
		log.Fatalf("Error parsing response: %v", err)
	}
	
	// Display the tasks
	fmt.Println("Today's tasks:")
	for _, task := range response.Tasks {
		fmt.Printf("- %s (Priority: %d)\n", task.Content, task.Priority)
	}
}
```

#### Creating a Task

```go
func createTask(client *mcp.Client, content, dueString string, priority int) {
	ctx := context.Background()
	
	// Prepare parameters
	params := map[string]interface{}{
		"content":   content,
		"dueString": dueString,
		"priority":  priority,
	}
	
	// Call the tool
	result, err := client.CallTool(ctx, "todoist_create_task", params)
	if err != nil {
		log.Fatalf("Error calling todoist_create_task: %v", err)
	}
	
	// Parse the response
	var response struct {
		Task struct {
			ID      string `json:"id"`
			Content string `json:"content"`
			Due     struct {
				String string `json:"string"`
			} `json:"due"`
		} `json:"task"`
	}

	if err := json.Unmarshal([]byte(result.Text), &response); err != nil {
		log.Fatalf("Error parsing response: %v", err)
	}

	// Display the created task
	fmt.Printf("Task created: %s\n", response.Task.Content)
	fmt.Printf("Due: %s\n", response.Task.Due.String)
}
```

## Integration with AI Assistants

MCP allows AI assistants to interact with Todoist through natural language. Here's how an AI assistant might use the Todoist MCP Server:

```
User: "I need to prepare for tomorrow's meeting."

AI: "I'll help you create a task for that. Let me add it to your Todoist."

[AI uses todoist_create_task behind the scenes]

AI: "I've added 'Prepare for meeting' to your Todoist with a due date of tomorrow. Would you like me to add any specific details to this task?"

User: "Yes, please add that I need to prepare the quarterly sales report."

[AI uses todoist_update_task behind the scenes]

AI: "I've updated the task with the details. Your task 'Prepare for meeting' now includes 'Need to prepare the quarterly sales report' in the description, due tomorrow."
```

This integration enables AI assistants to manage tasks and projects on behalf of users through natural conversation, without requiring users to directly interact with the Todoist interface.
