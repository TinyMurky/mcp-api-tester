// Package openapi will read openapi file from giving location
// use https://github.com/pb33f/libopenapi to read openapi
// more doc https://quobix.com/articles/parsing-openapi-using-go/
package openapi

// This instance.go is where to put variable that will be constantly reuse

// OpenAPIPointer Point to the OpenAPI instance that initialized by ReadFromFile
var OpenAPIPointer *OpenAPI
