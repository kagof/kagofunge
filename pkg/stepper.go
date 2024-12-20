package pkg

type Stepper interface {
	// Step Process the next instruction. Returns false if the program is terminated, true if there is a next step.
	Step() (bool, error)
}
