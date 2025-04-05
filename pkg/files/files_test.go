package files

import (
	"encoding/json"
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

func TestDirectoryTree(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a nested directory structure
	subDir1 := filepath.Join(tmpDir, "subdir1")
	subDir2 := filepath.Join(tmpDir, "subdir2")
	nestedDir := filepath.Join(subDir1, "nested")

	// Create directories
	os.Mkdir(subDir1, 0755)
	os.Mkdir(subDir2, 0755)
	os.Mkdir(nestedDir, 0755)

	// Create some files
	os.WriteFile(filepath.Join(tmpDir, "file1.txt"), []byte("content1"), 0644)
	os.WriteFile(filepath.Join(subDir1, "file2.txt"), []byte("content2"), 0644)
	os.WriteFile(filepath.Join(nestedDir, "file3.txt"), []byte("content3"), 0644)
	os.WriteFile(filepath.Join(subDir2, "file4.txt"), []byte("content4"), 0644)

	// Get directory tree using directoryTree function
	treeJSON, err := directoryTree(tmpDir)
	if err != nil {
		t.Errorf("failed to get directory tree: %v", err)
	}

	// Parse the JSON result
	var tree treeNode
	err = json.Unmarshal([]byte(treeJSON), &tree)
	if err != nil {
		t.Errorf("failed to parse directory tree JSON: %v", err)
	}

	// Verify the tree structure
	// 1. Check root properties
	if filepath.Base(tmpDir) != tree.Name {
		t.Errorf("expected root name to be %s, got %s", filepath.Base(tmpDir), tree.Name)
	}
	if tree.Type != "directory" {
		t.Errorf("expected root type to be directory, got %s", tree.Type)
	}
	if len(tree.Children) != 3 { // subdir1, subdir2, file1.txt
		t.Errorf("expected root to have 3 children, got %d", len(tree.Children))
	}

	// 2. Check if all expected files and directories are present
	// The tree should contain 3 directories and 4 files total
	var dirCount, fileCount int
	var countNodesRecursive func(node *treeNode)

	countNodesRecursive = func(node *treeNode) {
		if node.Type == "directory" {
			dirCount++
			for _, child := range node.Children {
				countNodesRecursive(child)
			}
		} else {
			fileCount++
		}
	}

	countNodesRecursive(&tree)

	// Subtract 1 from dirCount to exclude the root directory
	if dirCount-1 != 3 { // subdir1, subdir2, nested
		t.Errorf("expected 3 directories, got %d", dirCount-1)
	}
	if fileCount != 4 { // file1.txt, file2.txt, file3.txt, file4.txt
		t.Errorf("expected 4 files, got %d", fileCount)
	}

	// 3. Check for specific entries
	// Find subdir1 and check if it has file2.txt and nested directory
	var foundSubdir1, foundSubdir2, foundNested bool
	var foundFile1, foundFile2, foundFile3, foundFile4 bool

	for _, child := range tree.Children {
		if child.Name == "subdir1" && child.Type == "directory" {
			foundSubdir1 = true
			for _, subChild := range child.Children {
				if subChild.Name == "file2.txt" && subChild.Type == "file" {
					foundFile2 = true
				}
				if subChild.Name == "nested" && subChild.Type == "directory" {
					foundNested = true
					for _, nestedChild := range subChild.Children {
						if nestedChild.Name == "file3.txt" && nestedChild.Type == "file" {
							foundFile3 = true
						}
					}
				}
			}
		}
		if child.Name == "subdir2" && child.Type == "directory" {
			foundSubdir2 = true
			for _, subChild := range child.Children {
				if subChild.Name == "file4.txt" && subChild.Type == "file" {
					foundFile4 = true
				}
			}
		}
		if child.Name == "file1.txt" && child.Type == "file" {
			foundFile1 = true
		}
	}

	if !foundSubdir1 {
		t.Errorf("expected to find subdir1 directory")
	}
	if !foundSubdir2 {
		t.Errorf("expected to find subdir2 directory")
	}
	if !foundNested {
		t.Errorf("expected to find nested directory")
	}
	if !foundFile1 {
		t.Errorf("expected to find file1.txt")
	}
	if !foundFile2 {
		t.Errorf("expected to find file2.txt")
	}
	if !foundFile3 {
		t.Errorf("expected to find file3.txt")
	}
	if !foundFile4 {
		t.Errorf("expected to find file4.txt")
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
