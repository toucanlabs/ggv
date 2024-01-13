package asty

import (
	"encoding/json"
	"go/parser"
)

type Options struct {
	WithPositions  bool
	WithComments   bool
	WithReferences bool
	WithImports    bool
}

func SourceToJSON(input string, options Options) map[string]interface{} {
	marshaller := NewMarshaller(options)

	mode := parser.SkipObjectResolution
	if options.WithComments {
		mode |= parser.ParseComments
	}

	inFile, closeIn, err := OpenRead(input)
	if err != nil {
		return nil
	}
	defer closeIn()

	tree, err := parser.ParseFile(marshaller.FileSet(), input, inFile, mode)
	if err != nil {
		return nil
	}

	node := marshaller.MarshalFile(tree)
	b, err := node.MarshalJSON()
	if err != nil {
		return nil
	}
	var obj map[string]interface{}
	err = json.Unmarshal(b, &obj)
	if err != nil {
		return nil
	}

	return obj
}
