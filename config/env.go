package config

import (
	"os"
	"strconv"
)

func overridePropsFromEnv(p *Config) error {
	div0, err := fromEnvOrDefault("KGF_INTERPRETER_DIVIDE_BY_ZERO_BEHAVIOUR",
		div0Mapper,
		p.Interpreter.DivideByZeroBehaviour)
	if err != nil {
		return err
	}
	p.Interpreter.DivideByZeroBehaviour = div0
	mod0, err := fromEnvOrDefault("KGF_INTERPRETER_MODULUS_BY_ZERO_BEHAVIOUR",
		div0Mapper,
		p.Interpreter.ModulusByZeroBehaviour)
	if err != nil {
		return err
	}
	p.Interpreter.ModulusByZeroBehaviour = mod0
	poob, err := fromEnvOrDefault("KGF_INTERPRETER_PUT_OUT_OF_BOUNDS_BEHAVIOUR",
		oobMapper,
		p.Interpreter.PutOutOfBoundsBehaviour)
	if err != nil {
		return err
	}
	p.Interpreter.PutOutOfBoundsBehaviour = poob
	goob, err := fromEnvOrDefault("KGF_INTERPRETER_GET_OUT_OF_BOUNDS_BEHAVIOUR",
		oobMapper,
		p.Interpreter.GetOutOfBoundsBehaviour)
	if err != nil {
		return err
	}
	p.Interpreter.GetOutOfBoundsBehaviour = goob
	tSize, err := fromEnvOrDefault("KGF_INTERPRETER_ENFORCE_TORUS_SIZE_RESTRICTION",
		strconv.ParseBool,
		p.Interpreter.EnforceTorusSizeRestriction)
	if err != nil {
		return err
	}
	p.Interpreter.EnforceTorusSizeRestriction = tSize

	showT, err := fromEnvOrDefault("KGF_DEBUGGER_SHOW_TORUS",
		strconv.ParseBool,
		p.Debugger.ShowTorus)
	if err != nil {
		return err
	}
	p.Debugger.ShowTorus = showT
	showS, err := fromEnvOrDefault("KGF_DEBUGGER_SHOW_STACK",
		strconv.ParseBool,
		p.Debugger.ShowStack)
	if err != nil {
		return err
	}
	p.Debugger.ShowStack = showS
	colors, err := fromEnvOrDefault("KGF_DEBUGGER_ENABLE_COLORS",
		strconv.ParseBool,
		p.Debugger.EnableColors)
	if err != nil {
		return err
	}
	p.Debugger.EnableColors = colors

	return nil
}

func fromEnvOrDefault[T any](envVarName string, mapper func(string) (T, error), defaultVal T) (T, error) {
	fromEnv, exists := os.LookupEnv(envVarName)
	if !exists {
		return defaultVal, nil
	}
	return mapper(fromEnv)
}
