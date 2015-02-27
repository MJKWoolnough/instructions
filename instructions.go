package instructions

import (
	"io"

	"github.com/MJKWoolnough/tokeniser"
)

// New creates a new instruction parser from the given struct - exported
// methods on which will be turned into instructions
func New(functionStruct interface{}, data io.Reader) ([]Function, error) {
	f := make(functions)

	err := f.AddFunctions(functionStruct)
	if err != nil {
		return nil, err
	}

	functions := make([]Function, 0, 32)

	var funcToken *tokeniser.Item
	argTokens := make([]*tokeniser.Item, 0, 16)

	t := tokeniser.New(data, lexFunction)
	for {
		token, err := t.Next()

		if funcToken != nil && (token == nil || token.Typ == FUNCTION || token.Typ == COMMENT) {
			function, err := f.bind(funcToken.Val, argTokens...)
			if err != nil {
				return functions, err
			}
			functions = append(functions, function)
			argTokens = argTokens[:0]
			funcToken = nil
		}

		if err != nil {
			if err == io.EOF {
				break
			}
			return functions, err
		} else if token == nil {
			break
		}

		switch token.Typ {
		case FUNCTION:
			if token.Val != "" {
				funcToken = token
			}
		case COMMENT:
			token.Typ = STRING
			function, err := f.bind("Comment", token)
			if err != nil {
				return nil, err
			}
			functions = append(functions, function)
		case STRING, NUMBER:
			argTokens = append(argTokens, token)
		}
	}
	return functions, nil
}
