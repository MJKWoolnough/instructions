package instructions // import "vimagination.zapto.org/instructions"

import (
	"errors"
	"io"
)

// New creates a new instruction parser from the given value - exported
// methods on which will be turned into instructions.
func New(functionObj interface{}, data io.Reader) ([]Function, error) {
	f := make(functions)

	err := f.AddFunctions(functionObj)
	if err != nil {
		return nil, err
	}

	functions := make([]Function, 0, 32)

	var funcToken token
	argTokens := make([]token, 0, 16)

	t := newLexer(data)

	for {
		gtoken, err := t.GetToken()
		if funcToken.typ != 0 && (gtoken.typ == tokenFunction || gtoken.typ == tokenComment || gtoken.typ == tokenDone) {
			function, ferr := f.bind(funcToken.data, argTokens...)
			if ferr != nil {
				return functions, ferr
			}

			functions = append(functions, function)
			argTokens = argTokens[:0]
			funcToken.typ = 0
		}

		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return functions, err
		} else if gtoken.typ == tokenDone || gtoken.typ == tokenError {
			break
		}

		switch gtoken.typ {
		case tokenFunction:
			if gtoken.data != "" {
				funcToken = gtoken
			}
		case tokenComment:
			gtoken.typ = tokenString

			function, err := f.bind("Comment", gtoken)
			if err != nil {
				return nil, err
			}

			functions = append(functions, function)
		case tokenString, tokenNumber:
			argTokens = append(argTokens, gtoken)
		}
	}

	return functions, nil
}
