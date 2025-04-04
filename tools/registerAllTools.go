// Package tools collected all mcp tool,
// all tool need to implement Tool interface
package tools

import (
	"mcp-api-tester/tools/hello"

	"github.com/mark3labs/mcp-go/server"
)

// RegisterAllTools register all tools to server
func RegisterAllTools(srv *server.MCPServer) {
	hello := &hello.Hello{}
	hello.Register(srv)
}
