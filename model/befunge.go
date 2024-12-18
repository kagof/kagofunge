package model

import (
	"io"
)

type Direction uint8

const (
	right Direction = iota
	left  Direction = iota
	down  Direction = iota
	up    Direction = iota
)

type Befunge struct {
	Stack          *Stack[int]
	stringMode     bool
	direction      Direction
	Torus          *Torus
	ProgramCounter *Coords
	writer         io.Writer
	reader         io.Reader
}

func NewBefunge(s string, w io.Writer, r io.Reader) *Befunge {
	torus := NewTorus(s)
	return &Befunge{
		Stack:          NewStack[int](),
		stringMode:     false,
		direction:      right,
		Torus:          torus,
		ProgramCounter: NewCoords(0, 0),
		writer:         w,
		reader:         r,
	}
}

func (f *Befunge) stackPop() int {
	v, b := f.Stack.Pop()
	if b {
		return v
	}
	return 0
}

func (f *Befunge) stackPeek() int {
	v, b := f.Stack.Peek()
	if b {
		return v
	}
	return 0
}

func (f *Befunge) BefungeProcess() (bool, error) {
	var err error
	distToNext := 1
	stopped := false
	char := f.Torus.CharAt(f.ProgramCounter.X, f.ProgramCounter.Y)
	if f.stringMode && char != '"' {
		f.Stack.Push(int(char))
	} else {
		instruction := ParseInstruction(char)
		stopped, distToNext, err = instruction.PerformInstruction(f)
	}
	if stopped {
		return false, mapErr(f, err)
	}
	f.step(distToNext)
	return true, nil
}

func (f *Befunge) step(distToNext int) {
	switch f.direction {
	case right:
		f.ProgramCounter.X = f.Torus.ModWidth(f.ProgramCounter.X + distToNext)
	case left:
		f.ProgramCounter.X = f.Torus.ModWidth(f.ProgramCounter.X - distToNext)
	case down:
		f.ProgramCounter.Y = f.Torus.ModHeight(f.ProgramCounter.Y + distToNext)
	case up:
		f.ProgramCounter.Y = f.Torus.ModHeight(f.ProgramCounter.Y - distToNext)
	}
}

func mapErr(funge *Befunge, err error) error {
	if err == nil {
		return nil
	}
	return &BefungeExecutionError{
		X:   funge.ProgramCounter.X,
		Y:   funge.ProgramCounter.Y,
		Val: funge.Torus.CharAt(funge.ProgramCounter.X, funge.ProgramCounter.Y),
		Err: err,
	}
}
