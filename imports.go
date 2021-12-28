package traceable

import (
	"fmt"
	"strings"

	"golang.org/x/tools/go/packages"
)

type ImportedPackage struct {
	Name       string
	Path       string
	Duplicates []string
}

func createPackageMap(importPaths []string) (map[string]string, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName |
			packages.NeedTypesInfo |
			packages.NeedSyntax |
			packages.NeedTypes,
	}
	pkgs, err := packages.Load(cfg, importPaths...)
	if err != nil {
		return nil, err
	}

	pkgMap := make(map[string]string)
	found := map[string]struct{}{}
	for _, pkg := range pkgs {
		pkgMap[pkg.PkgPath] = pkg.Name
		found[pkg.PkgPath] = struct{}{}
	}

	var missing []string
	for _, importPath := range importPaths {
		if _, ok := found[importPath]; !ok {
			missing = append(missing, importPath)
		}
	}

	if len(missing) > 0 {
		return nil, fmt.Errorf("failed to load packages: %s", strings.Join(missing, ", "))
	}

	return pkgMap, nil
}
