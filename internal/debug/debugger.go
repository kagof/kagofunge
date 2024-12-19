package debug

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"github.com/kagof/gofunge/model"
	"io"
	"os"
	"slices"
	"strings"
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
)

type Debugger struct {
	befunge     *model.Befunge
	breakpoints []model.Vector2
	output      *strings.Builder
	stdout      io.Writer
	stepMode    bool
}

func NewDebugger(s string, stdout io.Writer, stdin io.Reader, breakpoints []model.Vector2) *Debugger {
	b := new(strings.Builder)
	return &Debugger{
		befunge:     model.NewBefunge(s, b, stdin),
		breakpoints: breakpoints,
		output:      b,
		stdout:      stdout,
	}
}

func (d *Debugger) Step() (bool, error) {
	if d.stepMode || slices.Contains(d.breakpoints, *d.befunge.InstructionPointer) {
		d.printDebug2()
		str, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			panic(err)
		}

		if strings.Contains(str, "s") || strings.Contains(str, "S") {
			d.stepMode = true
		} else {
			d.stepMode = false
		}
	}
	proceed, err := d.befunge.Step()
	if !proceed {
		fmt.Println(clearAndReturn)
		_, err := fmt.Fprint(d.stdout, d.output.String())
		if err != nil {
			panic(err)
		}
	}
	return proceed, err
}

func (d *Debugger) printDebug2() {
	fmt.Println(clearAndReturn)
	fmt.Printf(`%s: %d %s: %d %s: '%c'

%s:
%s

%s: %+v
%s: %s
[%s to continue, %s to step, %s to exit]`,
		bold.Sprint("x"),
		d.befunge.InstructionPointer.X,
		bold.Sprint("y"),
		d.befunge.InstructionPointer.Y,
		bold.Sprint("char"),
		d.befunge.Torus.CharAt(d.befunge.InstructionPointer.X, d.befunge.InstructionPointer.Y),
		bold.Sprint("torus"),
		d.torusToString(),
		bold.Sprint("stack"),
		d.befunge.Stack.Values,
		bold.Sprint("output"),
		d.output.String(),
		green.Sprint("return"),
		green.Sprint("s"),
		green.Sprint("ctrl+c"),
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
			currentPointer := *model.NewVector2(x, y)
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
		strBuilder.WriteRune('\n')
	}
	strBuilder.WriteString(cyan.Sprint("╚"))
	strBuilder.WriteString(cyan.Sprint(strings.Repeat("═", torus.Width)))
	strBuilder.WriteString(cyan.Sprint("╝"))
	return strBuilder.String()

}
