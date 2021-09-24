package paths

import (
	"io/fs"
	"os"
	"time"
)

type virtualSymlink struct {
	*virtualEntryBase
	target string
}

func newVirtualSymlink(name string, parent *virtualDir, perm os.FileMode, target string) *virtualSymlink {
	now := time.Now()
	return &virtualSymlink{
		target: target,
		virtualEntryBase: &virtualEntryBase{
			name:     name,
			mode:     perm.Perm() | fs.ModeSymlink,
			accessed: now,
			modified: now,
			parent:   parent,
		},
	}
}

func (s *virtualSymlink) Size() int64                { return 0 }
func (s *virtualSymlink) IsDir() bool                { return false }
func (s *virtualSymlink) Info() (fs.FileInfo, error) { return s, nil }

func (s *virtualSymlink) resolve() virtualEntry { return s.parent.resolve(s.target) }

func (s *virtualSymlink) resolveRecursive() virtualEntry {
	var (
		entry virtualEntry = s
		trail              = []*virtualSymlink{s}
	)
	for {
		if link, ok := entry.(*virtualSymlink); ok {
			for _, prev := range trail {
				if link == prev {
					// Recursion
					return nil
				}
			}
			trail = append(trail, link)
			entry = link.resolve()
		} else {
			return entry
		}
	}
}
