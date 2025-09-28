package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var SUPPORTED_COMPILERS = []string{
	"clang",
	"gcc",
	"zig",
	"cl",
	"bytes",
}

var foundCompailer string
var cwd string
var crunFolderPath string
var exePath string
var filename string

var ulog = &Ulog{0}

func CommandExists(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil // Simplified: directly return the boolean result
}

func getModTime(path string) (time.Time, error) {
	info, err := os.Stat(path)
	if err != nil {
		return time.Time{}, err
	}
	return info.ModTime(), nil
}

func shouldRecompile() bool {
	if flags.noCache {
		return true
	}

	sourceTime, err1 := getModTime(filename)
	exeTime, err2 := getModTime(exePath)

	if err1 != nil {
		return true
	}

	// If binary doesn't exist or source is newer, recompile
	if err2 != nil || sourceTime.After(exeTime) {
		return true
	}
	return false
}

func selectCompiler() {
	if flags.compiler != "" {
		// Check if manually specified compiler exists
		if CommandExists(flags.compiler) {
			foundCompailer = flags.compiler
			return
		} else {
			ulog.println("Specified compiler '%s' not found", flags.compiler)
			os.Exit(1)
		}
	}

	// Auto-detect compiler if not manually specified
	if foundCompailer == "" {
		for _, compailer := range SUPPORTED_COMPILERS {
			if CommandExists(compailer) {
				foundCompailer = compailer
				break
			}
		}
	}
}

func init() {
	var err error
	cwd, err = os.Getwd()
	if err != nil {
		fmt.Println("Failed to get the Current working directory")
		os.Exit(1)
	}

	crunFolderPath = cwd + string(os.PathSeparator) + ".crun"
	err = os.MkdirAll(crunFolderPath, 0755)
	if err != nil {
		fmt.Println("Failed to make the directory")
		os.Exit(1)
	}

	// Auto-detect compiler at startup
	for _, compailer := range SUPPORTED_COMPILERS {
		if CommandExists(compailer) {
			foundCompailer = compailer
			break
		}
	}
}

func runCommand(command string, args []string) bool {
	cmd := exec.Command(command, args...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		ulog.println("Failed to compile the source file", err)
		return false
	}

	return true
}

func compileSourceFile(sourceFilePath string) bool {
	var args []string

	switch foundCompailer {
	case "zig":
		args = []string{"cc", "-o", exePath, sourceFilePath}
	case "cl":
		args = []string{"/Fe:" + exePath, sourceFilePath}
	default:
		args = []string{"-o", exePath, sourceFilePath}
	}

	// Add extra flags if provided
	if flags.extraFlags != "" {
		extraArgs := strings.Fields(flags.extraFlags)
		// Insert extra flags before the source file
		if foundCompailer == "cl" {
			args = append(args[:1], append(extraArgs, args[1:]...)...)
		} else {
			args = append(args[:len(args)-1], append(extraArgs, args[len(args)-1])...)
		}
	}

	return runCommand(foundCompailer, args)
}

func setupExePath(filename string) {
	absFilename, err := filepath.Abs(filename)
	if err != nil {
		ulog.println("Error resolving absolute path: %s", err)
		return
	}

	var exeName string
	if flags.outputName != "" {
		exeName = flags.outputName
		if !strings.HasSuffix(exeName, ".exe") {
			exeName += ".exe"
		}
	} else {
		exeName = strings.TrimSuffix(filepath.Base(absFilename), filepath.Ext(absFilename)) + ".exe"
	}

	var outputDir string
	if flags.outputDir != "" {
		outputDir = flags.outputDir
		// Create output directory if it doesn't exist
		err := os.MkdirAll(outputDir, 0755)
		if err != nil {
			ulog.println("Failed to create output directory: %s", err)
			return
		}
	} else {
		outputDir = crunFolderPath
	}

	exePath = filepath.Join(outputDir, exeName)

	if !filepath.IsAbs(exePath) {
		exePath, err = filepath.Abs(exePath)
		if err != nil {
			ulog.println("Error resolving exePath: %s", err)
		}
	}
}

func clearLastLines(n int) {
	for i := 0; i < n; i++ {
		fmt.Print("\033[1A") // Move cursor up one line
		fmt.Print("\033[2K") // Clear entire line
	}
}

func findSuitableSourceFile() string {
	ulog.println("No file extension provided, trying common extensions...")
	possibleExtensions := []string{".c", ".cpp", ".cc", ".cxx", ".h", ".hpp", ".hh", ".hxx"}
	found := false
	for _, ext := range possibleExtensions {
		if _, err := os.Stat(filename + ext); err == nil {
			filename += ext
			found = true
			break
		}
	}
	if !found {
		ulog.println("No file found with common extensions")
		return ""
	}
	ulog.println("Found file: %s", filename)
	ulog.println("If this is not the intended file, please provide the correct filename with extension.")
	return filename
}

func main() {
	parseFlags()

	// check if running in non windows machine in that case running in external terminal is not supported
	if runtime.GOOS != "windows" {
		flags.runInNewTerminal = false
	} else {
		flags.runInNewTerminal = true
	}

	if len(os.Args) < 2 {
		ulog.println("Usage: crun [flags] <filename>")
		ulog.println("Use -h or --help for help")
		os.Exit(1)
	}

	filename = os.Args[1]

	ulog.println("Provided source file: %s", filename)

	// The script can take name with or without extension
	// so we need to check if filename is given with extension or not
	if filepath.Ext(filename) == "" {
		// No extension provided, try common ones
		findSuitableSourceFile()
	}

	setupExePath(filename)

	if !shouldRecompile() {
		ulog.println("No changes detected, skipping recompilation.")
		runBinary()
		return
	}

	selectCompiler()

	if foundCompailer == "" {
		ulog.println("Sorry no supported compiler found in your system")
		return
	}

	ulog.println("Using compiler: %s", foundCompailer)

	if compileSourceFile(filename) {
		ulog.println("Compiled successfully to: %s", exePath)
		runBinary()
	} else {
		ulog.println("Compilation failed.")
	}

}

func runBinary() {
	ulog.println("Running the binary...")

	ulog.clear()

	var args []string
	if flags.runArgs != "" {
		args = strings.Fields(flags.runArgs)
	}

	if !pathExists(exePath) {
		ulog.println("Executable not found: %s", exePath)
	}

	var err error

	if flags.runInNewTerminal {
		err = LaunchInExternalTerminal(exePath, args...)

	} else {
		cmd := exec.Command(exePath, args...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			ulog.println("Failed to run the binary: %s", err)
		}

	}

	if err != nil {
		fmt.Println("Failed to run the binary")
		return
	}
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
