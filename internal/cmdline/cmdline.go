package cmdline

import (
	"errors"
	"flag"
	"fmt"
	"github.com/kagof/gofunge/model"
	"io"
	"os"
	"strings"
)

type CommandArgs struct {
	DebugMode   bool
	Breakpoints []model.Vector2
	Out         io.Writer
	In          io.Reader
	Funge       string
}

func GetCommandArgs() (*CommandArgs, func(), error) {
	var debugMode bool
	var outputFile string
	var inputFile string
	var inline string
	var breakpoints breakpointList
	flag.BoolVar(&debugMode, "debug", false, "toggle on/off debug mode")
	flag.Var(&breakpoints, "breakpoint", "Breakpoints In the program, if In debug mode. Multiple supported")
	flag.StringVar(&outputFile, "o", "", "output file (default: stdout)")
	flag.StringVar(&inputFile, "i", "", "input file (default: stdin)")
	flag.StringVar(&inline, "inline", "", "an inline Befunge-93 program")
	flag.Usage = func() {
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "%s [OPTIONS] befungeFile\n\n", os.Args[0])
		_, _ = fmt.Fprintln(flag.CommandLine.Output(), "Examples:")
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "\t %s filename.bf\n", os.Args[0])
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "\t %s -debug -breakpoint '(0,0)' filename.bf\n", os.Args[0])
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "\t %s -debug -breakpoint 0,0 -breakpoint 1,2 filename.bf\n", os.Args[0])
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "\t %s -i input.txt -o output.txt filename.bf\n", os.Args[0])
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "\t %s -inline '\"olleh\",,,,,@'\n", os.Args[0])
		_, _ = fmt.Fprintln(flag.CommandLine.Output())
		flag.PrintDefaults()
	}
	flag.Parse()
	var err error

	if flag.NArg() < 1 && inline == "" {
		return nil, flag.Usage, errors.New("argument is required if no inline program provided")
	}
	var funge string
	if inline != "" {
		funge = inline
	} else {
		var file []byte
		file, err = os.ReadFile(flag.Arg(0))
		if err != nil {
			return nil, flag.Usage, errors.Join(errors.New(fmt.Sprintf("Cannot read file %s\n", flag.Arg(0))), err)
		}
		funge = string(file)
	}

	var out io.Writer
	if outputFile != "" {
		// Open the file (create if it doesn't exist, with write-only access)
		out, err = os.Create(outputFile)
		if err != nil {
			return nil, flag.Usage, errors.Join(errors.New(fmt.Sprintf("Cannot read output file %s\n", flag.Arg(0))), err)
		}
	} else {
		out = os.Stdout
	}
	var in io.Reader
	if inputFile != "" {
		// Open the file (create if it doesn't exist, with write-only access)
		in, err = os.Open(inputFile)
		if err != nil {
			return nil, flag.Usage, errors.Join(errors.New(fmt.Sprintf("Cannot read input file %s\n", flag.Arg(0))), err)
		}
	} else {
		in = os.Stdin

	}
	return &CommandArgs{
		DebugMode:   debugMode,
		Breakpoints: breakpoints,
		Out:         out,
		In:          in,
		Funge:       funge,
	}, flag.Usage, nil
}

type breakpointList []model.Vector2

func (b *breakpointList) String() string {
	stringSlice := make([]string, len(*b))
	for _, bp := range *b {
		stringSlice = append(stringSlice, bp.String())
	}
	return strings.Join(stringSlice, ", ")
}

func (b *breakpointList) Set(value string) error {
	bp, err := model.ParseVector2(value)
	if err != nil {
		return err
	}
	*b = append(*b, *bp)
	return nil
}
