// Package main can be start by sse or stdio
//
// use `-t sse` to start at  sse server
package main

import (
	"flag"
	"fmt"
	"log"
	getsingleapidetail "mcp-api-tester/tools/getSingleAPIDetail"
	listallapifromdocument "mcp-api-tester/tools/listAllAPIFromDocument"
	readopenapidocument "mcp-api-tester/tools/readOpenAPIDocument"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	var transport string
	flag.StringVar(&transport, "t", "stdio", "Transport type, how llm connect to mcp server (stdio or sse)")
	flag.StringVar(&transport, "transport", "stdio", "Transport type, how llm connect to mcp server (stdio or sse)")

	var port string

	flag.StringVar(&port, "p", "8000", "The port that sse server will listen to")
	flag.StringVar(&port, "sse-port", "8000", "The port that sse server will listen to")
	flag.Parse()

	err := run(transport, port)

	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

// run will start server base on transport type
func run(transport string, port string) error {
	srv := newMCPServer()

	switch transport {
	case "stdio":
		return server.ServeStdio(srv)
	case "sse":

		sseURL := fmt.Sprintf("http://localhost:%s", port)
		sseServer := server.NewSSEServer(
			srv,
			server.WithBaseURL(sseURL),
		)
		log.Printf("Server sse listening on %s", sseURL)

		ssePort := fmt.Sprintf(":%s", port)
		return sseServer.Start(ssePort)
	default:
		return fmt.Errorf("Only 'sse' and 'stdio can be used with -t flag, %s is not valid", transport)
	}
}

// newMCPServer will return MCPServer that register all tools
func newMCPServer() *server.MCPServer {
	srv := server.NewMCPServer(
		"mcp-api-tester",
		"0.0.1",
	)

	readopenapidocument.AddReadOpenAPIDocumentTool(srv)
	listallapifromdocument.AddListAllAPIFromDocumentTool(srv)
	getsingleapidetail.AddGetSingleAPIDetailTool(srv)

	return srv
}
