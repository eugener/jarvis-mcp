package files

import (
	"context"
	"errors"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func GetWriteFile() (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("write_file",
		mcp.WithDescription("Write file, given the path"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("Path for the file name to write"),
		),
		mcp.WithString("content",
			mcp.Required(),
			mcp.Description("Content to write to the file"),
		),
	), writeFileHandler
}

func writeFileHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {

	fileName, ok := request.Params.Arguments["path"].(string)
	if !ok {
		return nil, errors.New("file path is required")
	}

	content, ok := request.Params.Arguments["content"].(string)
	if !ok {
		return nil, errors.New("file content is required")
	}

	err := writeFile(fileName, content)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText("File written successfully"), nil
}
