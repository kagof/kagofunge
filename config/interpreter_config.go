package config

type InterpreterConfig struct {
	DivideByZeroBehaviour       DivideByZeroBehaviour `yaml:"divide-by-zero-behaviour"`
	ModulusByZeroBehaviour      DivideByZeroBehaviour `yaml:"modulus-by-zero-behaviour"`
	PutOutOfBoundsBehaviour     OutOfBoundsBehaviour  `yaml:"put-out-of-bounds-behaviour"`
	GetOutOfBoundsBehaviour     OutOfBoundsBehaviour  `yaml:"get-out-of-bounds-behaviour"`
	EnforceTorusSizeRestriction bool                  `yaml:"enforce-torus-size-restriction"`
	TorusSizeRestrictionWidth   int                   `yaml:"torus-size-restriction-width"`
	TorusSizeRestrictionHeight  int                   `yaml:"torus-size-restriction-height"`
}

type DivideByZeroBehaviour string

const (
	Div0PromptForInput DivideByZeroBehaviour = "PROMPT_FOR_INPUT"
	Div0ReturnZero     DivideByZeroBehaviour = "RETURN_ZERO"
	Div0Reflect        DivideByZeroBehaviour = "REFLECT"
	Div0Panic          DivideByZeroBehaviour = "PANIC"
)

var divideByZeroBehaviours = map[string]DivideByZeroBehaviour{
	"PROMPT_FOR_INPUT": Div0PromptForInput,
	"RETURN_ZERO":      Div0ReturnZero,
	"REFLECT":          Div0Reflect,
	"PANIC":            Div0Panic,
}

type OutOfBoundsBehaviour string

const (
	OobNoOp  OutOfBoundsBehaviour = "NO_OP"
	OobZero  OutOfBoundsBehaviour = "ZERO"
	OobWrap  OutOfBoundsBehaviour = "WRAP"
	OobPanic OutOfBoundsBehaviour = "PANIC"
)

var outOfBoundsBehaviours = map[string]OutOfBoundsBehaviour{
	"NO_OP": OobNoOp,
	"ZERO":  OobZero,
	"WRAP":  OobWrap,
	"PANIC": OobPanic,
}
