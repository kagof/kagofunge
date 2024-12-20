package pkg

import "fmt"

type BefungeExecutionError struct {
	X, Y int
	Val  rune
	Err  error
}

func (e *BefungeExecutionError) Error() string {
	return fmt.Sprintf("Befunge execution error at (%d, %d) with val %c: %v", e.X, e.Y, e.Val, e.Err)
}
