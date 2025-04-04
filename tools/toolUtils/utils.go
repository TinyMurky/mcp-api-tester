// Package toolUtils provides utility functions for tool handling.
//
// Portions of this file are adapted from:
// https://github.com/grafana/mcp-grafana/blob/main/tools.go
//
// Source project: Grafana mcp-grafana (Apache License 2.0)
// License: https://github.com/grafana/mcp-grafana/blob/main/LICENSE
//
// Modifications were made to fit this project's requirements.
//
// To study reflect package check: https://darjun.github.io/2021/05/27/godailylib/reflect/
package toolUtils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/invopop/jsonschema"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Tool Struct can be register to mcp server
type Tool struct {
	Tool    mcp.Tool
	Handler server.ToolHandlerFunc
}

// ToolHandlerFunc can be parse to generate Tool then be register
type ToolHandlerFunc[Args any, Response any] func(ctx context.Context, argument Args) (Response, error)

// Register Tool to MCP Server
func (tool *Tool) Register(srv *server.MCPServer) {
	srv.AddTool(tool.Tool, tool.Handler)
}

// MustTool creates a new Tool from the given name, description, and toolHandler.
// It panics if the tool cannot be created.
func MustTool[T any, R any](name, description string, toolHandler ToolHandlerFunc[T, R]) Tool {
	tool, handler, err := ConvertTool(name, description, toolHandler)
	if err != nil {
		panic(err)
	}
	return Tool{Tool: tool, Handler: handler}
}

// ConvertTool converts a toolHandler function to a Tool and ToolHandlerFunc.
//
// The toolHandler function must have two arguments: a context.Context and a struct
// to be used as the parameters for the tool. The second argument must not be a pointer,
// should be marshalable to JSON, and the fields should have a `jsonschema` tag with the
// description of the parameter.
func ConvertTool[Args any, Response any](
	name string,
	description string,
	toolHandler ToolHandlerFunc[Args, Response],
) (mcp.Tool, server.ToolHandlerFunc, error) {
	toolFuncSchema := createJSONSchemaFromToolHandlerFunc(toolHandler)

	// this contain property like mcp.Description...
	properties := make(map[string]interface{}, toolFuncSchema.Properties.Len())

	for pair := toolFuncSchema.Properties.Oldest(); pair != nil; pair = pair.Next() {
		properties[pair.Key] = pair.Value // value is *jsonschema.Schema
	}

	toolInputSchema := mcp.ToolInputSchema{
		Type:       toolFuncSchema.Type,
		Properties: properties, // I don't know why this work
		Required:   toolFuncSchema.Required,
	}

	return mcp.Tool{
		Name:        name,
		Description: description,
		InputSchema: toolInputSchema,
	}, nil, nil
}

func createMCPHandlerFunc[Args any, Response any](toolFunc ToolHandlerFunc[Args, Response]) server.ToolHandlerFunc {
	handlerValue := reflect.ValueOf(toolFunc)
	handlerType := handlerValue.Type()

	// Get the second parameter's type of the toolFunc, which represents the actual request argument type
	secondArgType := handlerType.In(1)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// The request.Params.Arguments is of type `map[string]interface{}` (decoded JSON)
		// We need to convert it into the actual argument type `Args` for toolFunc

		// Step 1: Marshal the map[string]interface{} back into JSON
		s, err := json.Marshal(request.Params.Arguments)

		if err != nil {
			return nil, fmt.Errorf("Marshal request.Params.Arguments failed, origin args: %v, error is: %w", request.Params.Arguments, err)
		}

		// Step 2: Create a new zero value of the expected argument type (Args)
		// - reflect.New(secondArgType) returns a reflect.Value of type *T (pointer to the type)
		// - .Interface() wraps that pointer into an interface{} so we can use it in json.Unmarshal
		//
		// For example:
		//   If secondArgType is `main.Cat`, reflect.New returns `*main.Cat`
		//   So the resulting value is: (type: *main.Cat, value: <*main.Cat Value>)
		//   This is still an interface{}, but it holds the concrete pointer to Cat under the hood
		unmarshaledArgs := reflect.New(secondArgType).Interface()

		// Step 3. Write marshal json from LLM Arguments input to the Args struct that toolFunc second Argument is
		if err := json.Unmarshal([]byte(s), unmarshaledArgs); err != nil {
			return nil, fmt.Errorf("request.Params.Arguments unmarshal to args failed, error: %w", err)
		}

		of := reflect.ValueOf(unmarshaledArgs)

		// Step 3.1 unmarshaledArgs should be a pointer point to args type of second input of toolFunc
		if of.Kind() != reflect.Ptr || !of.Elem().CanInterface() {
			return nil, errors.New("arguments must be a struct")
		}

		args := []reflect.Value{reflect.ValueOf(ctx), of.Elem()}

		output := handlerValue.Call(args)

		// Step 4, Check if output is validate
		if len(output) != 2 {
			return nil, errors.New("tool handler must return 2 values")
		}

		if !output[0].CanInterface() { // Can output[0] be convert to interface{}, aka no unexported lower case field
			return nil, errors.New("tool handler first return value must be interfaceable")
		}

		// Step 5, handle Error first, the second response of toolFunc Should be error
		var handlerErr error
		var ok bool

		if output[1].Kind() == reflect.Interface && !output[1].IsNil() { // make sure output[1] has value and is interface(error is interface)
			handlerErr, ok = output[1].Interface().(error)
			if !ok {
				return nil, errors.New("tool handler second return value must be error")
			}
		}

		// Step 5.1 after make sure  output[1] is error interface, return if not nil
		if handlerErr != nil {
			// We use NewToolResultError to tell LLM what error has happened
			// return nil, handlerErr
			return NewToolResultError(handlerErr.Error()), nil
		}

		// Step 6 Check if the first return value is nil (only for pointer, interface, map, etc.)
		isNilable := output[0].Kind() == reflect.Ptr ||
			output[0].Kind() == reflect.Interface ||
			output[0].Kind() == reflect.Map ||
			output[0].Kind() == reflect.Slice ||
			output[0].Kind() == reflect.Chan ||
			output[0].Kind() == reflect.Func

		if isNilable && output[0].IsNil() {
			return nil, nil
		}

		// Step7 use interface to change the type of output to any
		returnVal := output[0].Interface()
		returnType := output[0].Type()

		// Step 8 return base on case

		// Case 1: Already a *mcp.CallToolResult (case any to *mcp.CallToolResult to check if is ok)
		if callResult, ok := returnVal.(*mcp.CallToolResult); ok {
			return callResult, nil
		}

		// Case 2: An mcp.CallToolResult (not a pointer), need to return pointer
		if returnType.ConvertibleTo(reflect.TypeOf(mcp.CallToolResult{})) { // can be convert to mcp.CallToolResult{}
			callResult := returnVal.(mcp.CallToolResult)
			return &callResult, nil
		}

		// Case 3: String or *string
		if str, ok := returnVal.(string); ok {
			if str == "" {
				return nil, nil
			}
			return mcp.NewToolResultText(str), nil
		}

		if strPtr, ok := returnVal.(*string); ok {
			if strPtr == nil || *strPtr == "" {
				return nil, nil
			}
			return mcp.NewToolResultText(*strPtr), nil
		}

		// Case 4: Any other type - marshal to JSON
		jsonBytes, err := json.Marshal(returnVal)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal return value: %s", err)
		}

		return mcp.NewToolResultText(string(jsonBytes)), nil
	}
	return handler
}

var (
	// This special reflector can extract from struct description that start with `jsonschema`
	jsonSchemaReflector = jsonschema.Reflector{
		// BaseSchemaID:               "",
		// Anonymous:                  true,
		// AssignAnchor:               false,
		// AllowAdditionalProperties:  true,
		RequiredFromJSONSchemaTags: true, // generate a schema that requires any key tagged with `jsonschema:required`
		DoNotReference:             true, // This can let return of ReflectFromType (which is *jsonschema.Schema) actually get property
		// ExpandedStruct:             true,
		// FieldNameTag:               "",
		// IgnoredTypes:               nil,
		// Lookup:                     nil,
		// Mapper:                     nil,
		// Namer:                      nil,
		// KeyNamer:                   nil,
		// AdditionalFields:           nil,
		// CommentMap:                 nil,
	}
)

// createJSONSchemaFromToolHandlerFunc will get second args type
// of ToolHandlerFunc and get struct schema
func createJSONSchemaFromToolHandlerFunc(toolFunc any) *jsonschema.Schema {

	// this will return like func(type, type)
	funcType := reflect.ValueOf(toolFunc).Type()

	// if arg is struct type, we can get struct and parse by json
	secondArgType := funcType.In(1)
	schema := jsonSchemaReflector.ReflectFromType(secondArgType)

	return schema
}
