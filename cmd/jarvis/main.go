package main

import (
	"fmt"
	"jarvis_mcp/pkg/files"
	"jarvis_mcp/pkg/shell"

	"github.com/mark3labs/mcp-go/server"
)

type ToolRegistrator = func(mcpServer *server.MCPServer)

var registrators = []ToolRegistrator{
	shell.RegisterTools,
	files.RegisterTools,
}

func main() {
	// Create MCP server
	mcpServer := server.NewMCPServer(
		"jarvis-mcp",
		"1.0.0",
	)

	for _, registerTool := range registrators {
		registerTool(mcpServer)
	}

	// Start the stdio server
	if err := server.ServeStdio(mcpServer); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
