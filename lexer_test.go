package instructions

import (
	"io"
	"testing"

	"github.com/MJKWoolnough/memio"
)

func TestTokens(t *testing.T) {
	str := []byte("Hello\nWorld 1 2 3\n#Comment 3 5 7 1\nBeep \"This is a String\" 5.4\n#Comment 1\n")
	tokens := newLexer(memio.Open(str))
	tests := []token{
		{tokenFunction, "Hello"},
		{tokenFunction, "World"},
		{tokenNumber, "1"},
		{tokenNumber, "2"},
		{tokenNumber, "3"},
		{tokenComment, "Comment 3 5 7 1"},
		{tokenFunction, "Beep"},
		{tokenString, "This is a String"},
		{tokenNumber, "5.4"},
		{tokenComment, "Comment 1"},
	}
	for n, test := range tests {
		token, err := tokens.GetToken()
		if err != nil {
			t.Errorf("test %d: received unexpected error: %s", n+1, err)
		} else if test.typ != token.typ {
			t.Errorf("test %d: expecting token type %s, got %s", n+1, test.typ, token.typ)
		} else if test.data != token.data {
			t.Errorf("test %d: expecting token value %s, got %s", n+1, test.data, token.data)
		}
	}
	if token, err := tokens.GetToken(); err == nil {
		t.Errorf("test %d: expecting EOF error, nil received", len(tests)+1)
	} else if err != io.EOF {
		t.Errorf("test %d: expecting EOF error, received %q", len(tests)+1, err)
	} else if token.typ != tokenDone {
		t.Errorf("test %d: expecting done token, received %s", len(tests)+1, token.typ)
	}

}
