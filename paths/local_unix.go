//go:build !windows
// +build !windows

package paths

const root = "/"

func (l local) Root() string           { return root }
func (l local) SupportsSymlinks() bool { return true }
