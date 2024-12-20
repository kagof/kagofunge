package pkg

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

type Vector2 struct {
	X, Y int
}

func NewVector2(x int, y int) *Vector2 {
	return &Vector2{X: x, Y: y}
}

func (v *Vector2) String() string {
	return fmt.Sprintf("(%d,%d)", v.X, v.Y)
}

func (v *Vector2) Add(other *Vector2) *Vector2 {
	return NewVector2(v.X+other.X, v.Y+other.Y)
}

// ScaleToOne converts eg (-2, 4) to (-1, 1), and (3, 0) to (1, 0). This does not necessarily
// convert to a unit vector of magnitude 1, rather converts each dimension to 1 if it is positive, -1 if it is negative,
// or 0 if it is 0
func (v *Vector2) ScaleToOne() *Vector2 {
	absX := abs(v.X)
	absY := abs(v.Y)
	if (absX == 1 || absX == 0) && (absY == 1 || absY == 0) {
		return v
	}
	return NewVector2(v.X/ifZero(absX, 1), v.Y/ifZero(absY, 1))
}

func ifZero(first int, second int) int {
	if first == 0 {
		return second
	}
	return first
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func XPos() *Vector2 {
	return NewVector2(1, 0)
}

func XNeg() *Vector2 {
	return NewVector2(-1, 0)
}

func YPos() *Vector2 {
	return NewVector2(0, 1)
}

func YNeg() *Vector2 {
	return NewVector2(0, -1)
}

func ParseVector2(str string) (*Vector2, error) {
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

func getVals(fmtStr string, str string) (*Vector2, error) {
	var x, y int
	_, err := fmt.Sscanf(str, fmtStr, &x, &y)
	if err != nil {
		return nil, err
	}
	return NewVector2(x, y), err
}
