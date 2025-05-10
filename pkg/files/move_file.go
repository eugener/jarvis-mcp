package files

import (
	"context"
	"errors"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func GetMoveFile() (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("move_file",
		mcp.WithDescription("Move or rename files and directories to a specified location"),
		mcp.WithString("source",
			mcp.Required(),
			mcp.Description("The full path of the file or directory to be moved or renamed"),
		),
		mcp.WithString("destination",
			mcp.Required(),
			mcp.Description("The target path where the file or directory should be moved or renamed to"),
		),
	), moveFileHandler
}

func moveFileHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sourcePath, ok := request.Params.Arguments["source"].(string)
	if !ok {
		return nil, errors.New("source path is required")
	}

	destPath, ok := request.Params.Arguments["destination"].(string)
	if !ok {
		return nil, errors.New("destination path is required")
	}

	err := moveFile(sourcePath, destPath)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText("Successfully moved file from " + sourcePath + " to " + destPath), nil
}
