package traceable

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func Test_Generator_importPath(t *testing.T) {
	tests := []struct {
		name               string
		typeName           string
		rootPackage        string
		expectedImportPath string
	}{
		{
			name:               "returns import path",
			typeName:           "foobar/foo/bar.Barz",
			rootPackage:        "root/package/path",
			expectedImportPath: "foobar/foo/bar",
		},
		{
			name:               "returns root package",
			typeName:           "FooBarz",
			rootPackage:        "root/package/path",
			expectedImportPath: "root/package/path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := Generator{RootPackage: tt.rootPackage}

			qt.Check(t, g.importPath(tt.typeName), qt.Equals, tt.expectedImportPath)
		})
	}
}

func Test_getStructName(t *testing.T) {
	tests := []struct {
		name       string
		typeName   string
		structName string
	}{
		{
			name:       "bare type",
			typeName:   "Barz",
			structName: "Barz",
		},
		{
			name:       "type with package name",
			typeName:   "foo.Barz",
			structName: "Barz",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qt.Check(t, getStructName(tt.typeName), qt.Equals, tt.structName)
		})
	}
}
