package ansi_test

import (
	"fmt"
	. "github.com/hx/golib/ansi"
	. "github.com/hx/golib/testing"
	"testing"
)

func TestTruncate(t *testing.T) {
	type exp struct {
		original  string
		maxLength int
		expected  string
	}
	var (
		red   = Red.String()
		green = Green.String()
		reset = Reset

		fooBarBaz = "foo" + red + "bar" + green + "baz" + reset
		bazBarFoo = reset + green + "baz" + red + "bar" + green + "foo"
		nonsense  = red + "ab…" + green + "®🤩-" + reset
	)
	cases := []exp{
		{fooBarBaz, 9, fooBarBaz},
		{fooBarBaz, 8, "foo" + red + "bar" + green + "ba" + reset},
		{fooBarBaz, 7, "foo" + red + "bar" + green + "b" + reset},
		{fooBarBaz, 6, "foo" + red + "bar" + green + reset},
		{fooBarBaz, 5, "foo" + red + "ba" + green + reset},
		{fooBarBaz, 4, "foo" + red + "b" + green + reset},
		{fooBarBaz, 3, "foo" + red + green + reset},
		{fooBarBaz, 2, "fo" + red + green + reset},
		{fooBarBaz, 1, "f" + red + green + reset},
		{fooBarBaz, 0, ""},

		{bazBarFoo, 9, bazBarFoo},
		{bazBarFoo, 8, reset + green + "baz" + red + "bar" + green + "fo"},
		{bazBarFoo, 7, reset + green + "baz" + red + "bar" + green + "f"},
		{bazBarFoo, 6, reset + green + "baz" + red + "bar" + green},
		{bazBarFoo, 5, reset + green + "baz" + red + "ba" + green},
		{bazBarFoo, 4, reset + green + "baz" + red + "b" + green},
		{bazBarFoo, 3, reset + green + "baz" + red + green},
		{bazBarFoo, 2, reset + green + "ba" + red + green},
		{bazBarFoo, 1, reset + green + "b" + red + green},
		{bazBarFoo, 0, ""},

		{nonsense, 6, nonsense},
		{nonsense, 5, red + "ab…" + green + "®🤩" + reset},
		{nonsense, 4, red + "ab…" + green + "®" + reset},
		{nonsense, 3, red + "ab…" + green + reset},
		{nonsense, 2, red + "ab" + green + reset},
		{nonsense, 1, red + "a" + green + reset},
		{nonsense, 0, ""},
	}
	for i := range cases {
		c := cases[i]
		t.Run(fmt.Sprintf("(%s,%d) => %s", c.original, c.maxLength, c.expected), func(t *testing.T) {
			Equals(t, c.expected, Truncate(c.original, c.maxLength))
		})
	}
}
