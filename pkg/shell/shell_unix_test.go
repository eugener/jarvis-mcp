//go:build !windows
// +build !windows

package shell

import (
	"os"
	"strings"
	"testing"
)

// TestUnixSpecificCommands tests Unix-specific command execution
// This test will only run on Unix systems (Linux, macOS) due to the build tag
func TestUnixSpecificCommands(t *testing.T) {
	helper := NewTestHelper(t, "jarvis-mcp-unix-test")
	defer helper.Cleanup()
	
	// Create a test file
	helper.CreateTestFile("test.txt", "Hello, JARVIS MCP on Unix!")
	
	// Unix-specific test cases
	tests := []struct {
		name        string
		cmd         string
		workDir     string
		wantErr     bool
		wantContain string
	}{
		{
			name:        "ls command",
			cmd:         "ls",
			workDir:     helper.TempDir,
			wantErr:     false,
			wantContain: "test.txt",
		},
		{
			name:        "cat command",
			cmd:         "cat test.txt",
			workDir:     helper.TempDir,
			wantErr:     false,
			wantContain: "Hello, JARVIS MCP on Unix!",
		},
		{
			name:        "echo with Unix environment variable",
			cmd:         "echo $HOME",
			workDir:     "",
			wantErr:     false,
			wantContain: "", // Just check it doesn't error
		},
		{
			name:        "Unix command with pipe",
			cmd:         "ls | grep test",
			workDir:     helper.TempDir,
			wantErr:     false,
			wantContain: "test.txt",
		},
		{
			name:        "ps command",
			cmd:         "ps",
			workDir:     "",
			wantErr:     false,
			wantContain: "", // Just check it doesn't error
		},
		{
			name:        "whoami command",
			cmd:         "whoami",
			workDir:     "",
			wantErr:     false,
			wantContain: "", // Just check it doesn't error
		},
		{
			name:        "command with output redirection",
			cmd:         "echo 'Hello' > output.txt && cat output.txt",
			workDir:     helper.TempDir,
			wantErr:     false,
			wantContain: "Hello",
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
			} else if err != nil {
				t.Errorf("executeCommand() unexpected error = %v", err)
				return
			}
			
			// Check the output content if expected
			if tt.wantContain != "" && !strings.Contains(result, tt.wantContain) {
				t.Errorf("executeCommand() result = %v, want result containing %v", result, tt.wantContain)
			}
		})
	}
}

// TestUnixPathHandling tests Unix-specific path handling
func TestUnixPathHandling(t *testing.T) {
	helper := NewTestHelper(t, "jarvis-unix-paths")
	defer helper.Cleanup()
	
	// Create nested directories
	nestedPath := helper.TempDir + "/level1/level2/level3"
	err := os.MkdirAll(nestedPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create nested directories: %v", err)
	}
	
	// Create test file in nested directory
	nestedFile := nestedPath + "/nested.txt"
	err = os.WriteFile(nestedFile, []byte("Nested file content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create nested file: %v", err)
	}
	
	// Test cases for Unix path handling
	tests := []struct {
		name    string
		cmd     string
		workDir string
		wantErr bool
	}{
		{
			name:    "absolute path",
			cmd:     "ls " + helper.TempDir,
			workDir: "",
			wantErr: false,
		},
		{
			name:    "path with ~",
			cmd:     "echo test",
			workDir: "~", // This works on Unix but not directly testable
			wantErr: false,
		},
		{
			name:    "deep nested path",
			cmd:     "cat nested.txt",
			workDir: nestedPath,
			wantErr: false,
		},
		{
			name:    "relative path navigation",
			cmd:     "ls ../../",
			workDir: nestedPath,
			wantErr: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// For the ~ test, skip if workDir is ~
			if tt.workDir == "~" {
				t.Skip("Skipping ~ test as it can't be directly tested")
				return
			}
			
			_, err := executeCommand(tt.cmd, tt.workDir)
			if tt.wantErr && err == nil {
				t.Errorf("Expected error for %s, got none", tt.name)
			} else if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error for %s: %v", tt.name, err)
			}
		})
	}
}

// TestUnixSpecificOperators tests Unix-specific command operators
func TestUnixSpecificOperators(t *testing.T) {
	helper := NewTestHelper(t, "jarvis-unix-operators")
	defer helper.Cleanup()
	
	// Test Unix-specific operators and control characters
	tests := []struct {
		name    string
		cmd     string
		wantErr bool
	}{
		{
			name:    "background process with &",
			cmd:     "sleep 1 & echo 'Background process started'",
			wantErr: false,
		},
		{
			name:    "command substitution with backticks",
			cmd:     "echo `echo nested command`",
			wantErr: false,
		},
		{
			name:    "command substitution with $()",
			cmd:     "echo $(echo nested command)",
			wantErr: false,
		},
		{
			name:    "redirecting stderr",
			cmd:     "ls /nonexistent 2>/dev/null || echo 'Error redirected'",
			wantErr: false,
		},
		{
			name:    "heredoc syntax",
			cmd:     "cat << EOF\nThis is heredoc content\nEOF",
			wantErr: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := executeCommand(tt.cmd, helper.TempDir)
			
			if tt.wantErr && err == nil {
				t.Errorf("Expected error for command %s but got none", tt.cmd)
			} else if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error for command %s: %v", tt.cmd, err)
			}
			
			// For successful commands, verify we got output
			if !tt.wantErr && err == nil && len(result) == 0 {
				t.Errorf("Command succeeded but returned no output: %s", tt.cmd)
			}
		})
	}
}
