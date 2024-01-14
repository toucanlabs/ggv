package parser

import (
	"path/filepath"

	"github.com/k0kubun/pp/v3"
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
		jsonResult := asty.SourceToJSON(filepathString, asty.Options{
			WithPositions:  false,
			WithComments:   false,
			WithReferences: true,
			WithImports:    false,
		})

		result = append(result, jsonResult)
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
			"val":    e.Name,
			"group": 1,
		})
		for _, p := range e.UsagePkgs {
			links = append(links,
				map[string]interface{}{
					"source": e.Name,
					"target": p,
				})
		}

	}
	pp.Println(nodes)
	return map[string]interface{}{
		"nodes": nodes,
		"links": links,
	}
}
