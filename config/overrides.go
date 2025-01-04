package config

import "strconv"

func overridePropsFromMap(p *Config, overrides map[string]string) error {
	div0, err := fromMapOrDefault(overrides,
		"interpreter.divide-by-zero-behaviour",
		div0Mapper,
		p.Interpreter.DivideByZeroBehaviour)
	if err != nil {
		return err
	}
	p.Interpreter.DivideByZeroBehaviour = div0
	mod0, err := fromMapOrDefault(overrides,
		"interpreter.modulus-by-zero-behaviour",
		div0Mapper,
		p.Interpreter.ModulusByZeroBehaviour)
	if err != nil {
		return err
	}
	p.Interpreter.ModulusByZeroBehaviour = mod0
	poob, err := fromMapOrDefault(overrides,
		"interpreter.put-out-of-bounds-behaviour",
		oobMapper,
		p.Interpreter.PutOutOfBoundsBehaviour)
	if err != nil {
		return err
	}
	p.Interpreter.PutOutOfBoundsBehaviour = poob
	goob, err := fromMapOrDefault(overrides,
		"interpreter.get-out-of-bounds-behaviour",
		oobMapper,
		p.Interpreter.GetOutOfBoundsBehaviour)
	if err != nil {
		return err
	}
	p.Interpreter.GetOutOfBoundsBehaviour = goob
	tSize, err := fromMapOrDefault(overrides,
		"interpreter.enforce-torus-size-restriction",
		strconv.ParseBool,
		p.Interpreter.EnforceTorusSizeRestriction)
	if err != nil {
		return err
	}
	p.Interpreter.EnforceTorusSizeRestriction = tSize

	showT, err := fromMapOrDefault(overrides,
		"debugger.show-torus",
		strconv.ParseBool,
		p.Debugger.ShowTorus)
	if err != nil {
		return err
	}
	p.Debugger.ShowTorus = showT
	showS, err := fromMapOrDefault(overrides,
		"debugger.show-stack",
		strconv.ParseBool,
		p.Debugger.ShowStack)
	if err != nil {
		return err
	}
	p.Debugger.ShowStack = showS
	colors, err := fromMapOrDefault(overrides,
		"debugger.show-enable-colors",
		strconv.ParseBool,
		p.Debugger.EnableColors)
	if err != nil {
		return err
	}
	p.Debugger.EnableColors = colors

	return nil

}

func fromMapOrDefault[T any](overrides map[string]string, key string, mapper func(string) (T, error), defaultValue T) (T, error) {
	fromMap := overrides[key]
	if fromMap == "" {
		return defaultValue, nil
	}
	return mapper(fromMap)
}
