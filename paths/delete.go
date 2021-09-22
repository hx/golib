package paths

type DeleteFileEvent struct{ FileEvent }
type DeleteDirEvent struct{ Event }

func (p *Path) Delete() (err error) {
	stat, err := p.Stat()
	if err != nil {
		return
	}
	sys := p.tree.sys
	if stat.IsDir() {
		err = sys.RemoveAll(p.path)
		if err == nil {
			p.tree.dispatch(DeleteDirEvent{newEvent(p)})
		}
		return
	}
	err = sys.Remove(p.path)
	if err == nil {
		p.tree.dispatch(DeleteFileEvent{newFileEvent(p, stat.Size())})
	}
	return
}

func (p *Path) DeleteIfExists() error {
	if p.Exists() {
		return p.Delete()
	}
	return nil
}

func (p *Path) MustDelete()         { must(p.Delete()) }
func (p *Path) MustDeleteIfExists() { must(p.DeleteIfExists()) }

func (p *Path) KeepChildren(childNames ...string) (err error) {
	children, err := p.Children()
	if err != nil {
		return
	}

	if len(children) == 0 {
		return
	}

	childNameMap := make(map[string]struct{}, len(childNames))
	for _, name := range childNames {
		childNameMap[name] = struct{}{}
	}

	for _, child := range children {
		if _, found := childNameMap[child.Base()]; !found {
			if err = child.Delete(); err != nil {
				return
			}
		}
	}

	return
}
func (p *Path) MustKeepChildren(childNames ...string) { must(p.KeepChildren(childNames...)) }
