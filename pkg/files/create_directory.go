package files

import (
	"context"
	"errors"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func GetCreateDirectory() (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("create_directory",
		mcp.WithDescription("Create a new directory or ensure a directory exists"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("Path for the directory to create"),
		),
	), createDirectoryHandler
}

func createDirectoryHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	dirPath, ok := request.Params.Arguments["path"].(string)
	if !ok {
		return nil, errors.New("directory path is required")
	}

	err := createDirectory(dirPath)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText("Successfully created directory " + dirPath), nil
}
