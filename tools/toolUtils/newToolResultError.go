package toolUtils

import "github.com/mark3labs/mcp-go/mcp"

// NewToolResultError creates a new CallToolResult with an error message.
// Any errors that originate from the tool SHOULD be reported inside the result object.
// this function is not yet in mcp-go 0.17.0, but already in main branch of mcp-go
// copy from main branch temporarily, can be remove later
func NewToolResultError(text string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: text,
			},
		},
		IsError: true,
	}
}
