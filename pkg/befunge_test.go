package pkg

import (
	"github.com/kagof/kagofunge/config"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

const maxSteps = 10000

func TestBefunge(t *testing.T) {
	t.Parallel()
	asserts := assert.New(t)

	var cases = []struct {
		name     string
		funge    string
		input    string
		expected string
	}{
		{
			"hello_world",
			` >               v
 v"Hello, World!"<
 >:v
 ^,_@`,
			"",
			"Hello, World!",
		},
		{
			"cat",
			"~:!#@_,",
			"testing",
			"testing",
		},
		{
			"factorial",
			`&>:1-:v v *_$.@ 
 ^    _$>\:^`,
			"5",
			"120",
		},
		{
			"quine",
			"01->1# +# :# 0# g# ,# :# 5# 8# *# 4# +# -# _@",
			"",
			"01->1# +# :# 0# g# ,# :# 5# 8# *# 4# +# -# _@",
		},
		{
			"divide_by_zero",
			"10/.@",
			"3",
			"3", // dividing by zero should ask the user what it should equal
		},
		{
			"g_oob",
			"1-g.@",
			"",
			"0", // getting from out of bounds should not wrap, but should return 0
		},
		{
			"p_oob",
			"\" \"02-0p#@ #.<",
			"",
			"0", // putting out of bounds should not wrap, but should just discard
		},
	}

	cfg := config.DefaultConfig()
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var writer strings.Builder
			befunge := NewBefunge(&cfg, test.funge, &writer, strings.NewReader(test.input))

			var hasNext = true
			var err error
			var i = 0
			for hasNext && i < maxSteps {
				hasNext, err = befunge.Step()
				asserts.NoError(err, "no error expected while executing %s", test.name)
				i += 1
			}
			asserts.Less(i, maxSteps, "exceeded %d steps executing %s", maxSteps, test.name)
			asserts.Equal(test.expected, writer.String(), "%s output not as expected", test.name)
		})
	}
}
