//go:build windows
// +build windows

package shell

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestWindowsSpecificCommands tests Windows-specific command execution
// This test will only run on Windows due to the build tag
func TestWindowsSpecificCommands(t *testing.T) {
	helper := NewTestHelper(t, "jarvis-mcp-windows-test")
	defer helper.Cleanup()
	
	// Create a test file
	testFile := helper.CreateTestFile("test.txt", "Hello, JARVIS MCP on Windows!")
	
	// Windows-specific test cases
	tests := []struct {
		name        string
		cmd         string
		workDir     string
		wantErr     bool
		wantContain string
	}{
		{
			name:        "dir command",
			cmd:         "dir",
			workDir:     helper.TempDir,
			wantErr:     false,
			wantContain: "test.txt",
		},
		{
			name:        "type command (Windows equivalent of cat)",
			cmd:         "type test.txt",
			workDir:     helper.TempDir,
			wantErr:     false,
			wantContain: "Hello, JARVIS MCP on Windows!",
		},
		{
			name:        "echo with Windows environment variable",
			cmd:         "echo %TEMP%",
			workDir:     "",
			wantErr:     false,
			wantContain: "", // Just check it doesn't error
		},
		{
			name:        "Windows command with pipe",
			cmd:         "dir | find \"test\"",
			workDir:     helper.TempDir,
			wantErr:     false,
			wantContain: "test.txt",
		},
		{
			name:        "tasklist command",
			cmd:         "tasklist",
			workDir:     "",
			wantErr:     false,
			wantContain: "System", // System process should always be present
		},
		{
			name:        "whoami command",
			cmd:         "whoami",
			workDir:     "",
			wantErr:     false,
			wantContain: "", // Just check it doesn't error
		},
		{
			name:        "command with backslash path",
			cmd:         "echo Hello > output.txt && type output.txt",
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

// TestWindowsPathHandling tests Windows-specific path handling
func TestWindowsPathHandling(t *testing.T) {
	// Test Windows paths with backslashes
	testCases := []struct {
		name    string
		workDir string
		cmd     string
		wantErr bool
	}{
		{
			name:    "Windows path with backslashes",
			workDir: "C:\\Windows\\Temp",
			cmd:     "dir",
			wantErr: false,
		},
		{
			name:    "Windows path with mixed slashes",
			workDir: "C:/Windows/Temp",
			cmd:     "dir",
			wantErr: false,
		},
		{
			name:    "UNC path",
			workDir: "", // Don't use UNC path in workDir as it may not exist
			cmd:     "dir C:\\Windows",
			wantErr: false,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := executeCommand(tc.cmd, tc.workDir)
			if tc.wantErr && err == nil {
				t.Errorf("Expected error for workDir=%s, got none", tc.workDir)
			} else if !tc.wantErr && err != nil {
				t.Errorf("Unexpected error for workDir=%s: %v", tc.workDir, err)
			}
		})
	}
}

// TestWindowsCommandCharacters tests Windows-specific command characters and syntax
func TestWindowsCommandCharacters(t *testing.T) {
	helper := NewTestHelper(t, "jarvis-windows-chars")
	defer helper.Cleanup()
	
	tests := []struct {
		name    string
		cmd     string
		wantErr bool
	}{
		{
			name:    "Windows redirect output",
			cmd:     "echo test > test.txt && type test.txt",
			wantErr: false,
		},
		{
			name:    "Windows append output",
			cmd:     "echo line1 > test2.txt && echo line2 >> test2.txt && type test2.txt",
			wantErr: false,
		},
		{
			name:    "Windows AND operator",
			cmd:     "echo first && echo second",
			wantErr: false,
		},
		{
			name:    "Windows OR operator",
			cmd:     "nosuchcommand || echo fallback",
			wantErr: false, // Should succeed because of the fallback
		},
		{
			name:    "Windows special characters",
			cmd:     "echo %TEMP%",
			wantErr: false,
		},
		{
			name:    "Windows caret escape character",
			cmd:     "echo This is a test ^& symbol",
			wantErr: false,
		},
		{
			name:    "Windows batch file commands",
			cmd:     "set testvar=hello && echo %testvar%",
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
