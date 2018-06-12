package instructions // import "vimagination.zapto.org/instructions"

import "io"

// New creates a new instruction parser from the given value - exported
// methods on which will be turned into instructions
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
		token, err := t.GetToken()
		if funcToken.typ != 0 && (token.typ == tokenFunction || token.typ == tokenComment || token.typ == tokenDone) {
			function, err := f.bind(funcToken.data, argTokens...)
			if err != nil {
				return functions, err
			}
			functions = append(functions, function)
			argTokens = argTokens[:0]
			funcToken.typ = 0
		}

		if err != nil {
			if err == io.EOF {
				break
			}
			return functions, err
		} else if token.typ == tokenDone || token.typ == tokenError {
			break
		}

		switch token.typ {
		case tokenFunction:
			if token.data != "" {
				funcToken = token
			}
		case tokenComment:
			token.typ = tokenString
			function, err := f.bind("Comment", token)
			if err != nil {
				return nil, err
			}
			functions = append(functions, function)
		case tokenString, tokenNumber:
			argTokens = append(argTokens, token)
		}
	}
	return functions, nil
}
