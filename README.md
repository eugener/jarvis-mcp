# JARVIS MCP

> Just A Rather Very Intelligent System - Machine Command Proxy

JARVIS MCP is a lightweight server that provides secure access to local machine commands and file operations via a standardized API interface. Inspired by Tony Stark's AI assistant, JARVIS MCP acts as a bridge between applications and your local system.

## Overview

JARVIS MCP implements the Model-Code-Proxy (MCP) architecture to provide a secure, standardized way for applications to execute commands and perform file operations on a local machine. It serves as an intermediary layer that accepts requests through a well-defined API, executes operations in a controlled environment, and returns formatted results.

## Features

- **Command Execution**: Run shell commands on the local system with proper error handling
- **File Operations**: Read, write, and manage files on the local system
- **Working Directory Support**: Execute commands in specific directories
- **Robust Error Handling**: Detailed error messages and validation
- **Comprehensive Output**: Capture and return both stdout and stderr
- **Simple Integration**: Standard I/O interface for easy integration with various clients

## Installation

### Prerequisites

- Go 1.24.1 or higher
- Git (for cloning the repository)

### Building from Source

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd jarvis-mcp
   ```

2. Build the application using the provided script:
   ```bash
   ./build.sh
   ```

   The executable will be created in the `out` directory.

### Cross-Platform Build Instructions

#### Linux

```bash
# Build for Linux (current architecture)
GOOS=linux GOARCH=amd64 go build -o out/jarvis-mcp-linux-amd64 ./cmd/jarvis
chmod +x ./out/jarvis-mcp-linux-amd64

# For ARM64 (like Raspberry Pi)
GOOS=linux GOARCH=arm64 go build -o out/jarvis-mcp-linux-arm64 ./cmd/jarvis
chmod +x ./out/jarvis-mcp-linux-arm64
```

#### macOS

```bash
# Build for macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o out/jarvis-mcp-macos-intel ./cmd/jarvis
chmod +x ./out/jarvis-mcp-macos-intel

# Build for macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o out/jarvis-mcp-macos-arm64 ./cmd/jarvis
chmod +x ./out/jarvis-mcp-macos-arm64
```

#### Windows

```bash
# Build for Windows
GOOS=windows GOARCH=amd64 go build -o out/jarvis-mcp-windows-amd64.exe ./cmd/jarvis
```

## Usage

### Running the Server

Execute the binary:

```bash
# Linux/macOS
./out/jarvis-mcp

# Windows
.\out\jarvis-mcp-windows-amd64.exe
```

The server communicates via standard input/output, making it easy to integrate with various clients.

## Configuring with Claude Desktop

JARVIS MCP is designed to work seamlessly with Claude Desktop through its tools interface. Here's how to set it up:

### Setup Process

1. **Build JARVIS MCP** for your platform using the instructions above
2. **Open Claude Desktop** application
3. **Access Preferences**:
   - macOS: Click on "Claude" in the menu bar and select "Preferences"
   - Windows: Click on the settings gear icon in the top-right corner
4. **Navigate to the Tools Section** in the left sidebar
5. **Click "Add Tool"** to create a new tool configuration

### Configuring Command Execution Tool

1. Configure the **execute_command** tool:
   - **Name**: Execute Command
   - **Description**: Execute shell commands on your local machine
   - **Path**: Full path to your jarvis-mcp binary (e.g., `/Users/username/jarvis-mcp/out/jarvis-mcp`)
   - **Arguments**: Leave empty (the server uses stdin/stdout)
   - **Working Directory**: Optional; specify a default working directory

2. **Save** the configuration

### Configuring File Operation Tools

You can configure additional tools for specific file operations. For example:

1. Configure the **read_file** tool:
   - **Name**: Read File
   - **Description**: Read the contents of a file on your system
   - **Path**: Same path as your jarvis-mcp binary
   - **Arguments**: Leave empty

2. Configure the **write_file** tool:
   - **Name**: Write File
   - **Description**: Write content to a file on your system
   - **Path**: Same path as your jarvis-mcp binary
   - **Arguments**: Leave empty

3. Configure additional tools following the same pattern for:
   - **list_directory**: List directory contents
   - **create_directory**: Create new directories
   - **move_file**: Move or rename files
   - **search_files**: Search for files
   - **get_file_info**: Get file metadata

### Tool Usage in Conversations

Once configured, you can invoke these tools during conversations with Claude:

1. Type a request like "Please show me the contents of my .bashrc file"
2. Claude will display a tool selection interface
3. Select the appropriate tool (e.g., "Read File")
4. Claude will use JARVIS MCP to execute the operation
5. The results will be displayed in your conversation

### Platform-Specific Path Formats

#### macOS/Linux
```
/Users/username/path/to/jarvis-mcp/out/jarvis-mcp
```

#### Windows
```
C:\Users\username\path\to\jarvis-mcp\out\jarvis-mcp-windows-amd64.exe
```

### Troubleshooting

- **Tool Not Responding**: Ensure the binary path is correct and the file is executable
- **Permission Errors**: Check that Claude Desktop has permission to execute the binary
- **Path Issues**: Use absolute paths to avoid working directory problems
- **Execution Errors**: Ensure the tool has appropriate permissions to access requested files/directories

### API Reference

JARVIS MCP exposes the following tools through its API:

#### Command Tools

##### execute_command

Executes shell commands on the local system.

**Parameters:**
- `command` (string, required): The shell command to execute
- `working directory` (string, optional): Directory where the command should be executed

**Returns:**
- On success: Command output (stdout)
- On failure: Error message and any command output (stderr)

#### File System Tools

##### read_file

Reads the contents of a file.

**Parameters:**
- `path` (string, required): Path to the file to read

**Returns:**
- On success: File contents
- On failure: Error message

##### write_file

Writes content to a file.

**Parameters:**
- `path` (string, required): Path where the file will be written
- `content` (string, required): Content to write to the file

**Returns:**
- On success: Success message
- On failure: Error message

##### create_directory

Creates a new directory.

**Parameters:**
- `path` (string, required): Path for the directory to create

**Returns:**
- On success: Success message
- On failure: Error message

##### list_directory

Lists contents of a directory.

**Parameters:**
- `path` (string, required): Path for the directory to list

**Returns:**
- On success: List of files and directories with [FILE] and [DIR] indicators
- On failure: Error message

##### move_file

Moves or renames files and directories.

**Parameters:**
- `source` (string, required): Source path of the file or directory to move
- `destination` (string, required): Destination path where the file or directory will be moved to

**Returns:**
- On success: Success message
- On failure: Error message

##### search_files

Searches for files matching a pattern.

**Parameters:**
- `path` (string, required): Starting path for the search
- `pattern` (string, required): Search pattern to match file and directory names

**Returns:**
- On success: List of matching files
- On failure: Error message

##### get_file_info

Retrieves detailed metadata about a file or directory.

**Parameters:**
- `path` (string, required): Path for the file or directory to get information about

**Returns:**
- On success: JSON with file metadata (name, size, mode, modification time, etc.)
- On failure: Error message

## Architecture

JARVIS MCP is built on the [MCP Go framework](https://github.com/mark3labs/mcp-go), which implements the Model-Code-Proxy pattern. The architecture consists of:

1. **Request Handling**: Parsing and validating incoming requests
2. **Command Execution**: Running system commands in a controlled manner
3. **File Operations**: Reading from and writing to files on the local system
4. **Response Formatting**: Providing structured, informative responses

### Project Structure

```
jarvis-mcp/
├── build.sh                  # Build script
├── cmd/                      # Application entry points
│   └── jarvis/               # Main JARVIS MCP application
│       └── main.go           # Application entry point
├── pkg/                      # Library packages
│   ├── shell/                # Shell command execution package
│   │   └── execute_command.go # Command execution functionality
│   └── files/                # File operations package
│       ├── files.go          # Core file operation functions
│       ├── read_file.go      # Read file tool implementation
│       ├── write_file.go     # Write file tool implementation
│       ├── create_directory.go # Create directory tool implementation
│       ├── list_directory.go # List directory tool implementation
│       ├── move_file.go      # Move file tool implementation
│       ├── search_files.go   # Search files tool implementation
│       └── get_file_info.go  # Get file info tool implementation
├── go.mod                    # Go module definition
├── go.sum                    # Go module checksums
└── out/                      # Build outputs (gitignored)
    └── jarvis-mcp            # Compiled binary
```

## Security Considerations

JARVIS MCP provides direct access to execute commands and file operations on the local system. Consider the following security practices:

- Run with appropriate permissions (avoid running as root/administrator)
- Use in trusted environments only
- Consider implementing additional authorization mechanisms for production use
- Be cautious about which directories you allow command execution and file operations in
- Implement path validation to prevent unauthorized access to system files

### Platform-Specific Security Notes

#### Linux/macOS
- Run with a dedicated user with limited permissions
- Consider using a chroot environment to restrict file system access
- Use `chmod` to restrict executable permissions: `chmod 700 jarvis-mcp`

#### Windows
- Run as a standard user, not an administrator
- Consider using Windows Security features to restrict access
- Use folder/file permissions to limit access to sensitive directories

## Development

### Adding New Tools

To extend JARVIS MCP with additional functionality, create a new file in the appropriate package following this pattern:

```go
package mypackage

import (
    "context"
    "errors"
    
    "github.com/mark3labs/mcp-go/mcp"
    "github.com/mark3labs/mcp-go/server"
)

func GetMyTool() (tool mcp.Tool, handler server.ToolHandlerFunc) {
    return mcp.NewTool("my_tool",
        mcp.WithDescription("Description of the tool"),
        mcp.WithString("param_name",
            mcp.Required(),
            mcp.Description("Parameter description"),
        ),
    ), myToolHandler
}

func myToolHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    // Parameter validation
    param, ok := request.Params.Arguments["param_name"].(string)
    if !ok {
        return nil, errors.New("parameter is required")
    }
    
    // Tool implementation
    result, err := doSomething(param)
    if err != nil {
        return nil, err
    }
    
    return mcp.NewToolResultText(result), nil
}
```

Then register the tool in `cmd/jarvis/main.go`:

```go
mcpServer.AddTool(mypackage.GetMyTool())
```

## License

[Specify your license here]

## Acknowledgements

- Built with the [MCP Go framework](https://github.com/mark3labs/mcp-go)
- Inspired by Tony Stark's JARVIS from the Marvel Cinematic Universe
