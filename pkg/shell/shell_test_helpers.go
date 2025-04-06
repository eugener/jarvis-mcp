package shell

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// TestHelper contains common utilities for shell tests
type TestHelper struct {
	T         *testing.T
	TempDir   string
	TestFiles map[string]string // Filename -> Content
}

// NewTestHelper creates a new test helper with a temporary directory
func NewTestHelper(t *testing.T, dirPrefix string) *TestHelper {
	if dirPrefix == "" {
		dirPrefix = "jarvis-mcp-test"
	}
	
	tempDir, err := os.MkdirTemp("", dirPrefix)
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	
	return &TestHelper{
		T:         t,
		TempDir:   tempDir,
		TestFiles: make(map[string]string),
	}
}

// Cleanup removes the temporary directory
func (h *TestHelper) Cleanup() {
	os.RemoveAll(h.TempDir)
}

// CreateTestFile creates a test file with the given content
func (h *TestHelper) CreateTestFile(name, content string) string {
	filePath := filepath.Join(h.TempDir, name)
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		h.T.Fatalf("Failed to create test file: %v", err)
	}
	h.TestFiles[name] = content
	return filePath
}

// AssertCommandSuccess checks that a command succeeds and contains expected output
func (h *TestHelper) AssertCommandSuccess(cmd, workDir, expectedOutput string) {
	result, err := executeCommand(cmd, workDir)
	if err != nil {
		h.T.Errorf("executeCommand(%s) unexpected error: %v", cmd, err)
		return
	}
	
	if expectedOutput != "" && !strings.Contains(result, expectedOutput) {
		h.T.Errorf("executeCommand(%s) result doesn't contain %q, got: %q", cmd, expectedOutput, result)
	}
}

// AssertCommandError checks that a command fails with expected error message
func (h *TestHelper) AssertCommandError(cmd, workDir, expectedErrMsg string) {
	_, err := executeCommand(cmd, workDir)
	if err == nil {
		h.T.Errorf("executeCommand(%s) expected error but got none", cmd)
		return
	}
	
	if expectedErrMsg != "" && !strings.Contains(err.Error(), expectedErrMsg) {
		h.T.Errorf("executeCommand(%s) error doesn't contain %q, got: %v", cmd, expectedErrMsg, err)
	}
}

// PlatformCommands contains platform-specific command syntax
type PlatformCommands struct {
	ListDir       string // "ls" or "dir"
	ReadFile      string // "cat" or "type"
	EnvVarSyntax  string // "$VAR" or "%VAR%"
	PathSeparator string // "/" or "\"
}

// GetPlatformCommands returns the appropriate commands for the current platform
func GetPlatformCommands() PlatformCommands {
	if runtime.GOOS == "windows" {
		return PlatformCommands{
			ListDir:       "dir",
			ReadFile:      "type",
			EnvVarSyntax:  "%%PATH%%", // Escaped for string formatting
			PathSeparator: "\\",
		}
	}
	return PlatformCommands{
		ListDir:       "ls",
		ReadFile:      "cat",
		EnvVarSyntax:  "$PATH",
		PathSeparator: "/",
	}
}

// IsWindows returns true if running on Windows
func IsWindows() bool {
	return runtime.GOOS == "windows"
}
