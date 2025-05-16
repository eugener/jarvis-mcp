package files

import (
	"context"
	"errors"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func GetCreateDirectory() (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("create_directory",
		mcp.WithDescription("Create or verify the existence of a directory at the specified path"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("The filesystem path where the directory should be created or verified"),
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
