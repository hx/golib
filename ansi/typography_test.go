package ansi_test

import (
	"fmt"
	. "github.com/hx/golib/ansi"
	. "github.com/hx/golib/testing"
	"strings"
	"testing"
)

func TestEscapeSequence_String(t *testing.T) {
	type thing struct {
		actual   EscapeSequence
		expected string
	}
	clean := func(s string) string { return strings.Replace(s, "\033", "\\033", -1) }
	things := []thing{
		{Yellow | BlueBG | Bold, "\033[1;33;44m"},
		{Yellow | BlueBG, "\033[33;44m"},
		{Yellow | Bold, "\033[1;33m"},
		{Yellow, "\033[33m"},
		{Yellow | Bright, "\033[93m"},
	}
	for i := range things {
		var (
			thing    = things[i]
			expected = thing.expected
			actual   = thing.actual.String()
		)
		t.Run(fmt.Sprintf("Thing %d: %s", i, clean(expected)), func(t *testing.T) {
			Assert(t, expected == actual,
				fmt.Sprintf("Expected '%s' but got '%s'", clean(expected), clean(actual)))
		})
	}
}
