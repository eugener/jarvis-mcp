package shell

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestExecuteCommand(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "jarvis-mcp-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test file in the temp directory
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "Hello, JARVIS MCP!"
	err = os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

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
			cmd:         "cat test.txt",
			workDir:     tempDir,
			wantErr:     false,
			wantContain: testContent,
		},
		{
			name:        "invalid command",
			cmd:         "command_that_does_not_exist",
			workDir:     "",
			wantErr:     true,
			wantErrMsg:  "exit status",  // Just check for "exit status" which should be present in error
			wantContain: "Command failed",
		},
		{
			name:        "non-existent working directory",
			cmd:         "echo 'test'",
			workDir:     "/path/that/does/not/exist",
			wantErr:     true,
			wantErrMsg:  "does not exist",
		},
		{
			name:        "working directory is a file",
			cmd:         "echo 'test'",
			workDir:     testFile,
			wantErr:     true,
			wantErrMsg:  "is not a directory",
		},
		{
			name:        "command with environment variables",
			cmd:         "echo $HOME",
			workDir:     "",
			wantErr:     false,
			wantContain: os.Getenv("HOME"),
		},
		{
			name:        "multiple commands with pipe",
			cmd:         "echo 'test' | grep 'test'",
			workDir:     "",
			wantErr:     false,
			wantContain: "test",
		},
	}

	// Skip certain tests on Windows
	if runtime.GOOS == "windows" {
		// Filter out tests that won't work on Windows
		var windowsCompatibleTests []struct {
			name        string
			cmd         string
			workDir     string
			wantErr     bool
			wantErrMsg  string
			wantContain string
		}

		for _, test := range tests {
			// Skip tests with Unix-specific commands
			if strings.Contains(test.cmd, "cat ") {
				continue
			}
			
			// Adapt commands for Windows if needed
			if test.name == "command with environment variables" {
				test.cmd = "echo %HOME%"
			}
			
			windowsCompatibleTests = append(windowsCompatibleTests, test)
		}
		
		tests = windowsCompatibleTests
	}

	// Run tests
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

func TestExecuteCommand_WithInvalidWorkingDirectory(t *testing.T) {
	// Additional test case where we test a situation that caused a panic in the past
	_, err := executeCommand("echo 'test'", "")
	if err != nil {
		t.Errorf("executeCommand() with empty working directory should not error, got = %v", err)
	}
}

func TestExecuteCommand_WithLargeOutput(t *testing.T) {
	// Test handling of large command output
	cmd := ""
	if runtime.GOOS == "windows" {
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

func TestExecuteCommand_WithNonAsciiOutput(t *testing.T) {
	// Test handling of non-ASCII characters in output
	nonAscii := "你好，世界！"
	cmd := ""
	if runtime.GOOS == "windows" {
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

func TestExecuteCommandErrorHandling(t *testing.T) {
	// Test detailed error handling with expected error returns
	tests := []struct {
		name        string
		cmd         string
		workDir     string
		wantContain string
	}{
		{
			name:        "command not found",
			cmd:         "nonexistentcommand",
			workDir:     "",
			wantContain: "Command failed",
		},
		{
			name:        "permission denied",
			cmd:         "touch /permission_denied_test.txt",
			workDir:     "",
			wantContain: "Command failed",
		},
	}

	// Skip permission test on Windows as it works differently
	if runtime.GOOS == "windows" {
		var filteredTests []struct {
			name        string
			cmd         string
			workDir     string
			wantContain string
		}
		
		for _, test := range tests {
			if test.name != "permission denied" {
				filteredTests = append(filteredTests, test)
			}
		}
		
		tests = filteredTests
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := executeCommand(tt.cmd, tt.workDir)
			
			// We just want to verify the command is detected as failed
			// and proper output is returned, without asserting specific error messages
			if err == nil {
				t.Errorf("executeCommand() should return error for %s", tt.name)
				return
			}
			
			// Just check that the result contains the expected text
			if !strings.Contains(result, tt.wantContain) {
				t.Errorf("executeCommand() result doesn't contain expected text '%s', got: %v", 
					tt.wantContain, result)
			}
		})
	}
}
