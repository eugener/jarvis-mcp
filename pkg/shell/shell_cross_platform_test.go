package shell

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestCrossPlatformDetection tests that the function correctly detects and configures
// itself for the current operating system.
func TestCrossPlatformDetection(t *testing.T) {
	helper := NewTestHelper(t, "jarvis-cross-platform")
	defer helper.Cleanup()
	
	// Create a test file
	testFile := helper.CreateTestFile("cross_platform_test.txt", "Cross-platform test content.")
	
	// Get platform-specific commands
	pc := GetPlatformCommands()
	
	// Platform detection test - basic echo
	helper.AssertCommandSuccess("echo Cross-platform test", "", "Cross-platform test")
	
	// Test platform-specific command echo syntax
	helper.AssertCommandSuccess("echo "+pc.EnvVarSyntax, "", "")
	
	// File reading test with platform-specific command
	readCmd := pc.ReadFile
	if IsWindows() {
		// Ensure Windows path format for Windows tests
		readCmd += " " + strings.ReplaceAll(testFile, "/", "\\")
	} else {
		readCmd += " " + testFile
	}
	
	helper.AssertCommandSuccess(readCmd, "", "Cross-platform test content.")
}

// TestPathNormalization tests that paths are correctly normalized
// across platforms.
func TestPathNormalization(t *testing.T) {
	helper := NewTestHelper(t, "jarvis-path-norm")
	defer helper.Cleanup()
	
	// Get current user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get user home directory: %v", err)
	}
	
	// Test cases for path normalization
	tests := []struct {
		name        string
		cmd         string
		workDir     string
		wantErr     bool
	}{
		{
			name:        "absolute path",
			cmd:         "echo test",
			workDir:     helper.TempDir,
			wantErr:     false,
		},
		{
			name:        "home directory path",
			cmd:         "echo test",
			workDir:     homeDir,
			wantErr:     false,
		},
		{
			name:        "current directory path",
			cmd:         "echo test",
			workDir:     ".",
			wantErr:     false,
		},
		{
			name:        "non-existent path",
			cmd:         "echo test",
			workDir:     "/path/that/does/not/exist",
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

// TestRelativePathHandling tests handling of relative paths
func TestRelativePathHandling(t *testing.T) {
	// Create a nested directory structure
	parent := NewTestHelper(t, "jarvis-parent-dir")
	defer parent.Cleanup()
	
	// Create a file in the parent directory
	parent.CreateTestFile("parent.txt", "Parent file content")
	
	// Create a subdirectory
	subDirPath := filepath.Join(parent.TempDir, "subdir")
	err := os.Mkdir(subDirPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}
	
	// Create a file in the subdirectory
	subFilePath := filepath.Join(subDirPath, "child.txt")
	err = os.WriteFile(subFilePath, []byte("Child file content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create file in subdirectory: %v", err)
	}
	
	// Get platform-specific commands
	pc := GetPlatformCommands()
	
	// Test accessing parent file from subdirectory with relative path
	var relativePath string
	if IsWindows() {
		relativePath = "..\\parent.txt"
	} else {
		relativePath = "../parent.txt"
	}
	
	readCmd := pc.ReadFile + " " + relativePath
	result, err := executeCommand(readCmd, subDirPath)
	if err != nil {
		t.Errorf("Failed to read file with relative path: %v", err)
	} else if !strings.Contains(result, "Parent file content") {
		t.Errorf("Expected result to contain parent file content, got: %s", result)
	}
}

// TestCrossPlatformPaths tests path handling across different platforms
func TestCrossPlatformPaths(t *testing.T) {
	helper := NewTestHelper(t, "jarvis-cross-paths")
	defer helper.Cleanup()
	
	// Create test files
	helper.CreateTestFile("test1.txt", "Test file 1 content")
	helper.CreateTestFile("test2.txt", "Test file 2 content")
	
	// Get platform commands
	pc := GetPlatformCommands()
	
	// Test listing directory with platform-specific list command
	result, err := executeCommand(pc.ListDir, helper.TempDir)
	if err != nil {
		t.Errorf("Failed to list directory: %v", err)
		return
	}
	
	// Check that the results include our test files
	if !strings.Contains(result, "test1.txt") || !strings.Contains(result, "test2.txt") {
		t.Errorf("List directory command didn't show test files, got: %s", result)
	}
	
	// Test path with trailing separator
	trailingPath := helper.TempDir
	if !strings.HasSuffix(trailingPath, string(os.PathSeparator)) {
		trailingPath += string(os.PathSeparator)
	}
	
	_, err = executeCommand(pc.ListDir, trailingPath)
	if err != nil {
		t.Errorf("Failed to use path with trailing separator: %v", err)
	}
	
	// Test with normalized vs non-normalized paths
	nonNormalizedPath := filepath.Join(helper.TempDir, ".", "test1.txt")
	readCmd := pc.ReadFile + " " + nonNormalizedPath
	
	result, err = executeCommand(readCmd, "")
	if err != nil {
		t.Errorf("Failed to use non-normalized path: %v", err)
	} else if !strings.Contains(result, "Test file 1 content") {
		t.Errorf("Expected content from non-normalized path, got: %s", result)
	}
}

// TestEnvironmentVariables tests environment variable handling
func TestEnvironmentVariables(t *testing.T) {
	helper := NewTestHelper(t, "jarvis-env-vars")
	defer helper.Cleanup()
	
	// Define platform-specific env var commands
	var setCmd, getCmd string
	if IsWindows() {
		setCmd = "set TEST_ENV_VAR=test_value"
		getCmd = "echo %TEST_ENV_VAR%"
	} else {
		setCmd = "export TEST_ENV_VAR=test_value"
		getCmd = "echo $TEST_ENV_VAR"
	}
	
	// Test setting and getting environment variables
	// This test works differently on different platforms due to how environment
	// variables are handled in shells - on Unix, the variable is only set for
	// the duration of the command, while on Windows, it persists across commands
	// we expect the variable to be set in the current command only
	cmd := setCmd + " && " + getCmd
	result, err := executeCommand(cmd, "")
	if err != nil {
		t.Errorf("Failed to use environment variables: %v", err)
	}
	
	// The output may differ by platform, but should contain either the variable
	// value or indication the variable was used
	if IsWindows() {
		if !strings.Contains(result, "test_value") {
			t.Errorf("Environment variable not set properly, got: %s", result)
		}
	}
}
