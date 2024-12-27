package debug

import (
	"bufio"
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/kagof/kagofunge/internal"
	"github.com/kagof/kagofunge/pkg"
	"io"
	"os"
	"slices"
	"strings"
	"time"
	"unicode"
)

const clearAndReturn = "\033[2J\033[H"

var (
	bold                      = color.New(color.Bold)
	red                       = color.New(color.FgRed)
	redBg                     = color.New(color.BgRed)
	green                     = color.New(color.FgGreen)
	cyan                      = color.New(color.FgCyan)
	boldAndUnderlined         = color.New(color.Bold, color.Underline)
	redAndBoldAndUnderlined   = color.New(color.FgRed, color.Bold, color.Underline)
	redBgAndBoldAndUnderlined = color.New(color.BgRed, color.Bold, color.Underline)
	faint                     = color.New(color.Faint)
)

type Debugger struct {
	befunge     *pkg.Befunge
	breakpoints []pkg.Vector2
	output      *strings.Builder
	outfile     io.Writer
	autoSpeed   time.Duration
	reader      bufio.Reader
	stepMode    bool
	jumping     bool
	hasPrinted  bool
	isStarted   bool
	isFinished  bool
	stdinChan   chan string
	usingChanR  bool
}

func NewDebugger(s string, outFile io.Writer, inFile io.Reader, breakpoints []pkg.Vector2, speed time.Duration) *Debugger {
	b := new(strings.Builder)
	stdinChan := make(chan string)
	var fungeIn io.Reader
	var usingChanR bool
	// if the input file is the same as the debugger input, we have to read from the channel for both
	if inFile == os.Stdin {
		fungeIn = &chanReader{stdinChan}
		usingChanR = true
	} else {
		fungeIn = inFile
		usingChanR = false
	}

	return &Debugger{
		befunge:     pkg.NewBefunge(s, b, fungeIn),
		breakpoints: breakpoints,
		output:      b,
		outfile:     outFile,
		reader:      *bufio.NewReader(os.Stdin),
		autoSpeed:   speed,
		stdinChan:   stdinChan,
		usingChanR:  usingChanR,
	}
}

type chanReader struct {
	inChan chan string
}

func (c *chanReader) Read(p []byte) (n int, err error) {
	line, ok := <-c.inChan
	if !ok {
		return 0, io.EOF
	}
	return strings.NewReader(line).Read(p)
}

func (d *Debugger) paused() bool {
	return d.stepMode ||
		slices.Contains(d.breakpoints, *d.befunge.InstructionPointer)
}

func (d *Debugger) slowStepping() bool {
	return d.autoSpeed > 0 && !d.jumping
}

func (d *Debugger) Step() (bool, error) {
	// if this is the first step, start a go routine to read from stdin and output to a channel
	// this allows us to slow step through the program and be interrupted by keyboard input
	if !d.isStarted {
		d.isStarted = true
		go func() {
			for !d.isFinished {
				readString, err := d.reader.ReadString('\n')
				if err != nil {
					close(d.stdinChan)
					return
				}
				d.stdinChan <- readString
			}
		}()
	}

	char := d.befunge.CurrentChar()
	if d.nextStepWillPromptForInput(char) {
		d.printDebug(d.awaitingInputControls(char))
		d.hasPrinted = true
	} else if d.paused() {
		d.jumping = false

		d.printDebug(d.interruptedControls())
		d.hasPrinted = true

		str := <-d.stdinChan

		if strings.Contains(str, "c") || strings.Contains(str, "C") {
			d.stepMode = false
		} else if strings.Contains(str, "j") || strings.Contains(str, "J") {
			d.jumping = true
			d.stepMode = false
		} else {
			d.stepMode = true
		}
	} else if d.slowStepping() {
		d.printDebug(d.slowSteppingControls())
		d.hasPrinted = true

		// Create a context with a timeout of d.autoSpeed
		ctx, cancel := context.WithTimeout(context.Background(), d.autoSpeed)
		defer cancel()

		// Wait for either stdin input or the context to timeout
		select {
		case <-d.stdinChan:
			d.stepMode = true
			break
		case <-ctx.Done():
		}
	}

	proceed, err := d.befunge.Step()
	if !proceed {
		d.isFinished = true // stop the stdin go routine
		if d.hasPrinted {
			fmt.Println(clearAndReturn)
		}
		_, err := fmt.Fprint(d.outfile, d.output.String())
		if err != nil {
			panic(err)
		}
	}
	return proceed, err
}

// pretty hacky
func (d *Debugger) nextStepWillPromptForInput(char rune) bool {
	return d.usingChanR &&
		!d.befunge.StringMode &&
		((char == '~' || char == '&') ||
			(char == '/') || (char == '%') && d.befunge.StackPeek() == 0)
}

func (d *Debugger) awaitingInputControls(char rune) string {
	return fmt.Sprintf("[awaiting '%c' input] ", char)
}

func (d *Debugger) slowSteppingControls() string {
	return fmt.Sprintf("[%s to pause] ",
		green.Sprint("return"))
}

func (d *Debugger) interruptedControls() string {
	var jumpString string
	if d.autoSpeed > 0 {
		jumpString = fmt.Sprintf(", %s to jump to next breakpoint", green.Sprint("j"))
	}

	return fmt.Sprintf("[%s to step, %s to continue%s, %s to exit] ",
		green.Sprint("return"),
		green.Sprint("c"),
		jumpString,
		green.Sprint("ctrl+c"))
}

func (d *Debugger) printDebug(action string) {
	fmt.Printf(`%s%s: %d %s: %d %s: '%c'

%s:
%s

%s: [%s]
%s: %s
%s`,
		clearAndReturn,
		bold.Sprint("x"),
		d.befunge.InstructionPointer.X,
		bold.Sprint("y"),
		d.befunge.InstructionPointer.Y,
		bold.Sprint("char"),
		d.befunge.Torus.CharAt(d.befunge.InstructionPointer.X, d.befunge.InstructionPointer.Y),
		bold.Sprint("torus"),
		d.torusToString(),
		bold.Sprint("stack"),
		strings.Join(internal.MapSlice(d.befunge.Stack.Values, func(t int) string {
			var unicodeParen = ""
			if unicode.IsPrint(rune(t)) {
				unicodeParen = fmt.Sprintf(" (%c)", rune(t))
			}
			return fmt.Sprintf("%d%s", t, unicodeParen)
		}), ", "),
		bold.Sprint("output"),
		d.output.String(),
		action,
	)
}

func (d *Debugger) torusToString() string {
	strBuilder := new(strings.Builder)
	torus := d.befunge.Torus
	strBuilder.WriteString(cyan.Sprint("╔"))
	strBuilder.WriteString(cyan.Sprint(strings.Repeat("═", torus.Width)))
	strBuilder.WriteString(cyan.Sprint("╗"))
	strBuilder.WriteString("\n")
	for y, line := range torus.Chars {
		strBuilder.WriteString(cyan.Sprint("║"))
		for x, char := range line {
			currentPointer := *pkg.NewVector2(x, y)
			out := string(char)
			isBreakpoint := slices.Contains(d.breakpoints, currentPointer)
			isCursor := *d.befunge.InstructionPointer == currentPointer
			if isBreakpoint && isCursor {
				if char == ' ' {
					out = redBgAndBoldAndUnderlined.Sprint(out)
				} else {
					out = redAndBoldAndUnderlined.Sprint(out)
				}
			} else if isBreakpoint {
				if char == ' ' {
					out = redBg.Sprint(out)
				} else {
					out = red.Sprint(out)
				}
			} else if isCursor {
				out = boldAndUnderlined.Sprint(out)
			}
			strBuilder.WriteString(out)
		}
		strBuilder.WriteString(cyan.Sprint("║"))
		strBuilder.WriteString(faint.Sprint(y))
		strBuilder.WriteRune('\n')
	}
	strBuilder.WriteString(cyan.Sprint("╚"))
	strBuilder.WriteString(cyan.Sprint(strings.Repeat("═", torus.Width)))
	strBuilder.WriteString(cyan.Sprint("╝"))
	strBuilder.WriteRune('\n')
	var str0s strings.Builder
	var str10s strings.Builder
	var str100s strings.Builder
	for i := range torus.Width {
		mod10 := i % 10
		mod100 := i % 100
		str0s.WriteString(faint.Sprint(mod10))
		if mod100 == 0 && i >= 100 {
			str100s.WriteString(faint.Sprint((i / 100) % 10))
		} else {
			str100s.WriteString(" ")
		}
		if mod10 == 0 && i >= 10 {
			str10s.WriteString(faint.Sprint((i / 10) % 10))
		} else {
			str10s.WriteString(" ")
		}
	}
	if torus.Width > 100 {
		strBuilder.WriteString(" ")
		strBuilder.WriteString(str100s.String())
		strBuilder.WriteRune('\n')
	}
	if torus.Width > 10 {
		strBuilder.WriteString(" ")
		strBuilder.WriteString(str10s.String())
		strBuilder.WriteRune('\n')
	}
	strBuilder.WriteString(" ")
	strBuilder.WriteString(str0s.String())

	return strBuilder.String()

}
