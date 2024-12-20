package pkg

import (
	"github.com/kagof/kagofunge/internal"
	"strings"
	"unicode/utf8"
)

type Torus struct {
	Chars  [][]rune
	Width  int
	Height int
}

func (t *Torus) ModWidth(x int) int {
	return ((x % t.Width) + t.Width) % t.Width
}

func (t *Torus) ModHeight(y int) int {
	return ((y % t.Height) + t.Height) % t.Height
}

func (t *Torus) CharAt(x int, y int) rune {
	return t.Chars[t.ModHeight(y)][t.ModWidth(x)]
}

func (t *Torus) SetCharAt(x int, y int, v rune) {
	t.Chars[t.ModHeight(y)][t.ModWidth(x)] = v
}

func NewTorus(s string) *Torus {
	lines := strings.FieldsFunc(strings.ReplaceAll(s, "\r", ""), func(r rune) bool { return r == '\n' })
	var longestLine int = 0
	for _, line := range lines {
		numRunes := utf8.RuneCountInString(line)
		if numRunes > longestLine {
			longestLine = numRunes
		}
	}
	chars := internal.MapSlice(lines, func(line string) []rune {
		return []rune(pad(line, longestLine, ' '))
	})

	return &Torus{Chars: chars, Width: longestLine, Height: len(lines)}
}

func pad(s string, size int, padVal rune) string {
	if utf8.RuneCountInString(s) < size {
		return s + strings.Repeat(string(padVal), size-utf8.RuneCountInString(s))
	}
	return s
}
