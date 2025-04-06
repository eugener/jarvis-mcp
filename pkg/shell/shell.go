package shell

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// ExecuteCommand executes OS commands specified by the user and returns the command output.
// It takes the command string and an optional working directory, runs the command via the shell,
// and captures both stdout and stderr output.
func executeCommand(cmd string, workDir string) (string, error) {
	command := exec.Command("sh", "-c", cmd)
	command.Env = os.Environ() // Explicitly copy the current environment

	if strings.TrimSpace(workDir) != "" {
		dirInfo, err := os.Stat(workDir)
		if err != nil {
			if os.IsNotExist(err) {
				return "", fmt.Errorf("Path '%s' does not exist\n", workDir)
			}
			return "", fmt.Errorf("Error checking path: %v\n", err)
		}

		if !dirInfo.IsDir() {
			return "", fmt.Errorf("path '%s' exists but is not a directory", workDir)
		}

		command.Dir = workDir
	}

	// Execute command and capture output
	output, err := command.CombinedOutput()
	outputStr := string(output)

	// Format the response
	if err != nil {
		// Return both the error and any output
		resultText := fmt.Sprintf("Command failed: %s\n\nOutput:\n%s\n\nError: %v",
			cmd, outputStr, err)
		return resultText, err
	}

	resultText := fmt.Sprintf("Command executed successfully: %s\n\nOutput:\n%s",
		cmd, outputStr)
	return resultText, nil
}
