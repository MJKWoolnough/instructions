package instructions

import "fmt"

type noFunctionError struct {
	functionName string
}

func (n noFunctionError) Error() string {
	return "function not found: " + n.functionName
}

type duplicateFunctionError struct {
	functionName string
}

func (d duplicateFunctionError) Error() string {
	return "function already exists: " + d.functionName
}

type invalidNumArgsError struct {
	name           string
	variadic       bool
	expecting, got int
}

func (i invalidNumArgsError) Error() string {
	if i.variadic {
		return fmt.Sprintf("function %s: expecting %d (or more) arguments, got %d", i.name, i.expecting, i.got)
	}

	return fmt.Sprintf("function %s: expecting %d arguments, got %d", i.name, i.expecting, i.got)
}

type invalidArgTypeError struct {
	expecting, got string
}

func (i invalidArgTypeError) Error() string {
	return "expecting argument of type " + i.expecting + ", got " + i.got
}

type invalidMethodError struct {
	methodName string
}

func (i invalidMethodError) Error() string {
	return "method did not meet requirements: " + i.methodName
}
