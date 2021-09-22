package paths

func (p *Path) Children() (paths Paths, err error) {
	stat, err := p.Stat()
	if err != nil {
		return
	}
	if !stat.IsDir() {
		return nil, ErrNonDirectory
	}
	entries, err := p.tree.sys.ReadDir(p.path)
	if err != nil {
		return
	}
	paths = make(Paths, len(entries))
	for i, entry := range entries {
		paths[i] = p.Join(entry.Name())
	}
	return
}

func (p *Path) ChildrenIfIsDir() (Paths, error) {
	if p.IsDir() {
		return p.Children()
	}
	return nil, nil
}

func (p *Path) MustChildren() Paths        { return must1(p.Children()).(Paths) }
func (p *Path) MustChildrenIfIsDir() Paths { return must1(p.ChildrenIfIsDir()).(Paths) }

func (p *Path) Glob(pattern string) (paths Paths, err error) {
	sys := p.tree.sys
	entries, err := sys.Glob(sys.Join(p.path, pattern))
	if err != nil {
		return
	}
	paths = make(Paths, len(entries))
	for i, entry := range entries {
		paths[i] = p.Join(entry)
	}
	return
}
func (p *Path) MustGlob(pattern string) Paths { return must1(p.Glob(pattern)).(Paths) }
