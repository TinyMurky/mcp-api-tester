// Package listallapifromdocument will read Open API file
// and return { "name": "getUser", "method": "GET", "url": "/users/{id}" }
package listallapifromdocument

import (
	"context"
	"encoding/json"
	"fmt"
	openapi "mcp-api-tester/openAPI"
	"mcp-api-tester/tools"
	"mcp-api-tester/tools/toolUtils"

	"github.com/mark3labs/mcp-go/server"
)

// Param provide param for readOpenAPIDocument
// it will also be parse into tools description and mount to mcp server
type Param struct{}

func listAllAPIFromDocument(_ context.Context, _ Param) (string, error) {
	if openapi.OpenAPIPointer == nil {
		return "", fmt.Errorf("OpenAPI file hasn't read, please use %q tool to read file first ", tools.ReadOpenAPIDocument)
	}

	simplifyAPIs := openapi.OpenAPIPointer.ListAllAPIFromDocument()

	simplifyAPIBytes, err := json.Marshal(simplifyAPIs)

	if err != nil {
		return "", fmt.Errorf("Convert simplifyAPIs to json failed, error: %w", err)
	}

	return string(simplifyAPIBytes), nil
}

// ListAllAPIFromDocumentTool can register readOpenAPIDocument to MCP Server
var ListAllAPIFromDocumentTool = toolUtils.MustTool(
	tools.ListAllAPIFromDocument,
	fmt.Sprintf("%s will list all api and method from openAPI, Please use %q to tool OpenAPI file first", tools.ListAllAPIFromDocument, tools.ReadOpenAPIDocument),
	listAllAPIFromDocument,
)

// AddListAllAPIFromDocumentTool can register listAllAPIFromDocument to MCP Server
func AddListAllAPIFromDocumentTool(mcp *server.MCPServer) {
	ListAllAPIFromDocumentTool.Register(mcp)
}
