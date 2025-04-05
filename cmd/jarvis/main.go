package main

import (
	"fmt"
	"jarvis_mcp/pkg/files"
	"jarvis_mcp/pkg/shell"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Create MCP server
	mcpServer := server.NewMCPServer(
		"jarvis-mcp",
		"1.0.0",
	)

	mcpServer.AddTool(shell.GetExecuteCommand())

	mcpServer.AddTool(files.GetReadFile())
	mcpServer.AddTool(files.GetWriteFile())

	// Start the stdio server
	if err := server.ServeStdio(mcpServer); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
