package pkg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVector2_ScaleToOne(t *testing.T) {
	t.Parallel()
	asserts := assert.New(t)

	var cases = []struct {
		inX, inY, expectedX, expectedY int
	}{
		{0, 0, 0, 0},
		{0, 1, 0, 1},
		{1, 1, 1, 1},
		{-1, 1, -1, 1},
		{-5, 10, -1, 1},
		{10, 10, 1, 1},
		{10, 0, 1, 0},
		{0, 10, 0, 1},
		{-10, 0, -1, 0},
		{0, -10, 0, -1},
	}

	for _, test := range cases {
		v := Vector2{test.inX, test.inY}
		t.Run(v.String(), func(t *testing.T) {
			t.Parallel()
			result := v.ScaleToOne()
			asserts.Equal(test.expectedX, result.X, "X value mismatch")
			asserts.Equal(test.expectedY, result.Y, "Y value mismatch")
		})
	}
}
