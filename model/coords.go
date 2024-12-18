package model

import (
	"fmt"
	"regexp"
)

var (
	regexParenComma, _   = regexp.Compile(`\(\d+,\d+\)`)
	regexParenCommaSp, _ = regexp.Compile(`\(\d+, \d+\)`)
	regexParenSp, _      = regexp.Compile(`\(\d+ \d+\)`)
	regexBracComma, _    = regexp.Compile(`\[\d+,\d+]`)
	regexBracCommaSp, _  = regexp.Compile(`\[\d+, \d+]`)
	regexBracSp, _       = regexp.Compile(`\[\d+ \d+]`)
	regexComma, _        = regexp.Compile(`\d+,\d+`)
)

type Coords struct {
	X, Y int
}

func NewCoords(x int, y int) *Coords {
	return &Coords{X: x, Y: y}
}

func (c *Coords) String() string {
	return fmt.Sprintf("(%d,%d)", c.X, c.Y)
}

func ParseCoords(str string) (*Coords, error) {
	if regexParenComma.MatchString(str) {
		return getVals("(%d,%d)", str)
	} else if regexParenCommaSp.MatchString(str) {
		return getVals("(%d, %d)", str)
	} else if regexParenSp.MatchString(str) {
		return getVals("(%d %d)", str)
	} else if regexBracComma.MatchString(str) {
		return getVals("[%d,%d]", str)
	} else if regexBracCommaSp.MatchString(str) {
		return getVals("[%d, %d]", str)
	} else if regexBracSp.MatchString(str) {
		return getVals("[%d %d]", str)
	} else if regexComma.MatchString(str) {
		return getVals("%d,%d", str)
	}
	return nil, fmt.Errorf("invalid pointer value: %s", str)
}

func getVals(fmtStr string, str string) (*Coords, error) {
	var x, y int
	_, err := fmt.Sscanf(str, fmtStr, &x, &y)
	if err != nil {
		return nil, err
	}
	return NewCoords(x, y), err
}
