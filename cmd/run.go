package cmd

import (
	"github.com/kagof/kagofunge/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var runCmd = &cobra.Command{
	Use:     "run <program>",
	Aliases: []string{"execute"},
	Short:   "Run a Befunge-93 program",
	Example: `kagofunge run hello-world.bf
kagofunge run '<> #,:# _@#:"Hello, World!"' -I
kagofunge run hello-world.bf -o output.txt -i input.txt`,
	DisableAutoGenTag: true,
	Long: `run will execute a Befunge-93 program, exiting when either the program 
terminates or an unhandled error occurs.`,
	Args: cobra.ExactArgs(1),
	RunE: runRunE,
}

func runRunE(cmd *cobra.Command, args []string) error {
	flags := *cmd.Flags()

	befunge, err := getBefunge(flags, args)
	if err != nil {
		return err
	}
	hasNext := true
	for hasNext {
		hasNext, err = befunge.Step()
		if err != nil {
			return err
		}
	}

	return nil
}

func getBefunge(flags pflag.FlagSet, args []string) (pkg.Stepper, error) {
	program, outputFile, inputFile, err := getGlobals(flags, args)
	if err != nil {
		return nil, err
	}
	befunge := pkg.NewBefunge(program, outputFile, inputFile)
	return befunge, nil
}

func init() {
	rootCmd.AddCommand(runCmd)
}
