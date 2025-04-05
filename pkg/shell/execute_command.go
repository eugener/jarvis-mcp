package shell

import (
	"context"
	"errors"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// GetExecuteCommand returns the tool and handler for executing OS commands.
func GetExecuteCommand() (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("execute_command",
		mcp.WithDescription("Execute OS command"),
		mcp.WithString("command",
			mcp.Required(),
			mcp.Description("Full OS command to execute"),
		),
		mcp.WithString("working directory",
			mcp.Description("Working directory for the command"),
		),
	), executeCommandHandler
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

	result, err := executeCommand(cmd, workDir)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(result), nil
}
