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

func (p *Pkg) Funcs() (result []interface{}) {
	for _, el := range p.Files {
		filepathString := filepath.Join(p.Dir, el)
		// return external, internal funcs in file
		jsonResult := asty.ParseFile(filepathString, asty.Options{
			WithPositions:  false,
			WithComments:   false,
			WithReferences: true,
			WithImports:    false,
		})

		for k, _ := range jsonResult {
			result = append(result, fmt.Sprintf("%s::%s", el, k))
		}
	}

	return
}

type Graph interface {
	Data(pkgs []*Pkg) map[string]interface{}
}

func NewGraph() Graph {
	return &graph{}
}

type graph struct{}

func (g *graph) Data(pkgs []*Pkg) (result map[string]interface{}) {
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

		for _, el := range e.Funcs() {
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

	}
	return map[string]interface{}{
		"nodes": nodes,
		"links": links,
	}
}
