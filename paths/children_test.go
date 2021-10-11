package paths_test

import (
	"github.com/hx/golib/paths"
	. "github.com/hx/golib/testing"
	"runtime"
	"strings"
	"testing"
)

func reslash(str string) string {
	if runtime.GOOS == "windows" {
		return strings.ReplaceAll(str, "/", "\\")
	}
	return str
}

func TestPath_Glob(t *testing.T) {
	dir := pathTestFile().Parent()
	glob, err := dir.Glob("*_test.go")
	Ok(t, err)
	_, thisFile, _, _ := runtime.Caller(0)
	expected := reslash(thisFile)
	Assert(
		t,
		glob.Any(func(path *paths.Path) bool { return path.String() == expected }),
		"glob should include "+expected,
	)
}
