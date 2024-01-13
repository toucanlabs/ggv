package parser

import (
	"path/filepath"

	"github.com/toucan-labs/ggv/internal/asty"
)

type Pkg struct {
	Name  string
	Dir   string
	Files []string
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
