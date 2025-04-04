package main

import (
	"fmt"
	"jarvis_mcp/pkg/shell"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Create MCP server
	mcpServer := server.NewMCPServer(
		"jarvis-mcp",
		"1.0.0",
	)

	shell.RegisterTools(mcpServer)

	// Start the stdio server
	if err := server.ServeStdio(mcpServer); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
