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

	// shell tools
	mcpServer.AddTool(shell.GetExecuteCommand())

	// file system tools
	mcpServer.AddTool(files.GetReadFile())
	mcpServer.AddTool(files.GetWriteFile())
	mcpServer.AddTool(files.GetCreateDirectory())
	mcpServer.AddTool(files.GetListDirectory())
	mcpServer.AddTool(files.GetMoveFile())
	mcpServer.AddTool(files.GetSearchFiles())
	mcpServer.AddTool(files.GetFileInfo())

	// Start the stdio server
	if err := server.ServeStdio(mcpServer); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
