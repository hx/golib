package paths

import "os"

func (p *Path) Make() error                     { return p.MakeMode(0755) }
func (p *Path) MakeMode(mode os.FileMode) error { return p.tree.sys.MkdirAll(p.path, mode) }

func (p *Path) MustMake() *Path {
	must(p.Make())
	return p
}

func (p *Path) MustMakeMode(mode os.FileMode) *Path {
	must(p.MakeMode(mode))
	return p
}

func (p *Path) removeEmptyDirs(recursive bool) (removed Paths, err error) {
	children, err := p.Children()
	if err != nil {
		return
	}
	var empty bool
	for _, child := range children {
		if !child.IsDir() {
			continue
		}
		if recursive {
			var removedFromChild Paths
			removedFromChild, err = child.removeEmptyDirs(true)
			if err != nil {
				return
			}
			removed = append(removed, removedFromChild...)
		}
		empty, err = child.IsEmpty()
		if !empty || err != nil {
			continue
		}
		err = child.Delete()
		if err == nil {
			removed = append(removed, child)
		}
	}
	return
}

func (p *Path) removeEmptyDirsAndSelf(recursive bool) (removed Paths, err error) {
	removed, err = p.removeEmptyDirs(recursive)
	if err != nil {
		return
	}
	empty, err := p.IsEmpty()
	if err != nil {
		return
	}
	if empty {
		err = p.Delete()
	}
	if err == nil {
		removed = append(removed, p)
	}
	return
}

func (p *Path) RemoveEmptyDirs() (removed Paths, err error) {
	return p.removeEmptyDirs(false)
}
func (p *Path) RemoveEmptyDirsRecursive() (removed Paths, err error) {
	return p.removeEmptyDirs(true)
}
func (p *Path) MustRemoveEmptyDirs() (removed Paths) {
	return must1(p.removeEmptyDirs(false)).(Paths)
}
func (p *Path) MustRemoveEmptyDirsRecursive() (removed Paths) {
	return must1(p.removeEmptyDirs(true)).(Paths)
}
func (p *Path) RemoveEmptyDirsAndSelf() (removed Paths, err error) {
	return p.removeEmptyDirsAndSelf(false)
}
func (p *Path) RemoveEmptyDirsRecursiveAndSelf() (removed Paths, err error) {
	return p.removeEmptyDirsAndSelf(true)
}
func (p *Path) MustRemoveEmptyDirsAndSelf() (removed Paths) {
	return must1(p.removeEmptyDirsAndSelf(false)).(Paths)
}
func (p *Path) MustRemoveEmptyDirsRecursiveAndSelf() (removed Paths) {
	return must1(p.removeEmptyDirsAndSelf(true)).(Paths)
}
