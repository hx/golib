package paths_test

import (
	"fmt"
	"github.com/hx/golib/paths"
	. "github.com/hx/golib/testing"
	"runtime"
	"testing"
)

func TestPath_Glob(t *testing.T) {
	dir := pathTestFile().Parent()
	fmt.Println(dir)
	glob, err := dir.Glob("*_test.go")
	Ok(t, err)
	_, expected, _, _ := runtime.Caller(0)
	Assert(
		t,
		glob.Any(func(path *paths.Path) bool { return path.String() == expected }),
		"glob should include this file",
	)
}
