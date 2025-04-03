package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Create MCP server
	mcpServer := server.NewMCPServer(
		"mcp-test",
		"1.0.0",
	)

	// Add tool
	// tool := mcp.NewTool("hello_world",
	// 	mcp.WithDescription("Say hello to someone"),
	// 	mcp.WithString("name",
	// 		mcp.Required(),
	// 		mcp.Description("Name of the person to greet"),
	// 	),
	// )
	// mcpServer.AddTool(tool, helloHandler)

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

// func helloHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
// 	name, ok := request.Params.Arguments["name"].(string)
// 	if !ok {
// 		return nil, errors.New("name must be a string")
// 	}

// 	return mcp.NewToolResultText(fmt.Sprintf("Hello, %s!", name)), nil
// }

func executeCommandHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	cmd, ok := request.Params.Arguments["command"].(string)
	if !ok {
		return nil, errors.New("command must be a string")
	}

	command := exec.Command("sh", "-c", cmd)
	command.Env = os.Environ() // Explicitly copy the current

	if workDir, ok := request.Params.Arguments["working directory"].(string); ok {
		dirInfo, err := os.Stat(workDir)
		if err != nil {
			if os.IsNotExist(err) {
				return nil, fmt.Errorf("Path '%s' does not exist\n", dirInfo)
			}
			return nil, fmt.Errorf("Error checking path: %v\n", err)
		}
		command.Dir = workDir
	}

	// Execute command
	output, err := command.CombinedOutput()
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(fmt.Sprintf("Command executed: %s\nOutput:%s", cmd, output)), nil
}
