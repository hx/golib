package paths

type Paths []*Path

func (p Paths) Dirs() (dirs Paths) {
	for _, entry := range p {
		if entry.IsDir() {
			dirs = append(dirs, entry)
		}
	}
	return
}

func (p Paths) NonDirs() (nonDirs Paths) {
	for _, entry := range p {
		if entry.IsNonDir() {
			nonDirs = append(nonDirs, entry)
		}
	}
	return
}

func (p Paths) Any(predicate func(path *Path) bool) bool {
	for _, path := range p {
		if predicate(path) {
			return true
		}
	}
	return false
}

func (p Paths) All(predicate func(path *Path) bool) bool {
	for _, path := range p {
		if !predicate(path) {
			return false
		}
	}
	return true
}
