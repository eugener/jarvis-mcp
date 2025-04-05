package files

import (
	"context"
	"errors"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func GetDirectoryTree() (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("directory_tree",
		mcp.WithDescription("Get a recursive tree view of files and directories as a JSON structure"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("Path for the directory to generate tree from"),
		),
	), directoryTreeHandler
}

func directoryTreeHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	dirPath, ok := request.Params.Arguments["path"].(string)
	if !ok {
		return nil, errors.New("directory path is required")
	}

	treeJSON, err := directoryTree(dirPath)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(treeJSON), nil
}
