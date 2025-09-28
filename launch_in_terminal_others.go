//go:build !windows
// +build !windows

package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// LaunchInExternalTerminal launches a binary in a new terminal window with optional arguments.
// After execution, the terminal will prompt "Press Enter to exit..." and close once user presses Enter.
func LaunchInExternalTerminal(binaryPath string, args ...string) error {
	quotedBinary := fmt.Sprintf("\"%s\"", binaryPath)

	// Quote args
	var quotedArgs []string
	for _, arg := range args {
		quotedArgs = append(quotedArgs, fmt.Sprintf("\"%s\"", arg))
	}
	argLine := strings.Join(quotedArgs, " ")

	switch runtime.GOOS {
	case "darwin":
		// macOS: use AppleScript to run in Terminal.app
		// The "read -n 1" ensures it waits for user input before closing
		cmdStr := fmt.Sprintf("'%s %s; echo; echo Press Enter to exit...; read -n 1' ", binaryPath, argLine)
		script := fmt.Sprintf(`tell application "Terminal"
			activate
			do script %s
		end tell`, cmdStr)

		return exec.Command("osascript", "-e", script).Start()

	default:
		// Linux / Unix
		// Try common terminals
		terminalCandidates := []string{"gnome-terminal", "konsole", "xterm", "lxterminal", "xfce4-terminal"}
		var term string
		for _, t := range terminalCandidates {
			if _, err := exec.LookPath(t); err == nil {
				term = t
				break
			}
		}
		if term == "" {
			return fmt.Errorf("no supported terminal emulator found")
		}

		// Command that pauses
		shCmd := fmt.Sprintf("%s %s; echo; echo Press Enter to exit...; read -n 1", binaryPath, argLine)

		var cmd *exec.Cmd
		switch term {
		case "gnome-terminal", "xfce4-terminal", "lxterminal":
			cmd = exec.Command(term, "--", "bash", "-c", shCmd)
		case "konsole":
			cmd = exec.Command(term, "-e", "bash", "-c", shCmd)
		default: // fallback (xterm, etc.)
			cmd = exec.Command(term, "-e", "bash", "-c", shCmd)
		}
		return cmd.Start()
	}
}
