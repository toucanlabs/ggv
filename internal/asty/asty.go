package asty

import (
	"go/parser"
	"io"
	"os"
)

func noOpClose() error {
	return nil
}

func OpenRead(input string) (reader io.Reader, close func() error, err error) {
	if input == "" || input == "-" {
		return os.Stdin, noOpClose, nil
	}
	f, err := os.Open(input)
	if err != nil {
		return nil, nil, err
	}
	return f, func() error {
		return f.Close()
	}, nil
}

type Options struct {
	WithPositions  bool
	WithComments   bool
	WithReferences bool
	WithImports    bool
}

func ParseFile(input string, options Options) (map[any]any, map[any][]any) {
	marshaller := NewAstParser(options)

	mode := parser.SkipObjectResolution
	if options.WithComments {
		mode |= parser.ParseComments
	}

	inFile, closeIn, err := OpenRead(input)
	if err != nil {
		return nil, nil
	}
	defer closeIn()

	tree, err := parser.ParseFile(marshaller.FileSet(), input, inFile, mode)
	if err != nil {
		return nil, nil
	}

	marshaller.ParseFile(tree)
	return marshaller.internalFuncs, marshaller.refFuncs
}
