// +build !windows

package paths

func (l local) Root() string           { return "/" }
func (l local) SupportsSymlinks() bool { return true }
