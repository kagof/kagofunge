package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"io"
	"os"
	"strings"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kagofunge <run | debug> <program> [flags]",
	Short: "A Befunge-93 interpreter and debugger",
	Example: `kagofunge run hello-world.bf
kagofunge run '<> #,:# _@#:"Hello, World!"' -I
kagofunge run hello-world.bf -o output.txt -i input.txt

kagofunge debug hello-world.bf --breakpoint "(0,0)"
kagofunge debug '<> #,:# _@#:"Hello, World!"' -I -b 0,0 -b 8,0
kagofunge debug hello-world.bf -o output.txt -i input.txt -b '[1,1]'`,
	Version: "0.1.0",
	Long: `kagofunge is an interpreter and debugger for Befunge-93 written in Go.
For detailed usage, use kagofunge run --help or kagofunge debug --help.`,
	DisableAutoGenTag:     true,
	DisableFlagsInUseLine: true,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	rootCmd.SetUsageTemplate(strings.Replace(rootCmd.UsageTemplate(),
		"{{.CommandPath}} [command]",
		"{{.CommandPath}} <run|debug> <program> [flags]",
		1))

	rootCmd.PersistentFlags().StringP("output", "o", "", "Output file path. Default: stdout")
	rootCmd.PersistentFlags().StringP("input", "i", "", "Output file path. Default: stdin")
	rootCmd.PersistentFlags().BoolP("inline",
		"I",
		false,
		`If set, then the <program> is interpreted as an inline 
Befunge-93 program, otherwise it is interpreted as a 
path to a Befunge-93 program file.`)
	err := rootCmd.MarkPersistentFlagFilename("output")
	if err != nil {
		panic(err)
	}
	err = rootCmd.MarkPersistentFlagFilename("input")
	if err != nil {
		panic(err)
	}
}

func getOutputFile(flags pflag.FlagSet) (io.Writer, error) {
	var out io.WriteCloser
	flag, err := flags.GetString("output")
	if err != nil {
		out = nil
	} else if flag == "" {
		out = os.Stdout
	} else {
		// Open the file (create if it doesn't exist, with write-only access)
		out, err = os.Create(flag)
		if err != nil {
			out = nil
			err = errors.Join(errors.New(fmt.Sprintf("Cannot read output file %s", flag)), err)
		}
	}
	return out, err
}

func getInputFile(flags pflag.FlagSet) (io.Reader, error) {
	var in io.ReadCloser
	flag, err := flags.GetString("input")
	if err != nil {
		in = nil
	} else if flag == "" {
		in = os.Stdin
	} else {
		// Open the file (create if it doesn't exist, with write-only access)
		in, err = os.Open(flag)
		if err != nil {
			in = nil
			err = errors.Join(errors.New(fmt.Sprintf("Cannot read input file %s", flag)), err)
		}
	}
	return in, err
}

func getProgram(flags pflag.FlagSet, arg string) (string, error) {
	inline, err := flags.GetBool("inline")
	if err != nil {
		return "", err
	}
	var program string
	if inline {
		program = arg
	} else {
		var file []byte
		file, err = os.ReadFile(arg)
		if err != nil {
			return "", errors.Join(errors.New(fmt.Sprintf("Cannot read file %s", arg)), err)
		}
		program = string(file)
	}
	return program, nil
}

func getGlobals(flags pflag.FlagSet, args []string) (string, io.Writer, io.Reader, error) {
	outputFile, err := getOutputFile(flags)
	if err != nil {
		return "", nil, nil, err
	}

	inputFile, err2 := getInputFile(flags)
	if err2 != nil {
		return "", nil, nil, err2
	}

	program, err3 := getProgram(flags, args[0])
	if err3 != nil {
		return "", nil, nil, err3
	}

	return program, outputFile, inputFile, nil
}
