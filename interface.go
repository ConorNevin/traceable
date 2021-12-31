package traceable

type Interface struct {
	name       string
	importPath string
	methods    []Method
	imports    map[string]ImportedPackage
}

func (i Interface) hasMethod(m Method) bool {
	for _, im := range i.methods {
		if im.name == m.name {
			return true
		}
	}

	return false
}
