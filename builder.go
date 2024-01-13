package main

// import (
// 	"fmt"
// 	"go/build"
// 	"go/parser"
// 	"go/token"
// 	"log"
// 	"os"
// 	"path/filepath"
// 	"sort"
// 	"strings"

// 	"github.com/k0kubun/pp/v3"
// )

// var (
// 	buildTags    []string
// 	ignoreStdlib = false
// 	ignoreVendor = false
// 	maxLevel     = 256
// )

// type builder struct {
// 	ns           []node
// 	ls           []link
// 	erroredPkgs  map[string]bool
// 	ids          map[string]string
// 	buildContext build.Context
// 	ignored      map[string]bool
// }

// func NewBuilder() *builder {
// 	return &builder{
// 		ns:          []node{},
// 		ls:          []link{},
// 		erroredPkgs: make(map[string]bool),
// 		ids:         make(map[string]string),
// 		ignored: map[string]bool{
// 			"C": true,
// 		},
// 		buildContext: build.Default,
// 	}
// }

// const (
// 	STD_PKG = 1
// 	LIB_PKG = 2
// 	PRO_PKG = 3
// )

// func (b *builder) linkJson() []interface{} {
// 	obj := []interface{}{}
// 	for _, el := range b.ls {
// 		obj = append(obj, map[string]interface{}{
// 			"source": el.src,
// 			"target": el.des,
// 			"value":  3,
// 		})
// 	}
// 	return obj
// }

// func (b *builder) nodeJson() []interface{} {
// 	obj := []interface{}{}
// 	for _, el := range b.ns {
// 		obj = append(obj, map[string]interface{}{
// 			"id":    el.id,
// 			"val":   el.pkg,
// 			"group": el.group,
// 		})
// 	}
// 	return obj

// }

// func (b *builder) toJson() map[string]interface{} {
// 	return map[string]interface{}{
// 		"links": b.linkJson(),
// 		"nodes": b.nodeJson(),
// 	}
// }

// func (b *builder) Parse(packageString string) map[string]interface{} {
// 	b.buildContext.BuildTags = buildTags

// 	cwd, err := os.Getwd()
// 	if err != nil {
// 		log.Fatalf("failed to get cwd: %s", err)
// 	}

// 	pkgs, err := b.processPkg(cwd, packageString, 0, "")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// sort packages
// 	pkgKeys := []string{}
// 	for k := range pkgs {
// 		pkgKeys = append(pkgKeys, k)
// 	}

// 	sort.Strings(pkgKeys)

// 	for _, pkgName := range pkgKeys {

// 		pkg := pkgs[pkgName]
// 		pkgId := b.getId(pkgName)

// 		if b.isIgnored(pkg) {
// 			continue
// 		}

// 		if pkg.Goroot && !withGoroot {
// 			continue
// 		}

// 		b.ns = append(b.ns, b.newNode(pkgName, pkg))

// 		for _, imp := range getImports(pkg) {
// 			impPkg := pkgs[imp]
// 			if impPkg == nil || b.isIgnored(impPkg) {
// 				continue
// 			}

// 			if impPkg.Goroot && !withGoroot {
// 				continue
// 			}
// 			b.ls = append(b.ls, b.newLink(pkgId, imp))
// 		}
// 	}
// 	return b.toJson()
// }

// func (b *builder) newNode(pkgName string, pkg *build.Package) node {
// 	pkgId := b.getId(pkgName)
// 	var name string = pkgName
// 	for _, el := range []string{"github.com/", "golang.org/", "google.golang.org/"} {
// 		name = strings.Replace(name, el, "", 1)
// 	}
// 	n := &node{id: pkgId, pkg: name, group: LIB_PKG}

// 	if pkg.Goroot {
// 		n.group = STD_PKG
// 	}

// 	if strings.Contains(pkgId, pkgName) {
// 		n.group = PRO_PKG
// 	}

// 	if n.isProj() {
// 		for _, el := range pkg.GoFiles {
// 			filepathString := filepath.Join(pkg.Dir, el)
// 			if !strings.Contains(filepathString, "go/pkg/mod") {
// 				pp.Println(filepathString)
// 				fset := token.NewFileSet()
// 				pkgs, err := parser.ParseFile(fset, filepathString, nil, 0)
// 				if err != nil {
// 					log.Fatal(err)
// 				}
// 				pp.Println(pkgs)
// 			}
// 		}
// 	}

// 	return *n
// }

// func (b *builder) newLink(src, pkgName string) link {
// 	impId := b.getId(pkgName)
// 	l := link{src: src, des: impId}
// 	return l
// }

// func (b *builder) processPkg(root string, pkgName string, level int, importedBy string) (map[string]*build.Package, error) {
// 	pkgs := make(map[string]*build.Package)
// 	if err := b.processPackage(pkgs, root, pkgName, level, importedBy); err != nil {
// 		return pkgs, err
// 	}
// 	return pkgs, nil

// }

// func (b *builder) processPackage(pkgs map[string]*build.Package,
// 	root string,
// 	pkgName string,
// 	level int,
// 	importedBy string) error {

// 	if level++; level > maxLevel {
// 		return nil
// 	}
// 	if b.ignored[pkgName] {
// 		return nil
// 	}

// 	pkg, buildErr := b.buildContext.Import(pkgName, root, 0)
// 	if buildErr != nil {
// 		return fmt.Errorf("failed to import %s (imported at level %d by %s):\n%s", pkgName, level, importedBy, buildErr)
// 	}

// 	if b.isIgnored(pkg) {
// 		return nil
// 	}

// 	importPath := normalizeVendor(pkgName)
// 	if buildErr != nil {
// 		b.erroredPkgs[importPath] = true
// 	}

// 	pkgs[importPath] = pkg

// 	// Don't worry about dependencies for stdlib packages
// 	if pkg.Goroot && !withGoroot {
// 		return nil
// 	}

// 	for _, imp := range getImports(pkg) {
// 		if _, ok := pkgs[imp]; !ok {
// 			if err := b.processPackage(pkgs, pkg.Dir, imp, level, pkgName); err != nil {
// 				return err
// 			}
// 		}
// 	}
// 	return nil
// }

// func getImports(pkg *build.Package) []string {
// 	allImports := pkg.Imports
// 	var imports []string
// 	found := make(map[string]struct{})
// 	for _, imp := range allImports {
// 		if imp == normalizeVendor(pkg.ImportPath) {
// 			// Don't draw a self-reference when foo_test depends on foo.
// 			continue
// 		}
// 		if _, ok := found[imp]; ok {
// 			continue
// 		}
// 		found[imp] = struct{}{}
// 		imports = append(imports, imp)
// 	}
// 	return imports
// }

// func deriveNodeID(packageName string) string {
// 	id := "\"" + packageName + "\""
// 	return id
// }

// func (b *builder) getId(name string) string {
// 	id, ok := b.ids[name]
// 	if !ok {
// 		id = deriveNodeID(name)
// 		b.ids[name] = id
// 	}
// 	return id
// }

// func (b *builder) isIgnored(pkg *build.Package) bool {

// 	if ignoreVendor && isVendored(pkg.ImportPath) {
// 		return true
// 	}
// 	return b.ignored[normalizeVendor(pkg.ImportPath)] ||
// 		(pkg.Goroot && ignoreStdlib)
// }

// func isVendored(path string) bool {
// 	return strings.Contains(path, "/vendor/")
// }

// func normalizeVendor(path string) string {
// 	pieces := strings.Split(path, "vendor/")
// 	return pieces[len(pieces)-1]
// }
