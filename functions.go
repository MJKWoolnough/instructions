package instructions

import (
	"fmt"
	"reflect"
)

var errType = reflect.TypeOf((*error)(nil)).Elem()

type functionDesc struct {
	variadic  bool
	function  *reflect.Value
	arguments []argument
}

func comment(c string) {
	fmt.Println(c)
}

type functions map[string]functionDesc

func (f functions) AddFunctions(i interface{}) error {
	v := reflect.ValueOf(i)
	t := v.Type()

	for i := 0; i < v.NumMethod(); i++ {
		f.add(t.Method(i).Name, v.Method(i))
	}

	f.add("Comment", reflect.ValueOf(comment))

	return nil
}

func (f functions) add(name string, fn reflect.Value) error {
	if _, ok := f[name]; ok {
		return &duplicateFunctionError{name}
	}

	t := fn.Type()
	variadic := t.IsVariadic()

	if no := t.NumOut(); no > 1 {
		return &invalidMethodError{name}
	} else if no == 1 && !t.Out(0).Implements(errType) {
		return &invalidMethodError{name}
	}

	args := make([]argument, t.NumIn())

	for i := 0; i < len(args); i++ {
		argType := t.In(i)

		if variadic && i == len(args)-1 {
			argType = argType.Elem()
		}

		var a argument

		switch argType.Kind() {
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			a = &intT{bitSize: argType.Bits()}
		case reflect.Int:
			a = &intT{bitSize: 0}
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			a = &uintT{bitSize: argType.Bits()}
		case reflect.Uint:
			a = &uintT{bitSize: 0}
		case reflect.Float32, reflect.Float64:
			a = &floatT{bitSize: argType.Bits()}
		case reflect.String:
			a = new(stringT)
		default:
			return &invalidMethodError{name}
		}

		args[i] = a
	}

	f[name] = functionDesc{
		variadic:  variadic,
		function:  &fn,
		arguments: args,
	}

	return nil
}

func (f functions) bind(name string, args ...token) (Function, error) {
	fd, ok := f[name]
	if !ok {
		return nil, &noFunctionError{
			functionName: name,
		}
	} else if len(args) < len(fd.arguments)-1 || (len(args) > len(fd.arguments)) != fd.variadic {
		return nil, &invalidNumArgsError{
			name:      name,
			got:       len(args),
			expecting: len(fd.arguments),
			variadic:  fd.variadic,
		}
	}

	fp := &function{
		name:      name,
		function:  fd.function,
		arguments: make([]reflect.Value, len(args)),
	}

	for i := 0; i < len(args); i++ {
		j := i
		if j >= len(fd.arguments) {
			j = len(fd.arguments) - 1
		}

		err := fd.arguments[j].Set(args[i].data)
		if err != nil {
			return nil, err
		}

		fp.arguments[i] = fd.arguments[j].Get()
	}

	return fp, nil
}

// Function is used to interface the instructions to the methods.
type Function interface {
	Call() error
	Name() string
}

type function struct {
	name      string
	function  *reflect.Value
	arguments []reflect.Value
}

func (f *function) Call() error {
	returns := f.function.Call(f.arguments)
	if len(returns) == 0 {
		return nil
	}

	err := returns[0].Interface()
	if err == nil {
		return nil
	}

	return err.(error)
}

func (f *function) Name() string {
	return f.name
}
