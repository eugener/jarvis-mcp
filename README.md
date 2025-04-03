# JARVIS MCP

> Just A Rather Very Intelligent System - Machine Command Proxy

JARVIS MCP is a lightweight server that provides secure access to local machine commands via a standardized API interface. Inspired by Tony Stark's AI assistant, JARVIS MCP acts as a bridge between applications and your local system.

## Overview

JARVIS MCP implements the Model-Code-Proxy (MCP) architecture to provide a secure, standardized way for applications to execute commands on a local machine. It serves as an intermediary layer that accepts requests through a well-defined API, executes commands in a controlled environment, and returns formatted results.

## Features

- **Command Execution**: Run shell commands on the local system with proper error handling
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

2. Build the application:
   ```bash
   ./build.sh
   ```

The executable will be created in the `out` directory.

## Usage

### Running the Server

Execute the binary:

```bash
./out/jarvis-mcp
```

The server communicates via standard input/output, making it easy to integrate with various clients.

### Configuration with Claude Desktop

To use JARVIS MCP with Claude Desktop:

1. Open Claude Desktop preferences
2. Navigate to the "Tools" section
3. Add a new tool with the following configuration:
   - **Name**: Execute Command
   - **Description**: Execute shell commands on your local machine
   - **Path**: `/path/to/jarvis-mcp/out/jarvis-mcp`
   - **Arguments**: Leave empty (the server uses stdin/stdout)
   - **Working Directory**: `/path/to/preferred/default/directory`

4. Save the configuration

Once configured, you can invoke the "Execute Command" tool directly from conversations with Claude, allowing you to run system commands through natural language requests.

### API Reference

JARVIS MCP exposes the following tools through its API:

#### execute_command

Executes shell commands on the local system.

**Parameters:**
- `command` (string, required): The shell command to execute
- `working directory` (string, optional): Directory where the command should be executed

**Returns:**
- On success: Command output (stdout)
- On failure: Error message and any command output (stderr)

## Architecture

JARVIS MCP is built on the [MCP Go framework](https://github.com/mark3labs/mcp-go), which implements the Model-Code-Proxy pattern. The architecture consists of:

1. **Request Handling**: Parsing and validating incoming requests
2. **Command Execution**: Running system commands in a controlled manner
3. **Response Formatting**: Providing structured, informative responses

### Project Structure

```
jarvis-mcp/
├── build.sh                  # Build script
├── cmd/                      # Application entry points
│   └── jarvis/               # Main JARVIS MCP application
│       └── main.go           # Application entry point
├── pkg/                      # Library packages
│   └── shell/                # Shell command execution package
│       └── shell.go          # Command execution logic
├── go.mod                    # Go module definition
├── go.sum                    # Go module checksums
└── out/                      # Build outputs (gitignored)
    └── jarvis-mcp            # Compiled binary
```

## Security Considerations

JARVIS MCP provides direct access to execute commands on the local system. Consider the following security practices:

- Run with appropriate permissions (avoid running as root/administrator)
- Use in trusted environments only
- Consider implementing additional authorization mechanisms for production use
- Be cautious about which directories you allow command execution in

## Development

### Adding New Tools

To extend JARVIS MCP with additional functionality, you can add new tools following this pattern:

```go
// Define a new tool
newTool := mcp.NewTool("tool_name",
    mcp.WithDescription("Description of the tool"),
    mcp.WithString("param_name",
        mcp.Required(),
        mcp.Description("Parameter description"),
    ),
)

// Register the tool with a handler
mcpServer.AddTool(newTool, toolHandler)
```

## License

[Specify your license here]

## Acknowledgements

- Built with the [MCP Go framework](https://github.com/mark3labs/mcp-go)
- Inspired by Tony Stark's JARVIS from the Marvel Cinematic Universe