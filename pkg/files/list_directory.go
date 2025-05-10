package files

import (
	"context"
	"errors"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func GetListDirectory() (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("list_directory",
		mcp.WithDescription("Retrieve a comprehensive listing of all files and subdirectories within a specified directory path"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("The absolute or relative path of the directory to be listed"),
		),
	), listDirectoryHandler
}
func listDirectoryHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	dirPath, ok := request.Params.Arguments["path"].(string)
	if !ok {
		return nil, errors.New("directory path is required")
	}

	entries, err := listDirectory(dirPath)
	if err != nil {
		return nil, err
	}

	// Format the output for better readability
	output := strings.Join(entries, "\n")
	return mcp.NewToolResultText(output), nil
}
