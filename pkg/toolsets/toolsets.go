package toolsets

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ToolHandlerFunc is a handler function for a tool call.
type ToolHandlerFunc func(context.Context, *mcp.CallToolRequest) (*mcp.CallToolResult, error)

// ServerTool is a tool definition bound to a handler.
type ServerTool struct {
	Tool    mcp.Tool
	Handler ToolHandlerFunc
}

// NewServerTool creates a new server tool
func NewServerTool(tool mcp.Tool, handler ToolHandlerFunc) ServerTool {
	return ServerTool{Tool: tool, Handler: handler}
}

// Toolset represents a group of related tools
type Toolset struct {
	Name        string
	Description string
	Enabled     bool
	readOnly    bool
	writeTools  []ServerTool
	readTools   []ServerTool
}

// GetActiveTools returns all active tools in the toolset
func (t *Toolset) GetActiveTools() []ServerTool {
	if t.Enabled {
		if t.readOnly {
			return t.readTools
		}
		return append(t.readTools, t.writeTools...)
	}
	return nil
}

// GetAvailableTools returns all available tools in the toolset
func (t *Toolset) GetAvailableTools() []ServerTool {
	if t.readOnly {
		return t.readTools
	}
	return append(t.readTools, t.writeTools...)
}

// RegisterTools registers all tools in the toolset with the MCP server
func (t *Toolset) RegisterTools(s *mcp.Server) {
	if !t.Enabled {
		return
	}
	for _, tool := range t.readTools {
		s.AddTool(&tool.Tool, mcp.ToolHandler(tool.Handler))
	}
	if !t.readOnly {
		for _, tool := range t.writeTools {
			s.AddTool(&tool.Tool, mcp.ToolHandler(tool.Handler))
		}
	}
}

// SetReadOnly sets the toolset to read-only mode
func (t *Toolset) SetReadOnly() {
	t.readOnly = true
}

// AddWriteTools adds write tools to the toolset
func (t *Toolset) AddWriteTools(tools ...ServerTool) *Toolset {
	for _, tool := range tools {
		if tool.Tool.Annotations != nil && tool.Tool.Annotations.ReadOnlyHint {
			panic(fmt.Sprintf("tool (%s) is incorrectly annotated as read-only", tool.Tool.Name))
		}
	}
	if !t.readOnly {
		t.writeTools = append(t.writeTools, tools...)
	}
	return t
}

// AddReadTools adds read-only tools to the toolset
func (t *Toolset) AddReadTools(tools ...ServerTool) *Toolset {
	for _, tool := range tools {
		if tool.Tool.Annotations == nil || !tool.Tool.Annotations.ReadOnlyHint {
			panic(fmt.Sprintf("tool (%s) must be annotated as read-only", tool.Tool.Name))
		}
		tool.Tool.Annotations = &mcp.ToolAnnotations{
			ReadOnlyHint: true,
			Title:        tool.Tool.Annotations.Title,
		}
	}
	t.readTools = append(t.readTools, tools...)
	return t
}

// ToolsetGroup represents a group of toolsets
type ToolsetGroup struct {
	Toolsets     map[string]*Toolset
	everythingOn bool
	readOnly     bool
}

// NewToolsetGroup creates a new toolset group
func NewToolsetGroup(readOnly bool) *ToolsetGroup {
	return &ToolsetGroup{
		Toolsets:     make(map[string]*Toolset),
		everythingOn: false,
		readOnly:     readOnly,
	}
}

// AddToolset adds a toolset to the group
func (tg *ToolsetGroup) AddToolset(ts *Toolset) {
	if tg.readOnly {
		ts.SetReadOnly()
	}
	tg.Toolsets[ts.Name] = ts
}

// NewToolset creates a new toolset
func NewToolset(name string, description string) *Toolset {
	return &Toolset{
		Name:        name,
		Description: description,
		Enabled:     false,
		readOnly:    false,
	}
}

// IsEnabled checks if a toolset is enabled
func (tg *ToolsetGroup) IsEnabled(name string) bool {
	if tg.everythingOn {
		return true
	}

	toolset, exists := tg.Toolsets[name]
	if !exists {
		return false
	}
	return toolset.Enabled
}

// EnableToolsets enables multiple toolsets
func (tg *ToolsetGroup) EnableToolsets(names []string) error {
	for _, name := range names {
		if name == "all" {
			tg.everythingOn = true
			break
		}
		err := tg.EnableToolset(name)
		if err != nil {
			return err
		}
	}

	if tg.everythingOn {
		for name := range tg.Toolsets {
			err := tg.EnableToolset(name)
			if err != nil {
				return err
			}
		}
		return nil
	}
	return nil
}

// EnableToolset enables a specific toolset
func (tg *ToolsetGroup) EnableToolset(name string) error {
	toolset, exists := tg.Toolsets[name]
	if !exists {
		return fmt.Errorf("toolset %s does not exist", name)
	}
	toolset.Enabled = true
	tg.Toolsets[name] = toolset
	return nil
}

// RegisterTools registers all tools in the group with the MCP server
func (tg *ToolsetGroup) RegisterTools(s *mcp.Server) {
	for _, toolset := range tg.Toolsets {
		toolset.RegisterTools(s)
	}
}
