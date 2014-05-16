package instructions

import "testing"

type testErr struct {
	err string
}

func (t *testErr) Error() string {
	return t.err
}

type testFunctions struct {
	testA, testB, testC int
}

func (t *testFunctions) TestFuncA() {
	t.testA = 1
}

func (t *testFunctions) TestFuncB() error {
	return &testErr{"TestFuncB"}
}

func (t *testFunctions) TestFuncC(a int, b int) {
	t.testB = a + b
}

func (t *testFunctions) TestFuncD(a int16, b float32) error {
	return nil
}

func (t *testFunctions) TestFuncE(a int32, b ...byte) error {
	return nil
}

func (t *testFunctions) TestFuncF() error {
	return nil
}

func TestAddFunction(t *testing.T) {
	f := make(functions)
	tf := new(testFunctions)
	f.AddFunctions(tf)
	if len(f) != 7 { //Comment is automatically added
		t.Errorf("expecting 6 methods, got %d", len(f))
	}
	tests := []struct {
		name     string
		args     int
		variadic bool
	}{
		{"TestFuncA", 0, false},
		{"TestFuncB", 0, false},
		{"TestFuncC", 2, false},
		{"TestFuncD", 2, false},
		{"TestFuncE", 2, true},
	}
	for _, test := range tests {
		if m, ok := f[test.name]; !ok {
			t.Errorf("method %q does not exists", test.name)
		} else if len(m.arguments) != test.args {
			t.Errorf("expecting %d arguments, got %d", test.args, len(m.arguments))
		} else if m.variadic != test.variadic {
			t.Errorf("expecting variadic: %v, got variadic %v", test.variadic, m.variadic)
		}
	}
}

func TestCallFunction(t *testing.T) {
	f := make(functions)
	tf := new(testFunctions)
	f.AddFunctions(tf)
	a, err := f.bind("TestFuncA")
	if err != nil {
		t.Errorf("unexpected error: %q", err)
		return
	}
	err = a.Call()
	if err != nil {
		t.Errorf("unexpected error: %q", err)
		return
	} else if tf.testA != 1 {
		t.Errorf("expecting testA = 1, got %d", tf.testA)
		return
	}
	b, err := f.bind("TestFuncB")
	if err != nil {
		t.Errorf("unexpected error: %q", err)
		return
	}
	err = b.Call()
	if err == nil {
		t.Errorf("expecting error, nil received")
		return
	} else if _, ok := err.(*testErr); !ok {
		t.Errorf("unexpected error: %q", err)
		return
	}
	ff, err := f.bind("TestFuncF")
	if err != nil {
		t.Errorf("unexpected error: %q", err)
		return
	}
	err = ff.Call()
	if err != nil {
		t.Errorf("expecting nil, error received")
		return
	}
}
