package instructions

import (
	"fmt"
	"strings"
	"testing"
)

type instructs string

func (i *instructs) Func1() {
	(*i) = instructs("Func1")
}

func (i *instructs) Func2(arg1 string) {
	(*i) = instructs(arg1)
}

func (i *instructs) Func3(arg1 string, arg2 string) {
	(*i) = instructs(arg1 + arg2)
}

func (i *instructs) Func4(arg1 int, arg2 byte) {
	(*i) = instructs(fmt.Sprintf("%d - %d", arg1, arg2))
}

func (i *instructs) Func5(arg1 ...byte) {
	(*i) = instructs(arg1)
}

func (i *instructs) Func6(arg1 string, arg2 ...float32) {
	(*i) = instructs(fmt.Sprintf("%s - %v", arg1, arg2))
}

func TestInstructions(t *testing.T) {
	str := "Func1\nFunc2 \"abcxyz\"\nFunc2 \"ABCXYZ\"\nFunc3 \"123\" \"456\"\nFunc3 \"456\" \"123\"\nFunc4 14513 0x7f\nFunc5 72 101 108 0x6c 0x6f\nFunc6 \"Beep\" 3.14159 1.618\n"

	data := new(instructs)

	ins, err := New(data, strings.NewReader(str))

	if err != nil {
		fmt.Println(ins)
		t.Errorf("unexpected error: %q", err)
		return
	}

	tests := [...]struct {
		funcName, result string
	}{
		{"Func1", "Func1"},
		{"Func2", "abcxyz"},
		{"Func2", "ABCXYZ"},
		{"Func3", "123456"},
		{"Func3", "456123"},
		{"Func4", "14513 - 127"},
		{"Func5", "Hello"},
		{"Func6", "Beep - [3.14159 1.618]"},
	}
	if len(tests) != len(ins) {
		t.Errorf("expecting %d instructions, got %d", len(tests), len(ins))
		return
	}

	for n, test := range tests {
		ins[n].Call()
		if ins[n].Name() != test.funcName {
			t.Errorf("test %d: expecting function named %q, got %q", n+1, test.funcName, ins[n].Name())
		} else {
			ins[n].Call()
			if string(*data) != test.result {
				t.Errorf("test %d: expecting result %q, got %q", n+1, test.result, (*data))
			}
		}
	}

}
