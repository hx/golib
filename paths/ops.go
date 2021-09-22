package paths

import (
	"os"
)

type RenameEvent struct{ TargetEvent }
type SymlinkEvent struct{ TargetEvent }
type CopyEvent struct{ TargetEvent }
type CopyOverEvent struct{ TargetEvent }

func (p *Path) Rename(target *Path) (err error) {
	err = p.tree.sys.Rename(p.path, target.path)
	if err == nil {
		p.tree.dispatch(RenameEvent{newTargetEvent(p, target)})
	}
	return
}
func (p *Path) MustRename(target *Path) { must(p.Rename(target)) }

func (p *Path) SymlinkTo(target *Path) error {
	sys := p.tree.sys
	if !sys.SupportsSymlinks() {
		return p.CopyTo(target)
	}
	if target.IsNonDir() {
		stat, err := target.Stat()
		if err != nil {
			return err
		}
		if stat.Mode()&os.ModeSymlink != 0 {
			oldTarget, err := target.ReadLink()
			if err != nil {
				return err
			}
			if oldTarget.path == p.path {
				return nil
			}
		}
		if err = target.Delete(); err != nil {
			return err
		}
	}
	if err := sys.Symlink(p.path, target.path); err != nil {
		return err
	}
	p.tree.dispatch(SymlinkEvent{newTargetEvent(p, target)})
	return nil
}
func (p *Path) MustSymlinkTo(target *Path) { must(p.SymlinkTo(target)) }

func (p *Path) CopyTo(target *Path) error {
	existed := target.Exists()
	if existed {
		if equal, err := p.BytesAreEqual(target); err != nil {
			return err
		} else if equal {
			return nil
		}
	}
	writer, err := p.tree.sys.OpenFile(target.path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	err = p.ReadTo(writer)
	if err == nil {
		err = writer.Close()
	} else {
		_ = writer.Close()
	}
	if err != nil {
		return err
	}
	if existed {
		p.tree.dispatch(CopyOverEvent{newTargetEvent(p, target)})
	} else {
		p.tree.dispatch(CopyEvent{newTargetEvent(p, target)})
	}
	return nil
}
