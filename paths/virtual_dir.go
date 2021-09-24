package paths

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type virtualDir struct {
	*virtualEntryBase
	children []virtualEntry
}

func newVirtualDir() *virtualDir {
	now := time.Now()
	return &virtualDir{
		virtualEntryBase: &virtualEntryBase{
			mode:     0744 | fs.ModeDir,
			accessed: now,
			modified: now,
		},
	}
}

func (v *virtualDir) Size() int64                { return 0 }
func (v *virtualDir) IsDir() bool                { return true }
func (v *virtualDir) Info() (fs.FileInfo, error) { return v, nil }

func (v *virtualDir) root() *virtualDir {
	for v.parent != nil {
		v = v.parent
	}
	return v
}

func (v *virtualDir) resolve(name string) virtualEntry {
	if filepath.IsAbs(name) {
		if v.parent == nil {
			name = name[len(filepath.VolumeName(name))+1:]
		} else {
			return v.root().resolve(name)
		}
	}
	parts := strings.SplitN(name, string(os.PathSeparator), 2)
	if parts[0] == "." || parts[0] == "" {
		if len(parts) == 1 {
			return v
		}
		return v.resolve(parts[1])
	}
	if parts[0] == ".." {
		if len(parts) == 1 || v.parent == nil {
			return v.parent
		}
		return v.parent.resolve(parts[1])
	}
	for _, child := range v.children {
		if child.entry().name != parts[0] {
			continue
		}
		if len(parts) == 1 {
			return child
		}
		switch child := child.(type) {
		case *virtualDir:
			return child.resolve(parts[1])
		case *virtualSymlink:
			target := child.resolveRecursive()
			if target, ok := target.(*virtualDir); ok {
				return target.resolve(parts[1])
			}
		}
		return nil
	}
	return nil
}
