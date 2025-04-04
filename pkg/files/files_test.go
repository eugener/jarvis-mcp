package files

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/samber/lo"
)

func TestNormalizePath(t *testing.T) {

	// Test cases
	tests := []struct {
		name        string
		input       string
		expected    string
		expectError bool
	}{
		{
			name:        "Empty path",
			input:       "",
			expected:    "",
			expectError: true,
		},
		{
			name:        "Home directory reference",
			input:       getFirstUserHomeDir(false),
			expected:    getFirstUserHomeDir(true),
			expectError: false,
		},
		{
			name:        "Absolute path",
			input:       getFirstUserHomeDir(false),
			expected:    getFirstUserHomeDir(true),
			expectError: false,
		},
		{
			name:        "Non-existent path",
			input:       "/nonexistent/path",
			expected:    "/nonexistent/path",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := normalizePath(tt.input)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected an error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("did not expect an error but got: %v", err)
				}
			}
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestReadFile(t *testing.T) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write some content to the temporary file
	content := "Hello, World!"
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Read the content using readFile function
	readContent, err := readFile(tmpFile.Name())
	if err != nil {
		t.Errorf("failed to read file: %v", err)
	}

	// Verify the content
	if readContent != content {
		t.Errorf("expected %v, got %v", content, readContent)
	}
}

func TestWriteFile(t *testing.T) {
	// Create a temporary file path
	tmpFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Define the content to write
	content := "Hello, World!"

	// Write the content using writeFile function
	err = writeFile(tmpFile.Name(), content)
	if err != nil {
		t.Errorf("failed to write file: %v", err)
	}

	// Read the content back to verify
	readContent, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Errorf("failed to read file: %v", err)
	}

	// Verify the content
	if string(readContent) != content {
		t.Errorf("expected %v, got %v", content, string(readContent))
	}
}

func TestCreateDir(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Define the new directory path
	newDirPath := filepath.Join(tmpDir, "newdir")

	// Create the directory using createDirectory function
	err = createDirectory(newDirPath)
	if err != nil {
		t.Errorf("failed to create directory: %v", err)
	}

	// Verify the directory was created
	info, err := os.Stat(newDirPath)
	if err != nil {
		t.Errorf("failed to stat directory: %v", err)
	}
	if !info.IsDir() {
		t.Errorf("expected a directory but got a file")
	}
}

func TestListDir(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create some files and directories in the temporary directory
	os.Mkdir(filepath.Join(tmpDir, "subdir1"), 0755)
	os.Mkdir(filepath.Join(tmpDir, "subdir2"), 0755)
	os.WriteFile(filepath.Join(tmpDir, "file1.txt"), []byte("content1"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "file2.txt"), []byte("content2"), 0644)

	// List the directory using listDirectory function
	list, err := listDirectory(tmpDir)
	if err != nil {
		t.Errorf("failed to list directory: %v", err)
	}

	// Verify the directory contents
	expected := []string{
		"[FILE] file1.txt",
		"[FILE] file2.txt",
		"[DIR] subdir1",
		"[DIR] subdir2",
	}

	// t.Logf("%v", list)

	if !reflect.DeepEqual(list, expected) {
		t.Errorf("expected %v is not equal %v", expected, list)
	}

}

func TestMoveFile(t *testing.T) {
	// Create a temporary source file
	srcFile, err := os.CreateTemp("", "srcfile")
	if err != nil {
		t.Fatalf("failed to create temp source file: %v", err)
	}
	defer os.Remove(srcFile.Name())

	// Write some content to the source file
	content := "Hello, World!"
	if _, err := srcFile.WriteString(content); err != nil {
		t.Fatalf("failed to write to temp source file: %v", err)
	}
	srcFile.Close()

	// Create a temporary directory for the destination
	dstDir, err := os.MkdirTemp("", "dstdir")
	if err != nil {
		t.Fatalf("failed to create temp destination directory: %v", err)
	}
	defer os.RemoveAll(dstDir)

	// Define the destination file path
	dstFile := filepath.Join(dstDir, "dstfile")

	// Move the file using moveFile function
	err = moveFile(srcFile.Name(), dstFile)
	if err != nil {
		t.Errorf("failed to move file: %v", err)
	}

	// Verify the source file no longer exists
	if _, err := os.Stat(srcFile.Name()); !os.IsNotExist(err) {
		t.Errorf("expected source file to be moved, but it still exists")
	}

	// Verify the destination file exists with the correct content
	readContent, err := os.ReadFile(dstFile)
	if err != nil {
		t.Errorf("failed to read destination file: %v", err)
	}
	if string(readContent) != content {
		t.Errorf("expected %v, got %v", content, string(readContent))
	}
}

func TestSearchFiles(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create some files in the temporary directory
	os.WriteFile(filepath.Join(tmpDir, "file1.txt"), []byte("content1"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "file2.txt"), []byte("content2"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "another.txt"), []byte("content3"), 0644)

	// Search for files using searchFiles function
	pattern := "file"
	matchedFiles, err := searchFiles(tmpDir, pattern)
	if err != nil {
		t.Errorf("failed to search files: %v", err)
	}

	// Verify the matched files
	expected := []string{"file1.txt", "file2.txt"}
	if !reflect.DeepEqual(matchedFiles, expected) {
		t.Errorf("expected %v, got %v", expected, matchedFiles)
	}
}

func TestGetFileInfo(t *testing.T) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write some content to the temporary file
	content := "Hello, World!"
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Get file info using getFileInfo function
	info, err := getFileInfo(tmpFile.Name())
	if err != nil {
		t.Errorf("failed to get file info: %v", err)
	}

	// Get the actual file info using os.Stat for comparison
	actualInfo, err := os.Stat(tmpFile.Name())
	if err != nil {
		t.Fatalf("failed to stat temporary file: %v", err)
	}

	// Verify the file info
	expected := map[string]any{
		"Name":    actualInfo.Name(),
		"Size":    actualInfo.Size(),
		"Mode":    actualInfo.Mode(),
		"ModTime": actualInfo.ModTime(),
		"IsDir":   actualInfo.IsDir(),
	}

	if !reflect.DeepEqual(info, expected) {
		t.Errorf("expected %v, got %v", expected, info)
	}
}

// getHomeDir returns the home directory of the current user.
func getHomeDir() string {
	homeDir, _ := os.UserHomeDir()
	return homeDir
}

// getCurrentDir returns the current working directory.
func getCurrentDir() string {
	currentDir, _ := os.Getwd()
	return currentDir
}

func getFirstUserHomeDir(absolute bool) string {

	homeDir := getHomeDir()
	entries, _ := os.ReadDir(homeDir)
	entries = lo.Filter(entries, func(e os.DirEntry, index int) bool {
		return !strings.HasPrefix(e.Name(), ".") && e.IsDir()
	})
	dir := filepath.Join(homeDir, entries[0].Name())
	// println(dir)
	if absolute {
		return dir
	}
	return strings.Replace(dir, homeDir, "~/", 1)
}
