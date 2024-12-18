package main

import (
	"github.com/kagof/gofunge/internal/cmdline"
	"github.com/kagof/gofunge/internal/debug"
	"github.com/kagof/gofunge/model"
	"io"
	"log"
	"os"
)

func main() {
	l := log.New(os.Stderr, "", 0)
	args, usage, err := cmdline.GetCommandArgs()
	if err != nil {
		usage()
		l.Fatalln(err)
	}
	var befunge BefungeProcessor
	if args.DebugMode {
		befunge = debug.NewDebugger(args.Funge, args.Out, args.In, args.Breakpoints)
	} else {
		if len(args.Breakpoints) > 0 {
			usage()
			l.Fatalln("Breakpoints only supported in debug mode")
		}
		befunge = model.NewBefunge(args.Funge, args.Out, args.In)
	}

	proceed := true
	for proceed {
		proceed, err = befunge.BefungeProcess()
		if err != nil {
			l.Fatalln(err)
		}
	}
	switch oc := args.Out.(type) {
	case io.Closer:
		defer func(oc io.Closer) {
			err := oc.Close()
			if err != nil {
				log.Fatalf("Failed to close output file: %v", err)
			}
		}(oc)
	}
	switch ic := args.In.(type) {
	case io.Closer:
		defer func(ic io.Closer) {
			err := ic.Close()
			if err != nil {
				log.Fatalf("Failed to close input file: %v", err)
			}
		}(ic)
	}
}
