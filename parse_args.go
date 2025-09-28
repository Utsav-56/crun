package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

// Command line flags
type Flags struct {
	verbose          bool
	noCache          bool
	help             bool
	compiler         string
	extraFlags       string
	outputName       string
	outputDir        string
	runArgs          string
	runInNewTerminal bool
}

var flags Flags

// Flag aliases mapping
var flagAliases = map[string]string{
	"--verbose":         "-v",
	"--recompile":       "-n",
	"--help":            "-h",
	"--compiler":        "-c",
	"--extra":           "-e",
	"--output":          "-o",
	"--directory":       "-d",
	"--run-args":        "-r",
	"--no-new-terminal": "-std",
}

func parseFlags() {

	// Process aliases first
	args := processAliases(os.Args[1:])

	// Parse flags manually
	var nonFlagArgs []string

	for i := 0; i < len(args); i++ {
		arg := args[i]

		switch arg {
		case "-v":
			flags.verbose = true
		case "-n":
			flags.noCache = true
		case "-h":
			flags.help = true
		case "-c":
			if i+1 < len(args) {
				flags.compiler = args[i+1]
				i++
			} else {
				fmt.Println("Error: -c flag requires a compiler name")
				os.Exit(1)
			}
		case "-e":
			if i+1 < len(args) {
				flags.extraFlags = args[i+1]
				i++
			} else {
				fmt.Println("Error: -e flag requires extra flags")
				os.Exit(1)
			}
		case "-o":
			if i+1 < len(args) {
				flags.outputName = args[i+1]
				i++
			} else {
				fmt.Println("Error: -o flag requires output name")
				os.Exit(1)
			}
		case "-d":
			if i+1 < len(args) {
				flags.outputDir = args[i+1]
				i++
			} else {
				fmt.Println("Error: -d flag requires directory path")
				os.Exit(1)
			}
		case "-r":
			if i+1 < len(args) {
				flags.runArgs = args[i+1]
				i++
			} else {
				fmt.Println("Error: -r flag requires run arguments")
				os.Exit(1)
			}
		case "-std":
			flags.runInNewTerminal = false
			ulog.println("-std flag is set so, will not open a new terminal window")
		default:
			if strings.HasPrefix(arg, "-") {
				fmt.Printf("Error: Unknown flag %s\n", arg)
				os.Exit(1)
			}
			nonFlagArgs = append(nonFlagArgs, arg)
		}
	}

	// Set the remaining args back to os.Args for compatibility
	os.Args = append([]string{os.Args[0]}, nonFlagArgs...)

	if flags.help {
		showHelp()
		os.Exit(0)
	}
}

func processAliases(args []string) []string {
	var processedArgs []string

	for _, arg := range args {
		if alias, exists := flagAliases[arg]; exists {
			processedArgs = append(processedArgs, alias)
		} else {
			processedArgs = append(processedArgs, arg)
		}
	}

	return processedArgs
}

func addFlagAlias(longFlag, shortFlag string) {
	flagAliases[longFlag] = shortFlag
}

func showHelp() {
	fmt.Println("crun - Compile and run C/C++ files quickly")
	fmt.Println("\nUsage: crun [flags] <filename>")
	fmt.Println("\nFlags:")
	fmt.Println("  -v, --verbose       Verbose mode - don't clear ulog output")
	fmt.Println("  -n, --recompile     Always recompile source file")
	fmt.Println("  -h, --help          Show this help")
	fmt.Println("  -c, --compiler      Manually choose compiler (clang, gcc, zig, cl, bytes)")
	fmt.Println("  -e, --extra         Extra flags to pass to compiler")
	fmt.Println("  -o, --output        Output binary name")
	fmt.Println("  -d, --directory     Directory to store the binary")
	fmt.Println("  -r, --run-args      Arguments to pass to the binary")

	if runtime.GOOS == "windows" {
		fmt.Println("  -std, --no-new-terminal  Do not open a new terminal window to run the binary (on by default) supported only on Windows")
	}
	fmt.Println("\nSupported compilers:", strings.Join(supportedCompilers, ", "))
	fmt.Println("\nExamples:")
	fmt.Println("  crun main.c")
	fmt.Println("  crun -v main.cpp")
	fmt.Println("  crun --verbose --recompile main.c")
	fmt.Println("  crun -c gcc -e \"-O2 -Wall\" main.c")
	fmt.Println("  crun --compiler clang --extra \"-g -fsanitize=address\" main.c")
	fmt.Println("  crun -o myprogram -d ./bin main.c")
	fmt.Println("  crun -r \"arg1 arg2\" main.c")

	if runtime.GOOS == "windows" {
		fmt.Println("  crun -std main.c")
	}
}
