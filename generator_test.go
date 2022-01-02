package traceable

import (
	"fmt"
	"regexp"
	"strings"
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

func TestGenerator_generate(t *testing.T) {
	pm := map[string]string{
		"context":                               "context",
		"github.com/opentracing/opentracing-go": "opentracing",
	}

	tests := []struct {
		name              string
		inter             Interface
		expectedFunctions map[string]string
		expectedImports   []string
	}{
		{
			name: "One function interface",
			inter: Interface{
				name: "FooBar",
				methods: []Method{
					{
						name: "Foo",
						args: []*Parameter{
							{
								typ: &Type{pkg: "context", value: "Context"},
							},
						},
						returns: []*Parameter{
							{
								typ: &Type{value: "error"},
							},
						},
					},
				},
			},
			expectedFunctions: map[string]string{
				"Foo": "func (t *TracedFooBar) Foo(a0 context.Context) error {",
			},
			expectedImports: []string{
				"context",
				"github.com/opentracing/opentracing-go",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Generator{
				packageMap: pm,
				Interface:  tt.inter,
			}
			g.generate(tt.inter.name)

			lines := strings.Split(g.buf.String(), "\n")

			for _, method := range g.Interface.methods {
				idx := findMethodLines(t, method.name, lines)

				got, ok := tt.expectedFunctions[method.name]
				if !ok {
					continue
				}

				qt.Check(t, lines[idx], qt.Equals, got)
			}

			qt.Check(t, findImports(t, lines), qt.ContentEquals, tt.expectedImports)
		})
	}
}

func findMethodLines(t *testing.T, methodName string, lines []string) int {
	t.Helper()
	r := regexp.MustCompile(fmt.Sprintf(`func\s+\(.*\)\s*%s`, methodName))
	for i, line := range lines {
		if r.MatchString(line) {
			return i
		}
	}

	t.Fatalf("unable to find 'func (.*) %s'", methodName)
	return -1 // unreachable
}

func findImports(t *testing.T, lines []string) []string {
	t.Helper()

	var (
		foundImportBlock bool
		imports          []string
	)

	for _, line := range lines {
		switch line {
		case "import(":
			foundImportBlock = true
		case ")":
			return imports
		default:
			if foundImportBlock {
				imports = append(imports, line[1:(len(line)-1)])
			}
		}
	}

	t.Fatal("unable to find imported block")
	return nil // unreachable
}
