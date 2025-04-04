package files

import (
	"fmt"
	"jarvis_mcp/pkg/utils"
	"os"
	"path/filepath"
	"strings"

	"github.com/samber/lo"
)

// normalizePath takes a path string, resolves any home directory references (~),
// converts it to an absolute path, and verifies that the path exists.
// Returns the validated absolute path and any error encountered.
func normalizePath(path string) (string, error) {
	// Handle empty path
	if path == "" {
		return "", os.ErrNotExist
	}

	// Expand home directory reference
	if strings.HasPrefix(path, "~") {
		userDir, err := os.UserHomeDir()
		if err != nil {
			return "", err // Return empty string consistently on error
		}
		path = filepath.Join(userDir, path[1:]) // More reliable than simple string replacement
	}

	// Convert to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err // Return empty string consistently on error
	}

	// Check if the path exists
	_, err = os.Stat(absPath)
	if err != nil {
		return absPath, err // Return the absolute path even if it doesn't exist
	}

	return absPath, nil
}

// readFile reads the content of a file at the given path and returns it as a string.
func readFile(path string) (string, error) {
	// Validate and normalize the file path
	path, err := normalizePath(path)
	if err != nil {
		return "", err
	}

	// Read the file directly into a string to avoid unnecessary conversions
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	// Pre-allocate the string with the exact size of the byte slice
	// This is more efficient than the implicit conversion
	return string(data), nil
}

// writeFile writes the given content to a file at the specified path.
func writeFile(path string, content string) error {

	// Ensure parent directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Write the file with appropriate permissions
	// 0644 = -rw-r--r-- (owner can read/write, group/others can read)
	return os.WriteFile(path, []byte(content), 0644)
}

// readFiles reads the content of multiple files at the given paths and returns them as a slice of strings.
func readFiles(paths []string) ([]string, error) {
	var contents []string
	for _, path := range paths {
		content, err := readFile(path)
		if err != nil {
			return nil, err
		}
		contents = append(contents, content)
	}
	return contents, nil
}

// createDirectory creates a directory at the specified path with appropriate permissions.
// Returns an error if the directory cannot be created.
func createDirectory(path string) error {
	// Create the directory with appropriate permissions
	// 0755 = drwxr-xr-x (owner can read/write/execute, group/others can read/execute)
	return os.Mkdir(path, 0755)
}

// listDirectory lists the contents of the directory at the given path.
// Returns a slice of strings with directory entries prefixed with [DIR] or [FILE].
func listDirectory(path string) ([]string, error) {
	// Validate and normalize the directory path
	path, err := normalizePath(path)
	if err != nil {
		return nil, err
	}

	// List the contents of the directory
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	// Extract the names of the files and directories
	names := lo.Map(files, func(file os.DirEntry, index int) string {
		prefix := utils.IfElse(file.IsDir(), "[DIR]", "[FILE]")
		return fmt.Sprintf("%s %s", prefix, file.Name())
	})

	return names, nil
}

// moveFile moves a file from the source path to the destination path.
// Returns an error if the operation fails.
func moveFile(src, dst string) error {
	// Validate and normalize the source and destination paths
	src, err := normalizePath(src)
	if err != nil {
		return err
	}

	// Move the file from source to destination
	return os.Rename(src, dst)
}

// searchFiles searches for files in the specified directory that match the given pattern.
// Returns a slice of file names that contain the pattern and any error encountered.
func searchFiles(path string, pattern string) ([]string, error) {
	// Validate and normalize the directory path
	path, err := normalizePath(path)
	if err != nil {
		return nil, err
	}

	// List the contents of the directory
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	names := lo.FilterMap(files, func(file os.DirEntry, index int) (string, bool) {
		if strings.Contains(file.Name(), pattern) {
			return file.Name(), true
		}
		return "", false
	})

	return names, nil
}

// getFileInfo returns file information for the given path as a map.
// The map contains the file's name, size, mode, modification time, and whether it's a directory.
func getFileInfo(path string) (map[string]any, error) {
	// Validate and normalize the file path
	path, err := normalizePath(path)
	if err != nil {
		return nil, err
	}

	// Get the file information
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	// Return the file information as a map
	return map[string]any{
		"Name":    info.Name(),
		"Size":    info.Size(),
		"Mode":    info.Mode(),
		"ModTime": info.ModTime(),
		"IsDir":   info.IsDir(),
	}, nil
}
