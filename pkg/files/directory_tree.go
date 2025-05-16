package files

import (
	"context"
	"errors"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func GetDirectoryTree() (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("directory_tree",
		mcp.WithDescription("Retrieve a detailed, recursive tree structure of files and directories in JSON format"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("The file system path of the directory for which to generate the tree structure"),
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
