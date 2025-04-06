package shell

import (
	"strings"
	"testing"
)

// TestBasicCommands tests basic command execution functionality
func TestBasicCommands(t *testing.T) {
	helper := NewTestHelper(t, "jarvis-basic-test")
	defer helper.Cleanup()
	
	// Create a test file
	helper.CreateTestFile("test.txt", "Hello, JARVIS MCP!")
	
	// Get platform-specific commands
	pc := GetPlatformCommands()
	
	// Test cases
	tests := []struct {
		name        string
		cmd         string
		workDir     string
		wantErr     bool
		wantErrMsg  string
		wantContain string
	}{
		{
			name:        "echo command",
			cmd:         "echo 'Hello, World!'",
			workDir:     "",
			wantErr:     false,
			wantContain: "Hello, World!",
		},
		{
			name:        "command with working directory",
			cmd:         pc.ReadFile + " test.txt",
			workDir:     helper.TempDir,
			wantErr:     false,
			wantContain: "Hello, JARVIS MCP!",
		},
		{
			name:        "invalid command",
			cmd:         "command_that_does_not_exist",
			workDir:     "",
			wantErr:     true,
			wantErrMsg:  "exit status",
			wantContain: "Command failed",
		},
		{
			name:        "command with environment variables",
			cmd:         "echo " + pc.EnvVarSyntax,
			workDir:     "",
			wantErr:     false,
			wantContain: "", // Just check it doesn't error
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := executeCommand(tt.cmd, tt.workDir)
			
			// Check error cases
			if tt.wantErr {
				if err == nil {
					t.Errorf("executeCommand() error = nil, wantErr = true")
					return
				}
				
				if tt.wantErrMsg != "" && !strings.Contains(err.Error(), tt.wantErrMsg) {
					t.Errorf("executeCommand() error = %v, want error containing %v", err, tt.wantErrMsg)
				}
			} else if err != nil {
				t.Errorf("executeCommand() unexpected error = %v", err)
				return
			}
			
			// Check the output content for all cases
			if tt.wantContain != "" && !strings.Contains(result, tt.wantContain) {
				t.Errorf("executeCommand() result = %v, want result containing %v", result, tt.wantContain)
			}
		})
	}
}

// TestWithInvalidWorkingDirectory tests handling of invalid working directory
func TestWithInvalidWorkingDirectory(t *testing.T) {
	// Test empty working directory
	_, err := executeCommand("echo 'test'", "")
	if err != nil {
		t.Errorf("executeCommand() with empty working directory should not error, got = %v", err)
	}
	
	// Get a test helper for access to temp directory
	helper := NewTestHelper(t, "jarvis-workdir-test")
	defer helper.Cleanup()
	
	// Create a test file to use as an invalid working directory
	testFile := helper.CreateTestFile("file.txt", "Not a directory")
	
	// Test cases for invalid working directories
	tests := []struct {
		name    string
		workDir string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "non-existent working directory",
			workDir: "/path/that/does/not/exist",
			wantErr: true,
			errMsg:  "does not exist",
		},
		{
			name:    "file as working directory",
			workDir: testFile,
			wantErr: true,
			errMsg:  "is not a directory",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// For working directory errors, we only check the error message
			_, err := executeCommand("echo 'test'", tt.workDir)
			
			if !tt.wantErr && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			} else if tt.wantErr && err == nil {
				t.Errorf("Expected error for invalid working directory but got none")
			} else if tt.wantErr && err != nil && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("Error message doesn't contain %q, got: %v", tt.errMsg, err)
			}
		})
	}
}

// TestWithLargeOutput tests handling of commands that produce large output
func TestWithLargeOutput(t *testing.T) {
	// Use appropriate command based on platform
	var cmd string
	if IsWindows() {
		cmd = "dir /s C:\\Windows"
	} else {
		cmd = "find /usr -type f -name '*.go' 2>/dev/null || echo 'No Go files found'"
	}
	
	result, err := executeCommand(cmd, "")
	if err != nil {
		t.Errorf("executeCommand() with large output should not error, got = %v", err)
	}
	
	if len(result) < 100 {
		t.Errorf("executeCommand() with large output returned suspiciously small result: %d bytes", len(result))
	}
}

// TestWithNonAsciiOutput tests handling of non-ASCII characters in command output
func TestWithNonAsciiOutput(t *testing.T) {
	nonAscii := "你好，世界！"
	var cmd string
	if IsWindows() {
		cmd = "echo " + nonAscii
	} else {
		cmd = "echo '" + nonAscii + "'"
	}
	
	result, err := executeCommand(cmd, "")
	if err != nil {
		t.Errorf("executeCommand() with non-ASCII output should not error, got = %v", err)
	}
	
	if !strings.Contains(result, nonAscii) {
		t.Errorf("executeCommand() failed to correctly handle non-ASCII output")
	}
}
