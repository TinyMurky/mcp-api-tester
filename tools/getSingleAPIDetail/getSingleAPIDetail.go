// Package getsingleapidetail will return specific data of
package getsingleapidetail

import (
	"context"
	"fmt"
	openapi "mcp-api-tester/openAPI"
	"mcp-api-tester/tools"
	toolutils "mcp-api-tester/tools/toolUtils"

	v3high "github.com/pb33f/libopenapi/datamodel/high/v3"

	"github.com/mark3labs/mcp-go/server"
)

// Param provide param for readOpenAPIDocument
// it will also be parse into tools description and mount to mcp server
type Param struct {
	URLPath string `json:"urlPath" jsonschema:"required,description=The url path is  the route url that you want to search"`
	Method  string `json:"method" jsonschema:"required,description=Http method that you want to check of certain url path,enum=get,enum=put,enum=post,enum=delete,enum=options,enum=head,enum=patch,enum=trace"`
}

func getSingleAPIDetail(_ context.Context, args Param) (*v3high.Operation, error) {
	if openapi.OpenAPIPointer == nil {
		return nil, fmt.Errorf("OpenAPI file hasn't read, please use %q tool to read file first ", tools.ReadOpenAPIDocument)
	}

	operation, err := openapi.OpenAPIPointer.GetOneAPIByPath(args.URLPath, args.Method)

	if err != nil {
		return nil, err
	}

	return operation, nil
}

// GetSingleAPIDetailTool can register getSingleAPIDetail to MCP Server
var GetSingleAPIDetailTool = toolutils.MustTool(
	tools.GetSingleAPIDetail,
	fmt.Sprintf("%s will return api details of a single method of certain url", tools.GetSingleAPIDetail),
	getSingleAPIDetail,
)

// AddGetSingleAPIDetailTool can register readOpenAPIDocument to MCP Server
func AddGetSingleAPIDetailTool(mcp *server.MCPServer) {
	GetSingleAPIDetailTool.Register(mcp)
}
