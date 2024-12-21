package pkg

import (
	"bufio"
	"io"
)

type Befunge struct {
	writer             io.Writer
	reader             *bufio.Reader
	Stack              *Stack[int]
	Torus              *Torus
	InstructionPointer *Vector2
	stringMode         bool
	delta              *Vector2
	halted             bool
}

func NewBefunge(s string, w io.Writer, r io.Reader) *Befunge {
	torus := NewTorus(s)
	return &Befunge{
		writer:             w,
		reader:             bufio.NewReader(r),
		Stack:              NewStack[int](),
		Torus:              torus,
		InstructionPointer: NewVector2(0, 0),
		stringMode:         false,
		delta:              XPos(),
		halted:             false,
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

func (f *Befunge) Step() (bool, error) {
	var err error
	char := f.Torus.CharAt(f.InstructionPointer.X, f.InstructionPointer.Y)
	if f.stringMode && char != '"' {
		f.Stack.Push(int(char))
	} else {
		instruction := ParseInstruction(char)
		err = instruction.PerformInstruction(f)
	}
	if f.halted {
		return false, mapErr(f, err)
	}
	f.step()
	return true, nil
}

func (f *Befunge) step() {
	if f.delta.X != 0 { //saving some modulus operations
		f.InstructionPointer.X = f.Torus.ModWidth(f.InstructionPointer.X + f.delta.X)
	}
	if f.delta.Y != 0 { //saving some modulus operations
		f.InstructionPointer.Y = f.Torus.ModHeight(f.InstructionPointer.Y + f.delta.Y)
	}
	f.delta = f.delta.ScaleToOne()
}

func mapErr(funge *Befunge, err error) error {
	if err == nil {
		return nil
	}
	return &BefungeExecutionError{
		X:   funge.InstructionPointer.X,
		Y:   funge.InstructionPointer.Y,
		Val: funge.Torus.CharAt(funge.InstructionPointer.X, funge.InstructionPointer.Y),
		Err: err,
	}
}
