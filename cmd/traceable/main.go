package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"

	"github.com/ConorNevin/traceable"
)

var (
	typeNames = flag.String("types", "", "comma-separated list of type names; must be set")
	output    = flag.String("output", "", "output file name; default srcdir/traced_<type>.go")
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("traceable: ")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		// Default: process whole package in current directory.
		args = []string{"."}
	}
	types := strings.Split(*typeNames, ",")

	if err := run(args, types); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args, types []string) error {
	g := newGenerator()

	for _, typeName := range types {
		idx := strings.IndexRune(typeName, '.')
		if idx != -1 {
			args = append(args, typeName[:idx])
		}
	}

	g.ParsePackage(args)
	g.GenerateAll(types)

	dst := os.Stdout
	if len(*output) > 0 {
		if err := os.MkdirAll(filepath.Dir(*output), os.ModePerm); err != nil {
			log.Fatalf("unable to create directory: %s", err)
		}
		f, err := os.Create(*output)
		if err != nil {
			log.Fatalf("failed opening destination file: %s", err)
		}
		defer f.Close()
		dst = f
	}

	if _, err := dst.Write(g.Format()); err != nil {
		log.Fatalf("writing output: %s", err)
	}

	return nil
}

func newGenerator() traceable.Generator {
	var g traceable.Generator

	dstPath, err := filepath.Abs(filepath.Dir(*output))
	if err != nil {
		log.Println("unable to determine destination file path:", err)
	}

	pkgPath, err := parsePackageImport(dstPath)
	if err != nil {
		log.Println("unable to infer output package name", err)
	}

	g.OutputPackagePath = pkgPath
	g.RootPackage = getRootPackage()

	return g
}

func getRootPackage() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("get current directory failed: %v", err)
	}
	packageName, err := packageNameOfDir(dir)
	if err != nil {
		log.Fatalf("parse package name failed: %v", err)
	}

	return packageName
}

func parsePackageImport(srcDir string) (string, error) {
	moduleMode := os.Getenv("GO111MODULE")
	// trying to find the module
	if moduleMode != "off" {
		currentDir := srcDir
		for {
			dat, err := ioutil.ReadFile(filepath.Join(currentDir, "go.mod"))
			if os.IsNotExist(err) {
				if currentDir == filepath.Dir(currentDir) {
					// at the root
					break
				}
				currentDir = filepath.Dir(currentDir)
				continue
			} else if err != nil {
				return "", err
			}
			modulePath := modfile.ModulePath(dat)
			return filepath.ToSlash(filepath.Join(modulePath, strings.TrimPrefix(srcDir, currentDir))), nil
		}
	}
	// fall back to GOPATH mode
	goPaths := os.Getenv("GOPATH")
	if goPaths == "" {
		return "", fmt.Errorf("GOPATH is not set")
	}

	goPathList := strings.Split(goPaths, string(os.PathListSeparator))
	for _, goPath := range goPathList {
		sourceRoot := filepath.Join(goPath, "src") + string(os.PathSeparator)
		if strings.HasPrefix(srcDir, sourceRoot) {
			return filepath.ToSlash(strings.TrimPrefix(srcDir, sourceRoot)), nil
		}
	}

	return "", errors.New("package not found")
}

// packageNameOfDir get package import path via dir
func packageNameOfDir(srcDir string) (string, error) {
	files, err := ioutil.ReadDir(srcDir)
	if err != nil {
		log.Fatal(err)
	}

	var goFilePath string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".go") {
			goFilePath = file.Name()
			break
		}
	}
	if goFilePath == "" {
		return "", fmt.Errorf("go source file not found %s", srcDir)
	}

	packageImport, err := parsePackageImport(srcDir)
	if err != nil {
		return "", err
	}
	return packageImport, nil
}
