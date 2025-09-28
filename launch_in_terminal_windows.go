//go:build windows
// +build windows

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// LaunchInExternalTerminal launches a binary in a new terminal window with optional arguments.
func LaunchInExternalTerminal(binaryPath string, args ...string) error {
	// Quote the binary path if it contains spaces
	quotedBinary := fmt.Sprintf("\"%s\"", binaryPath)

	// Join arguments safely
	argLine := ""
	if len(args) > 0 {
		for i, arg := range args {
			args[i] = fmt.Sprintf("\"%s\"", arg)
		}
		argLine = " " + strings.Join(args, " ")
	}

	// Final command to run inside new terminal
	commandLine := fmt.Sprintf("start \"\" cmd /c \"%s%s & echo Press any key to exit... & pause > nul\"", quotedBinary, argLine)

	// Create temporary .bat file
	batContent := "@echo off\n" + commandLine + "\n"
	tmpBat := filepath.Join(os.TempDir(), "launch_external_terminal.bat")

	err := os.WriteFile(tmpBat, []byte(batContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to write bat file: %w", err)
	}

	// Launch the .bat file
	return exec.Command("cmd.exe", "/C", tmpBat).Start()
}
