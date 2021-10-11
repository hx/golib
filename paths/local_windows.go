//go:build windows
// +build windows

package paths

import (
	"os"
	"strings"
)

var root = strings.ToUpper(os.Getenv("SYSTEMDRIVE")) + "\\"

func (l local) Root() string           { return root }
func (l local) SupportsSymlinks() bool { return false }
