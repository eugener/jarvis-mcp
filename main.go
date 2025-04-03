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

	command := exec.Command("sh", "-c", cmd)
	command.Env = os.Environ() // Explicitly copy the current

	if workDir, ok := request.Params.Arguments["working directory"].(string); ok && workDir != "" {
		dirInfo, err := os.Stat(workDir)
		if err != nil {
			if os.IsNotExist(err) {
				return nil, fmt.Errorf("Path '%s' does not exist\n", workDir)
			}
			return nil, fmt.Errorf("Error checking path: %v\n", err)
		}

		if !dirInfo.IsDir() {
			return nil, fmt.Errorf("'%s' is not a directory", workDir)
		}

		command.Dir = workDir
	}

	// Execute command and capture output
	output, err := command.CombinedOutput()
	outputStr := string(output)

	// Format the response
	var resultText string
	if err != nil {
		// Return both the error and any output
		resultText = fmt.Sprintf("Command failed: %s\n\nOutput:\n%s\n\nError: %v",
			cmd, outputStr, err)
	} else {
		resultText = fmt.Sprintf("Command executed successfully: %s\n\nOutput:\n%s",
			cmd, outputStr)
	}

	return mcp.NewToolResultText(resultText), nil
}
