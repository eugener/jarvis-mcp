package files

import (
	"context"
	"errors"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func GetReadFile() (tool mcp.Tool, handler server.ToolHandlerFunc) {

	return mcp.NewTool("read_file",
		mcp.WithDescription("Read file, given the path"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("Path for the file name to read"),
		),
	), readFileHandler
}

func readFileHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {

	fileName, ok := request.Params.Arguments["path"].(string)
	if !ok {
		return nil, errors.New("file path is required")
	}

	content, err := readFile(fileName)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText("File read successfully. Content: " + content), nil
}
