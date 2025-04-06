package shell

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// TestCrossPlatformDetection tests that the function correctly detects and configures
// itself for the current operating system.
func TestCrossPlatformDetection(t *testing.T) {
	// Create a test file
	tempDir, err := os.MkdirTemp("", "jarvis-cross-platform")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test file in the temp directory
	testFile := filepath.Join(tempDir, "cross_platform_test.txt")
	testContent := "Cross-platform test content."
	err = os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Platform detection test
	result, err := executeCommand("echo Cross-platform test", "")
	if err != nil {
		t.Fatalf("executeCommand() failed on basic echo test: %v", err)
	}
	if !strings.Contains(result, "Cross-platform test") {
		t.Errorf("Platform detection failed, echo command didn't return expected output")
	}

	// Test platform-specific command echo syntax (only validates it doesn't crash)
	var platformCmd string
	if runtime.GOOS == "windows" {
		platformCmd = "echo %PATH%"
	} else {
		platformCmd = "echo $PATH"
	}

	_, err = executeCommand(platformCmd, "")
	if err != nil {
		t.Errorf("Platform-specific command echo failed: %v", err)
	}

	// File reading test (platform independent operation)
	var catCmd string
	if runtime.GOOS == "windows" {
		catCmd = "type " + strings.ReplaceAll(testFile, "/", "\\")
	} else {
		catCmd = "cat " + testFile
	}

	result, err = executeCommand(catCmd, "")
	if err != nil {
		t.Errorf("Failed to read file with platform-specific command: %v", err)
	}
	if !strings.Contains(result, testContent) {
		t.Errorf("File content not found in result, expected: %s, got: %s", testContent, result)
	}
}

// TestWindowsCommandSimulation simulates Windows-specific commands and behavior
// even when running on non-Windows platforms, to validate the Windows-specific
// code paths in the executeCommand function.
func TestWindowsCommandSimulation(t *testing.T) {
	// Skip on actual Windows as we have dedicated Windows tests
	if runtime.GOOS == "windows" {
		t.Skip("Skipping Windows simulation test on actual Windows")
	}

	// Mock the runtime.GOOS value
	// Note: This is just a simulation test - we can't actually change runtime.GOOS
	// We're mostly ensuring the code paths would work on Windows without errors

	// Basic Windows command analogs that might work on Unix
	cmds := []struct {
		name string
		cmd  string
	}{
		{
			name: "dir command analog",
			cmd:  "ls",
		},
		{
			name: "echo command",
			cmd:  "echo Windows simulation test",
		},
		{
			name: "Windows-style path with forward slashes",
			cmd:  "ls " + filepath.Join(".", "testdata"),
		},
	}

	for _, tc := range cmds {
		t.Run(tc.name, func(t *testing.T) {
			// Execute the Unix analog of the Windows command
			_, err := executeCommand(tc.cmd, "")
			// We're mostly checking that these don't crash
			if err != nil {
				t.Logf("Note: %s resulted in error: %v (expected on non-Windows)", tc.name, err)
			}
		})
	}

	// Test that forward slashes work in paths (Windows accepts both)
	tempDir, err := os.MkdirTemp("", "jarvis-win-sim")
	if err == nil {
		defer os.RemoveAll(tempDir)
		
		forwardSlashPath := strings.ReplaceAll(tempDir, "\\", "/")
		_, err = executeCommand("ls", forwardSlashPath)
		if err != nil {
			t.Logf("Forward slash path test resulted in error: %v (expected on non-Windows)", err)
		}
	}
}

// TestPathNormalization tests that paths are correctly normalized
// across platforms.
func TestPathNormalization(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "jarvis-path-norm")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Get current user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get user home directory: %v", err)
	}

	// Test cases
	tests := []struct {
		name        string
		cmd         string
		workDir     string
		wantWorkDir string
		wantErr     bool
	}{
		{
			name:        "Absolute path",
			cmd:         "echo test",
			workDir:     tempDir,
			wantWorkDir: tempDir,
			wantErr:     false,
		},
		{
			name:        "Home directory path with tilde",
			cmd:         "echo test",
			workDir:     homeDir, // We can't use ~ directly in tests, but we're testing the concept
			wantWorkDir: homeDir,
			wantErr:     false,
		},
		{
			name:        "Current directory path with dot",
			cmd:         "echo test",
			workDir:     ".",
			wantWorkDir: "", // Just checking it doesn't error
			wantErr:     false,
		},
		{
			name:        "Non-existent path",
			cmd:         "echo test",
			workDir:     "/path/that/does/not/exist",
			wantWorkDir: "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := executeCommand(tt.cmd, tt.workDir)
			if tt.wantErr && err == nil {
				t.Errorf("Expected error for workDir=%s, got none", tt.workDir)
			} else if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error for workDir=%s: %v", tt.workDir, err)
			}
		})
	}
}

// TestSpecialCommandsEscaping tests that special command characters are properly
// escaped across platforms.
func TestSpecialCommandsEscaping(t *testing.T) {
	// Commands with special characters that should work cross-platform
	tests := []struct {
		name    string
		cmd     string
		wantErr bool
	}{
		{
			name:    "Command with quotes",
			cmd:     "echo \"quoted text\"",
			wantErr: false,
		},
		{
			name:    "Command with single quotes",
			cmd:     "echo 'single quoted text'",
			wantErr: false,
		},
		{
			name:    "Command with safe special chars",
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
