// Package readopenapidocument will read
package readopenapidocument

import (
	"context"
	"fmt"
	openapi "mcp-api-tester/openAPI"
	"mcp-api-tester/tools"
	"mcp-api-tester/tools/toolUtils"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Param provide param for readOpenAPIDocument
// it will also be parse into tools description and mount to mcp server
type Param struct {
	OpenAPIPath string `json:"openAPIPath" jsonschema:"required,description=The path that lead to OpenAPI.yaml or OpenAPI.json or any OpenAPI file, try Absolute path if not work"`
}

func readOpenAPIDocument(_ context.Context, args Param) (*mcp.CallToolResult, error) {
	_, err := openapi.ReadFromPath(args.OpenAPIPath)

	if err != nil {
		return mcp.NewToolResultText("error"), err
	}

	return mcp.NewToolResultText("success"), nil
}

// ReadOpenAPIDocumentTool can register readOpenAPIDocument to MCP Server
var ReadOpenAPIDocumentTool = toolUtils.MustTool(
	tools.ReadOpenAPIDocument,
	fmt.Sprintf("%s will read OpenAPI file by given Path, Please run this tool first to load OpenAPI file before using other tools", tools.ReadOpenAPIDocument),
	readOpenAPIDocument,
)

// AddReadOpenAPIDocumentTool can register readOpenAPIDocument to MCP Server
func AddReadOpenAPIDocumentTool(mcp *server.MCPServer) {
	ReadOpenAPIDocumentTool.Register(mcp)
}
