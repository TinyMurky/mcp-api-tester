package tools

// GetSingleAPIDetail is the tool that return OpenAPI operation detail
const GetSingleAPIDetail = "GetSingleAPIDetail"

// ListAllAPIFromDocument is the tool name of listAllAPIFromDocument
const ListAllAPIFromDocument = "ListAllAPIFromDocument"

// ReadOpenAPIDocument is the tool name odf readOpenAPIDocument
const ReadOpenAPIDocument = "ReadOpenAPIDocument"

// ToolNames has all tools' name in this project
var ToolNames = map[string]string{
	GetSingleAPIDetail:     GetSingleAPIDetail,
	ListAllAPIFromDocument: ListAllAPIFromDocument,
	ReadOpenAPIDocument:    ReadOpenAPIDocument,
}
