package files

import (
	"context"
	"errors"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func GetSearchFiles() (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("search_files",
		mcp.WithDescription("Recursively search for files and directories matching a pattern"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("Starting path for the search"),
		),
		mcp.WithString("pattern",
			mcp.Required(),
			mcp.Description("Search pattern to match file and directory names"),
		),
	), searchFilesHandler
}

func searchFilesHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	dirPath, ok := request.Params.Arguments["path"].(string)
	if !ok {
		return nil, errors.New("search path is required")
	}

	pattern, ok := request.Params.Arguments["pattern"].(string)
	if !ok {
		return nil, errors.New("search pattern is required")
	}

	files, err := searchFiles(dirPath, pattern)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return mcp.NewToolResultText("No files matching the pattern were found"), nil
	}

	// Format the output for better readability
	output := "Files matching pattern '" + pattern + "':\n" + strings.Join(files, "\n")
	return mcp.NewToolResultText(output), nil
}
