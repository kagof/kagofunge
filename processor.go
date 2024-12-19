package main

type Stepper interface {
	// Step Process the next instruction. Returns false if the program is terminated.
	Step() (bool, error)
}
