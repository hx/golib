// +build windows

package paths

func (l local) Root() string           { return os.GetEnv("SYSTEMDRIVE") + "\\" }
func (l local) SupportsSymlinks() bool { return false }
