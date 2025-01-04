package pkg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTorus_parse(t *testing.T) {
	t.Parallel()
	asserts := assert.New(t)

	var cases = []struct {
		name                 string
		str                  string
		expectedWidth        int
		expectedHeight       int
		expectedCharAtOrigin rune
		expectedCharAtEnd    rune
	}{
		{
			name:                 "trivial",
			str:                  "@",
			expectedWidth:        1,
			expectedHeight:       1,
			expectedCharAtOrigin: '@',
			expectedCharAtEnd:    '@',
		},
		{
			name: "multiline_different_lengths",
			str: `v     v <
 @
 ^      <
>      #^          <`,
			expectedWidth:        20,
			expectedHeight:       4,
			expectedCharAtOrigin: 'v',
			expectedCharAtEnd:    '<',
		},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			torus := NewTorus(test.str, -1, -1)
			asserts.Equal(test.expectedWidth, torus.Width, "Width mismatch")
			asserts.Equal(test.expectedHeight, torus.Height, "Height mismatch")
			asserts.Equal(test.expectedCharAtOrigin, torus.CharAt(0, 0), "CharAt(0, 0) mismatch")
			asserts.Equal(test.expectedCharAtEnd, torus.CharAt(torus.Width-1, torus.Height-1),
				"CharAt(torus.Width-1, torus.Height-1) mismatch")
		})
	}
}
