package traceable

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func Test_createPackageMap(t *testing.T) {
	tests := []struct {
		name        string
		importPath  string
		wantPackage string
	}{
		{"standard library", "context", "context"},
		{"third party", "golang.org/x/tools/present", "present"},
	}

	importPaths := make([]string, len(tests))
	for i := range importPaths {
		importPaths[i] = tests[i].importPath
	}

	packageMap, err := createPackageMap(importPaths)
	qt.Assert(t, err, qt.IsNil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPackageName, gotOk := packageMap[tt.importPath]
			qt.Assert(t, gotOk, qt.IsTrue)
			qt.Check(t, gotPackageName, qt.Equals, tt.wantPackage)
		})
	}
}
