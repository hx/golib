// Package testing provides simple assertion tools for writing unit tests.
package testing

import (
	"fmt"
	"reflect"
	t "testing"
)

// Assert fails the test if the condition is false.
func Assert(tb t.TB, condition bool, msg string, v ...interface{}) {
	tb.Helper()
	if !condition {
		tb.Logf(msg, v...)
		tb.FailNow()
	}
}

// Ok fails the test if an err is not nil.
func Ok(tb t.TB, err error) {
	tb.Helper()
	Assert(tb, err == nil, fmt.Sprintf("%v", err))
}

// Equals fails the test if exp is not equal to act.
func Equals(tb t.TB, exp, act interface{}) {
	tb.Helper()
	Assert(tb, reflect.DeepEqual(exp, act), "Expected %v, but got %v", exp, act)
}

// NotEquals fails the test if exp is equal to act.
func NotEquals(tb t.TB, exp, act interface{}) {
	tb.Helper()
	Assert(tb, !reflect.DeepEqual(exp, act), "Expected anything but %v", exp)
}

// Panic fails the test if function does not panic. Returns the panic.
func Panic(tb t.TB, fn func()) (err interface{}) {
	defer func() {
		err = recover()
		Assert(tb, err != nil, "Expected a panic")
	}()
	fn()
	return
}
