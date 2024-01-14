package parser

import (
	"fmt"
	"go/build"
	"log"
	"os"
	"strings"

	"github.com/toucan-labs/ggv/internal/utils"
)

func NewParser() *Parser {
	p := &Parser{
		ctx:          build.Default,
		withGoroot:   true,
		ignoreVendor: true,
		ignoreStdlib: true,
		ignored: map[string]bool{
			"C": true,
		},
	}
	return p
}

type Parser struct {
	ctx          build.Context
	ignored      map[string]bool
	withGoroot   bool
	ignoreVendor bool
	ignoreStdlib bool
}

type PkgRule func(*build.Package) bool

func (p *Parser) Parse(pkgName string) []*Pkg {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get cwd: %s", err)
	}

	var pkgs map[string]*Pkg = make(map[string]*Pkg)

	err = p.lookupPackage(pkgs, cwd, pkgName, func(pk *build.Package) bool {
		paths := []string{
			"go/pkg/mod",
			"google.golang.org",
			"golang.org",
		}

		result := (pk.Goroot && !p.withGoroot) ||
			p.isIgnored(pk) ||
			utils.InOps[string](pk.Dir, paths, func(el, i string) bool {
				return strings.Contains(i, el)
			})
		return result
	})

	if err != nil {
		log.Fatal(err)
	}
	list := []*Pkg{}
	for _, el := range pkgs {
		list = append(list, el)
	}
	return list
}

func (p *Parser) getPkg(root string, pkgName string, isIgnore PkgRule) *build.Package {
	pkg, buildErr := p.ctx.Import(pkgName, root, 0)
	if buildErr != nil {
		return nil
	}
	if isIgnore(pkg) {
		return nil
	}
	return pkg
}

func (p *Parser) lookupPackage(pkgs map[string]*Pkg, root string, pkgName string, isIgnore PkgRule) (err error) {
	pkg, buildErr := p.ctx.Import(pkgName, root, 0)
	if buildErr != nil {
		err = fmt.Errorf("failed to import %s:\n%s", pkgName, buildErr)
		return
	}

	if isIgnore(pkg) {
		return nil
	}

	newPkg := &Pkg{
		Dir:   pkg.Dir,
		Name:  p.normalizeVendor(pkgName),
		Files: pkg.GoFiles,
	}

	importPath := p.normalizeVendor(pkgName)
	pkgs[importPath] = newPkg

	for _, imp := range p.getImports(pkg) {
		imported := p.getPkg(pkg.Dir, imp, isIgnore)
		if imported != nil && !isIgnore(imported) {
			newPkg.UsagePkgs = append(newPkg.UsagePkgs, imp)
		}
		if _, ok := pkgs[imp]; !ok {
			err = p.lookupPackage(pkgs, pkg.Dir, imp, isIgnore)
			if err != nil {
				return
			}
		}

	}
	return
}

func (p *Parser) getImports(pkg *build.Package) (imports []string) {
	found := make(map[string]struct{})
	for _, imp := range pkg.Imports {
		if _, ok := found[imp]; ok {
			continue
		}
		found[imp] = struct{}{}
		imports = append(imports, imp)
	}
	return
}

func (p *Parser) isIgnored(pkg *build.Package) bool {
	if p.ignoreVendor && p.isVendored(pkg.ImportPath) {
		return true
	}
	path := p.normalizeVendor(pkg.ImportPath)
	if _, ok := p.ignored[path]; ok {
		return true
	}
	return pkg.Goroot && p.ignoreStdlib
}

func (p *Parser) isVendored(path string) bool {
	return strings.Contains(path, "/vendor/")
}

func (p *Parser) normalizeVendor(path string) string {
	pieces := strings.Split(path, "vendor/")
	return pieces[len(pieces)-1]
}
