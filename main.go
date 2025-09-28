package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var (
	supportedCompilers                 = []string{"clang", "gcc", "zig", "cl", "bytes"}
	foundCompiler                      string
	cwd, crunFolder, exePath, filename string
	ulog                               = &Ulog{0}
)

// ---------- Helpers ----------

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func getModTime(path string) time.Time {
	info, err := os.Stat(path)
	if err != nil {
		return time.Time{}
	}
	return info.ModTime()
}

func mustMakeDir(path string) {
	if err := os.MkdirAll(path, 0755); err != nil {
		ulog.println("Failed to create directory %s: %v", path, err)
		os.Exit(1)
	}
}

// ---------- Core logic ----------

func shouldRecompile() bool {
	if flags.noCache {
		return true
	}
	srcTime, exeTime := getModTime(filename), getModTime(exePath)
	return srcTime.IsZero() || exeTime.IsZero() || srcTime.After(exeTime)
}

func detectCompiler(preferred string) string {
	if preferred != "" {
		if commandExists(preferred) {
			return preferred
		}
		ulog.println("Specified compiler '%s' not found", preferred)
		os.Exit(1)
	}

	for _, c := range supportedCompilers {
		if commandExists(c) {
			return c
		}
	}
	return ""
}

func runCommand(cmd string, args ...string) error {
	c := exec.Command(cmd, args...)
	c.Stdin, c.Stdout, c.Stderr = os.Stdin, os.Stdout, os.Stderr
	return c.Run()
}

func compile(source string) bool {
	var args []string
	switch foundCompiler {
	case "zig":
		args = []string{"cc", "-o", exePath, source}
	case "cl":
		args = []string{"/Fe:" + exePath, source}
	default:
		args = []string{"-o", exePath, source}
	}

	if flags.extraFlags != "" {
		extra := strings.Fields(flags.extraFlags)
		if foundCompiler == "cl" {
			args = append(args[:1], append(extra, args[1:]...)...)
		} else {
			args = append(args[:len(args)-1], append(extra, args[len(args)-1])...)
		}
	}

	if err := runCommand(foundCompiler, args...); err != nil {
		ulog.println("Compilation failed: %v", err)
		return false
	}
	return true
}

func setupExePath(src string) {
	absSrc, _ := filepath.Abs(src)

	name := flags.outputName
	if name == "" {
		name = strings.TrimSuffix(filepath.Base(absSrc), filepath.Ext(absSrc))
	}
	if !strings.HasSuffix(name, ".exe") {
		name += ".exe"
	}

	outputDir := flags.outputDir
	if outputDir == "" {
		outputDir = crunFolder
	}
	mustMakeDir(outputDir)

	exePath, _ = filepath.Abs(filepath.Join(outputDir, name))
}

func findSource() string {
	if filepath.Ext(filename) != "" {
		return filename
	}

	ulog.println("No extension provided, trying common ones...")
	for _, ext := range []string{".c", ".cpp", ".cc", ".cxx", ".h", ".hpp", ".hh", ".hxx"} {
		if pathExists(filename + ext) {
			ulog.println("Found: %s", filename+ext)
			return filename + ext
		}
	}
	ulog.println("No matching source file found")
	return ""
}

// ---------- Runner ----------

func runBinary() {
	ulog.println("Running binary...")

	args := []string{}
	if flags.runArgs != "" {
		args = strings.Fields(flags.runArgs)
	}

	if !pathExists(exePath) {
		ulog.println("Executable not found: %s", exePath)
		return
	}

	var err error
	if flags.runInNewTerminal {
		err = LaunchInExternalTerminal(exePath, args...)
	} else {
		ulog.clear()
		err = runCommand(exePath, args...)
	}
	if err != nil {
		ulog.println("Failed to run binary: %v", err)
	}
}

// ---------- Init & Main ----------

func init() {
	var err error
	if cwd, err = os.Getwd(); err != nil {
		fmt.Println("Failed to get working dir:", err)
		os.Exit(1)
	}
	crunFolder = filepath.Join(cwd, ".crun")
	mustMakeDir(crunFolder)
	flags.runInNewTerminal = false
}

func main() {

	parseFlags()

	if len(os.Args) < 2 {
		ulog.println("Usage: crun [flags] <filename>")
		os.Exit(1)
	}

	filename = os.Args[1]
	filename = findSource()
	if filename == "" {
		os.Exit(1)
	}

	setupExePath(filename)

	if !shouldRecompile() {
		ulog.println("No changes detected, skipping recompilation.")
		runBinary()
		return
	}

	foundCompiler = detectCompiler(flags.compiler)
	if foundCompiler == "" {
		ulog.println("No supported compiler found")
		return
	}

	ulog.println("Using compiler: %s", foundCompiler)

	if compile(filename) {
		ulog.println("Compiled successfully: %s", exePath)
		runBinary()
	}
}
