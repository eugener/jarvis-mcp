package files

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func GetFileInfo() (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("get_file_info",
		mcp.WithDescription("Retrieve comprehensive metadata and attributes for a specified file or directory"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("The absolute or relative path of the file or directory to retrieve metadata for"),
		),
	), getFileInfoHandler
}

func getFileInfoHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	filePath, ok := request.Params.Arguments["path"].(string)
	if !ok {
		return nil, errors.New("file path is required")
	}

	info, err := getFileInfo(filePath)
	if err != nil {
		return nil, err
	}

	// Format the output as JSON for better readability
	jsonData, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error formatting file info: %v", err)
	}

	return mcp.NewToolResultText(string(jsonData)), nil
}
