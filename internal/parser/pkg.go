package parser

import (
	"fmt"
	"path/filepath"

	"github.com/toucan-labs/ggv/internal/asty"
)

type Pkg struct {
	UsagePkgs []string
	Name      string
	Dir       string
	Files     []string
}

func (p *Pkg) Funcs() (nodes, links []interface{}) {
	for _, el := range p.Files {
		filepathString := filepath.Join(p.Dir, el)
		_, ref := asty.ParseFile(filepathString, asty.Options{
			WithPositions:  false,
			WithComments:   false,
			WithReferences: true,
			WithImports:    false,
		})

		for internalFunName, list := range ref {
			if len(list) > 0 {
				fn := fmt.Sprintf("%s::%s", el, internalFunName)
				nodes = append(nodes, fn)
				for _, el := range list {
					nodes = append(nodes, el)
					links = append(links, map[string]interface{}{
						"source": fn,
						"target": el,
					})
				}
			}
		}
	}

	return
}

type Graph interface {
	Data(includeInternalFuncs bool, pkgs []*Pkg) map[string]interface{}
}

func NewGraph() Graph {
	return &graph{}
}

type graph struct{}

func (g *graph) Data(includeInternalFuncs bool, pkgs []*Pkg) (result map[string]interface{}) {
	nodes := []interface{}{}
	links := []interface{}{}
	for _, e := range pkgs {
		nodes = append(nodes, map[string]interface{}{
			"id":    e.Name,
			"val":   e.Name,
			"group": 1,
		})
		for _, p := range e.UsagePkgs {
			links = append(links, map[string]interface{}{
				"source": e.Name,
				"target": p,
			})
		}

		if includeInternalFuncs {
			n, l := e.Funcs()
			for _, el := range n {
				nodes = append(nodes, map[string]interface{}{
					"id":    el,
					"val":   el,
					"group": 2,
				})
				links = append(links, map[string]interface{}{
					"source": e.Name,
					"target": el,
				})
			}
			links = append(links, l...)
		}

	}
	return map[string]interface{}{
		"nodes": nodes,
		"links": links,
	}
}
