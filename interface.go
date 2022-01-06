package traceable

type Interface struct {
	name       string
	importPath string
	methods    []Method
}

func (i *Interface) hasMethod(m Method) bool {
	for _, im := range i.methods {
		if im.name == m.name {
			return true
		}
	}

	return false
}

func (i *Interface) imports() map[string]struct{} {
	imports := make(map[string]struct{})
	for _, m := range i.methods {
		mi := m.imports()
		for ip := range mi {
			imports[ip] = struct{}{}
		}
	}
	return imports
}
