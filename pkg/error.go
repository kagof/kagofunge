package pkg

import "fmt"

type BefungeExecutionError struct {
	X, Y int
	Val  rune
	Err  error
}

func (e *BefungeExecutionError) Error() string {
	return fmt.Sprintf("Befunge execution error at position (%d, %d) '%c': %v", e.X, e.Y, e.Val, e.Err.Error())
}

func (e *BefungeExecutionError) Unwrap() error {
	return e.Err
}

func (e *BefungeExecutionError) String() (s string) {
	return e.Error()
}
