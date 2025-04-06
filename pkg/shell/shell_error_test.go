package shell

import (
	"strings"
	"testing"
)

// TestCommandErrorHandling tests error handling for command execution
func TestCommandErrorHandling(t *testing.T) {
	helper := NewTestHelper(t, "jarvis-error-test")
	defer helper.Cleanup()
	
	// Test cases for command errors
	tests := []struct {
		name       string
		cmd        string
		skipOnWin  bool // Skip this test on Windows
	}{
		{
			name:      "command not found",
			cmd:       "nonexistentcommand",
			skipOnWin: false,
		},
		{
			name:      "permission denied",
			cmd:       "touch /permission_denied_test.txt",
			skipOnWin: true, // Skip on Windows as permissions work differently
		},
		{
			name:      "command with syntax error",
			cmd:       "echo 'missing quote",
			skipOnWin: false,
		},
	}
	
	for _, tt := range tests {
		// Skip tests that shouldn't run on Windows
		if tt.skipOnWin && IsWindows() {
			t.Logf("Skipping %s on Windows", tt.name)
			continue
		}
		
		t.Run(tt.name, func(t *testing.T) {
			result, err := executeCommand(tt.cmd, helper.TempDir)
			
			// Command execution errors should be propagated
			if err == nil {
				t.Errorf("executeCommand() should return error for %s", tt.name)
				return
			}
			
			// Result should contain error indication
			if !strings.Contains(result, "Command failed") {
				t.Errorf("executeCommand() result doesn't contain error indication for %s", tt.name)
			}
		})
	}
}

// TestCommandPipes tests error handling for pipe operations
func TestCommandPipes(t *testing.T) {
	helper := NewTestHelper(t, "jarvis-pipes-test")
	defer helper.Cleanup()
	
	// Define find command based on platform
	findCmd := "grep"
	if IsWindows() {
		findCmd = "find"
	}
	
	// Unix shells may not propagate errors from the first command in a pipe
	// if the second command succeeds, so we need to handle this differently
	// on different platforms
	failPipeCmd := ""
	if IsWindows() {
		failPipeCmd = "nonexistentcommand | echo 'fallback'"
	} else {
		// On Unix, use a command that will fail even with piping
		failPipeCmd = "nonexistentcommand && echo 'this will not run'"
	}
	
	// Test different pipe scenarios
	tests := []struct {
		name      string
		cmd       string
		wantErr   bool
	}{
		{
			name:    "successful pipe",
			cmd:     "echo 'test' | " + findCmd + " 'test'",
			wantErr: false,
		},
		{
			name:    "failing command",
			cmd:     failPipeCmd,
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := executeCommand(tt.cmd, helper.TempDir)
			
			if tt.wantErr && err == nil {
				t.Errorf("executeCommand() should return error for %s", tt.name)
			} else if !tt.wantErr && err != nil {
				t.Errorf("executeCommand() unexpected error for %s: %v", tt.name, err)
			}
			
			// For error cases, verify error message in result
			if tt.wantErr && err != nil && !strings.Contains(result, "Command failed") {
				t.Errorf("executeCommand() result doesn't contain error indication")
			}
		})
	}
}

// TestSpecialCommandsEscaping tests escape handling in commands
func TestSpecialCommandsEscaping(t *testing.T) {
	helper := NewTestHelper(t, "jarvis-escaping-test")
	defer helper.Cleanup()
	
	// Commands with special characters that should work cross-platform
	tests := []struct {
		name    string
		cmd     string
		wantErr bool
	}{
		{
			name:    "command with double quotes",
			cmd:     "echo \"quoted text\"",
			wantErr: false,
		},
		{
			name:    "command with single quotes",
			cmd:     "echo 'single quoted text'",
			wantErr: false,
		},
		{
			name:    "command with safe special chars",
			cmd:     "echo special chars",
			wantErr: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := executeCommand(tt.cmd, "")
			if tt.wantErr && err == nil {
				t.Errorf("Expected error for cmd=%s, got none", tt.cmd)
			} else if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error for cmd=%s: %v", tt.cmd, err)
			}
		})
	}
}
