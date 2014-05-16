package instructions

import "fmt"

type noFunction struct {
	functionName string
}

func (n noFunction) Error() string {
	return "function not found: " + n.functionName
}

type duplicateFunction struct {
	functionName string
}

func (d duplicateFunction) Error() string {
	return "function already exists: " + d.functionName
}

type invalidNumArgs struct {
	name           string
	variadic       bool
	expecting, got int
}

func (i invalidNumArgs) Error() string {
	if i.variadic {
		return fmt.Sprintf("function %s: expecting %d (or more) arguments, got %d", i.name, i.expecting, i.got)
	} else {
		return fmt.Sprintf("function %s: expecting %d arguments, got %d", i.name, i.expecting, i.got)
	}
}

type invalidArgType struct {
	expecting, got string
}

func (i invalidArgType) Error() string {
	return "expecting argument of type " + i.expecting + ", got " + i.got
}

type invalidMethod struct {
	methodName string
}

func (i invalidMethod) Error() string {
	return "method did not meet requirements: " + i.methodName
}
