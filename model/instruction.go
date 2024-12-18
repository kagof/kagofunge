package model

import (
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
)

type InstructionPerformer interface {
	PerformInstruction(f *Befunge) (bool, int, error)
}

type stringMode struct {
}

func (s stringMode) PerformInstruction(f *Befunge) (bool, int, error) {
	f.stringMode = !f.stringMode
	return false, 1, nil
}

type num struct {
	val int
}

func (n num) PerformInstruction(f *Befunge) (bool, int, error) {
	f.Stack.Push(n.val)
	return false, 1, nil
}

type op2 struct {
	operator func(int, int) (int, error)
}

func (o op2) PerformInstruction(f *Befunge) (bool, int, error) {
	a, b := f.stackPop(), f.stackPop()
	res, err := o.operator(a, b)
	if err != nil {
		return true, 1, err
	}
	f.Stack.Push(res)
	return false, 1, nil
}

type not struct {
}

func (n not) PerformInstruction(f *Befunge) (bool, int, error) {
	if f.stackPop() == 0 {
		f.Stack.Push(1)
	} else {
		f.Stack.Push(0)
	}
	return false, 1, nil
}

type dir struct {
	direction func() Direction
}

func (d dir) PerformInstruction(f *Befunge) (bool, int, error) {
	f.direction = d.direction()
	return false, 1, nil
}

type conditionalDir struct {
	zeroDir Direction
	elseDir Direction
}

func (d conditionalDir) PerformInstruction(f *Befunge) (bool, int, error) {
	if f.stackPop() == 0 {
		f.direction = d.zeroDir
	} else {
		f.direction = d.elseDir
	}
	return false, 1, nil
}

type dup struct {
}

func (d dup) PerformInstruction(f *Befunge) (bool, int, error) {
	f.Stack.Push(f.stackPeek())
	return false, 1, nil
}

type swap struct {
}

func (s swap) PerformInstruction(f *Befunge) (bool, int, error) {
	a, b := f.stackPop(), f.stackPop()
	f.Stack.Push(a)
	f.Stack.Push(b)
	return false, 1, nil
}

type discard struct {
}

func (d discard) PerformInstruction(f *Befunge) (bool, int, error) {
	f.stackPop()
	return false, 1, nil
}

type write struct {
	writeFun func(io.Writer, int) (int, error)
}

func (w write) PerformInstruction(f *Befunge) (bool, int, error) {
	_, err := w.writeFun(f.writer, f.stackPop())
	if err != nil {
		return true, 1, err
	}
	return false, 1, nil
}

type skip struct {
}

func (s skip) PerformInstruction(f *Befunge) (bool, int, error) {
	return false, 2, nil
}

type put struct {
}

func (p put) PerformInstruction(f *Befunge) (bool, int, error) {
	y, x, v := f.stackPop(), f.stackPop(), rune(f.stackPop())
	f.Torus.SetCharAt(x, y, v)
	return false, 1, nil
}

type get struct {
}

func (g get) PerformInstruction(f *Befunge) (bool, int, error) {
	y, x := f.stackPop(), f.stackPop()
	f.Stack.Push(int(f.Torus.CharAt(x, y)))
	return false, 1, nil
}

type read struct {
	readFun func(io.Reader) (int, int, error)
}

func (r read) PerformInstruction(f *Befunge) (bool, int, error) {
	i, _, err := r.readFun(f.reader)
	if err == nil {
		f.Stack.Push(i)
		return false, 1, nil
	} else {
		return true, 1, err
	}
}

type terminate struct {
}

func (t terminate) PerformInstruction(f *Befunge) (bool, int, error) {
	return true, 1, nil
}

type noop struct {
}

func (n noop) PerformInstruction(f *Befunge) (bool, int, error) {
	return false, 1, nil
}

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
		return dir{direction: func() Direction {
			return right
		}}
	case char == '<':
		return dir{direction: func() Direction {
			return left
		}}
	case char == 'v':
		return dir{direction: func() Direction {
			return down
		}}
	case char == '^':
		return dir{direction: func() Direction {
			return up
		}}
	case char == '?':
		return dir{direction: func() Direction {
			return Direction(rand.IntN(4))
		}}
	case char == '_':
		return conditionalDir{
			zeroDir: right,
			elseDir: left,
		}
	case char == '|':
		return conditionalDir{
			zeroDir: down,
			elseDir: up,
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
		return write{writeFun: func(w io.Writer, a int) (int, error) {
			return fmt.Fprintf(w, "%d ", a)
		}}
	case char == ',':
		return write{writeFun: func(w io.Writer, a int) (int, error) {
			return fmt.Fprint(w, string(rune(a)))
		}}
	case char == '#':
		return skip{}
	case char == 'p':
		return put{}
	case char == 'g':
		return get{}
	case char == '&':
		return read{readFun: func(r io.Reader) (int, int, error) {
			var i int
			n, err := fmt.Fscanf(r, "%d\n", &i)
			return i, n, err
		}}
	case char == '~':
		return read{readFun: func(r io.Reader) (int, int, error) {
			var ch rune
			n, err := fmt.Fscanf(r, "%c\n", &ch)
			return int(ch), n, err
		}}
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
