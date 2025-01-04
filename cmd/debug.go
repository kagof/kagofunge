package cmd

import (
	"github.com/kagof/kagofunge/internal"
	"github.com/kagof/kagofunge/internal/debug"
	"github.com/kagof/kagofunge/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"time"
)

var debugCmd = &cobra.Command{
	Use:   "debug <program>",
	Short: "Debug a Befunge-93 program",
	Example: `kagofunge debug hello-world.bf --breakpoint "(0,0)"
kagofunge debug '<> #,:# _@#:"Hello, World!"' -I -b 0,0 -b 8,0 -I -b 0,0 -b 15,0
kagofunge debug hello-world.bf -o output.txt -i input.txt -b '[1,1]'`,
	Long: `debug will execute a Befunge-93 program, exiting when either the program 
terminates or an unhandled error occurs.

Breakpoints can be set using the --breakpoint/-b flag, which will interrupt the 
program's execution and display information about the current state of the 
program to the caller.'`,
	Args:              cobra.ExactArgs(1),
	DisableAutoGenTag: true,
	RunE:              debugRunE,
}

func debugRunE(cmd *cobra.Command, args []string) error {
	flags := *cmd.Flags()

	befunge, err := getDebugger(flags, args)
	if err != nil {
		return err
	}
	cmd.SilenceUsage = true // don't print usage for errors past this point
	hasNext := true
	for hasNext {
		hasNext, err = befunge.Step()
		if err != nil {
			return err
		}
	}

	return nil
}

func getDebugger(flags pflag.FlagSet, args []string) (pkg.Stepper, error) {
	config, program, outputFile, inputFile, err := getGlobals(flags, args)
	if err != nil {
		return nil, err
	}
	var breakpoints []pkg.Vector2
	breakpoints, err = getBreakpoints(flags)
	if err != nil {
		return nil, err
	}
	var speed time.Duration
	speed, err = flags.GetDuration("speed")
	if err != nil {
		return nil, err
	}

	befunge := debug.NewDebugger(config, program, outputFile, inputFile, breakpoints, speed)
	return befunge, nil
}

func getBreakpoints(flags pflag.FlagSet) ([]pkg.Vector2, error) {
	breakpointStrings, err := flags.GetStringArray("breakpoint")
	if err != nil {
		return nil, err
	}
	return internal.MapSlicePE(breakpointStrings, func(str string) (*pkg.Vector2, error) {
		vec, err := pkg.ParseVector2(str)
		return vec, err
	})
}

func init() {
	rootCmd.AddCommand(debugCmd)
	debugCmd.Flags().StringArrayP("breakpoint",
		"b",
		nil,
		`Breakpoints to set in the program while 
executing. can be in the formats (x,y), (x y), 
[x,y], [x y], or x,y.`)
	debugCmd.Flags().DurationP("speed",
		"s",
		0,
		`If set, the program will progress automatically
at the specified speed. Should be a duration. 
Eg 100ms, 1s`)
}
