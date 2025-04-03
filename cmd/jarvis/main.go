package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"jarvis_mcp/pkg/shell"
)

func main() {
	// Create MCP server
	mcpServer := server.NewMCPServer(
		"jarvis-mcp",
		"1.0.0",
	)

	cmdTool := mcp.NewTool("execute_command",
		mcp.WithDescription("Execute OS command"),
		mcp.WithString("command",
			mcp.Required(),
			mcp.Description("Full OS command to execute"),
		),
		mcp.WithString("working directory",
			mcp.Description("Working directory for the command"),
		),
	)
	mcpServer.AddTool(cmdTool, executeCommandHandler)

	// Start the stdio server
	if err := server.ServeStdio(mcpServer); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

// executeCommandHandler executes OS commands specified by the user and returns the command output.
// It takes the command string from the request parameters, runs it via the shell, and handles
// optional working directory configuration. The function captures both stdout and stderr output.
func executeCommandHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	cmd, ok := request.Params.Arguments["command"].(string)
	if !ok {
		return nil, errors.New("command must be a string")
	}

	workDir := ""
	if workDirVal, ok := request.Params.Arguments["working directory"].(string); ok {
		workDir = workDirVal
	}

	result, err := shell.ExecuteCommand(cmd, workDir)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(result), nil
}
