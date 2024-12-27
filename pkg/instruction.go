package pkg

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"strconv"
	"strings"
)

type InstructionPerformer interface {
	PerformInstruction(f *Befunge) error
}

type stringMode struct {
}

func (s stringMode) PerformInstruction(f *Befunge) error {
	f.StringMode = !f.StringMode
	return nil
}

type num struct {
	val int
}

func (n num) PerformInstruction(f *Befunge) error {
	f.Stack.Push(n.val)
	return nil
}

type op2 struct {
	operator func(int, int) (int, error)
}

func (o op2) PerformInstruction(f *Befunge) error {
	a, b := f.stackPop(), f.stackPop()
	res, err := o.operator(a, b)
	if err != nil {
		if strings.Contains(err.Error(), "divide by zero") {
			// the befunge spec states that division/modulus by 0 should result in asking the user for desired outcome
			return intRead.PerformInstruction(f)
		} else {
			return err
		}
	}
	f.Stack.Push(res)
	return nil
}

type not struct {
}

func (n not) PerformInstruction(f *Befunge) error {
	if f.stackPop() == 0 {
		f.Stack.Push(1)
	} else {
		f.Stack.Push(0)
	}
	return nil
}

type dir struct {
	delta func() *Vector2
}

func (d dir) PerformInstruction(f *Befunge) error {
	f.delta = d.delta()
	return nil
}

type conditionalDir struct {
	zeroDir func() *Vector2
	elseDir func() *Vector2
}

func (d conditionalDir) PerformInstruction(f *Befunge) error {
	if f.stackPop() == 0 {
		f.delta = d.zeroDir()
	} else {
		f.delta = d.elseDir()
	}
	return nil
}

type dup struct {
}

func (d dup) PerformInstruction(f *Befunge) error {
	f.Stack.Push(f.StackPeek())
	return nil
}

type swap struct {
}

func (s swap) PerformInstruction(f *Befunge) error {
	a, b := f.stackPop(), f.stackPop()
	f.Stack.Push(a)
	f.Stack.Push(b)
	return nil
}

type discard struct {
}

func (d discard) PerformInstruction(f *Befunge) error {
	f.stackPop()
	return nil
}

type write struct {
	writeFun func(io.Writer, int) (int, error)
}

func (w write) PerformInstruction(f *Befunge) error {
	_, err := w.writeFun(f.writer, f.stackPop())
	if err != nil {
		f.halted = true
		return err
	}
	return nil
}

type skip struct {
}

func (s skip) PerformInstruction(f *Befunge) error {
	f.delta = f.delta.Add(f.delta)
	return nil
}

type put struct {
}

func (p put) PerformInstruction(f *Befunge) error {
	y, x, v := f.stackPop(), f.stackPop(), rune(f.stackPop())
	if y >= f.Torus.Height || y < 0 || x >= f.Torus.Width || x < 0 {
		return nil
	}
	f.Torus.SetCharAt(x, y, v)
	return nil
}

type get struct {
}

func (g get) PerformInstruction(f *Befunge) error {
	y, x := f.stackPop(), f.stackPop()
	if y >= f.Torus.Height || y < 0 || x >= f.Torus.Width || x < 0 {
		return nil
	}
	f.Stack.Push(int(f.Torus.CharAt(x, y)))
	return nil
}

type read struct {
	readFun func(*bufio.Reader) (int, error)
}

func (r read) PerformInstruction(f *Befunge) error {
	i, err := r.readFun(f.reader)
	if err == nil {
		f.Stack.Push(i)
		return nil
	} else {
		f.halted = true
		return err
	}
}

type terminate struct {
}

func (t terminate) PerformInstruction(f *Befunge) error {
	f.halted = true
	return nil
}

type noop struct {
}

func (n noop) PerformInstruction(f *Befunge) error {
	return nil
}

var intWrite = write{writeFun: func(w io.Writer, a int) (int, error) {
	return fmt.Fprintf(w, "%d", a)
}}

var charWrite = write{writeFun: func(w io.Writer, a int) (int, error) {
	return fmt.Fprint(w, string(rune(a)))
}}

var intRead = read{readFun: func(r *bufio.Reader) (int, error) {
	for {
		b, _, err1 := r.ReadLine()
		// error reading from the Reader, return
		if err1 != nil {
			if err1 == io.EOF {
				// Don't error for EOF, just return 0
				return 0, nil
			}
			return 0, err1
		}

		// Try to convert the input to an integer. If it fails, re-prompt for input
		num, err2 := strconv.Atoi(strings.TrimSpace(string(b)))
		if err2 == nil {
			// int parsed
			return num, nil
		}
	}
}}

var charRead = read{readFun: func(r *bufio.Reader) (int, error) {
	for {
		ch, _, err := r.ReadRune()
		if err != nil {
			if err == io.EOF {
				// Don't error for EOF, just return 0
				return 0, nil
			}
			return 0, err
		}

		if ch != '\n' && ch != '\r' {
			return int(ch), nil
		} // else ignore empty lines
	}
}}

func ParseInstruction(char rune) InstructionPerformer {
	switch {
	case '0' <= char && char <= '9':
		return num{val: int(char - '0')}
	case char == '+':
		return op2{operator: func(a int, b int) (int, error) {
			return a + b, nil
		}}
	case char == '-':
		return op2{operator: func(a int, b int) (int, error) {
			return b - a, nil
		}}
	case char == '*':
		return op2{operator: func(a int, b int) (int, error) {
			return a * b, nil
		}}
	case char == '/':
		return op2{operator: func(a int, b int) (int, error) {
			if a == 0 {
				return 0, errors.New("divide by zero")
			}
			return b / a, nil
		}}
	case char == '%':
		return op2{operator: func(a int, b int) (int, error) {
			if a == 0 {
				return 0, errors.New("divide by zero")
			}
			return b % a, nil
		}}
	case char == '!':
		return not{}
	case char == '`':
		return op2{operator: func(a int, b int) (int, error) {
			return boolToInt(b > a), nil
		}}
	case char == '>':
		return dir{delta: func() *Vector2 {
			return XPos()
		}}
	case char == '<':
		return dir{delta: func() *Vector2 {
			return XNeg()
		}}
	case char == 'v':
		return dir{delta: func() *Vector2 {
			return YPos()
		}}
	case char == '^':
		return dir{delta: func() *Vector2 {
			return YNeg()
		}}
	case char == '?':
		return dir{delta: func() *Vector2 {
			switch rand.IntN(4) {
			case 0:
				return XPos()
			case 1:
				return XNeg()
			case 2:
				return YPos()
			case 3:
				fallthrough
			default:
				return YNeg()
			}
		}}
	case char == '_':
		return conditionalDir{
			zeroDir: XPos,
			elseDir: XNeg,
		}
	case char == '|':
		return conditionalDir{
			zeroDir: YPos,
			elseDir: YNeg,
		}
	case char == '"':
		return stringMode{}
	case char == ':':
		return dup{}
	case char == '\\':
		return swap{}
	case char == '$':
		return discard{}
	case char == '.':
		return intWrite
	case char == ',':
		return charWrite
	case char == '#':
		return skip{}
	case char == 'p':
		return put{}
	case char == 'g':
		return get{}
	case char == '&':
		return intRead
	case char == '~':
		return charRead
	case char == '@':
		return terminate{}
	default:
		return noop{}
	}
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
