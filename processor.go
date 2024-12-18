package main

type BefungeProcessor interface {
	// BefungeProcess Process the next instruction. Returns false if the program is terminated.
	BefungeProcess() (bool, error)
}
