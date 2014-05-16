package instructions

import (
	"github.com/MJKWoolnough/memio"
	"github.com/MJKWoolnough/tokeniser"
	"io"
	"testing"
)

func TestTokens(t *testing.T) {
	str := []byte("Hello\nWorld 1, 2, 3\n#Comment 3 5 7 1\nBeep \"This is a String\", 5.4\n#Comment 1\nShouldError")
	tokens := tokeniser.New(memio.Open(str), lexFunction)
	tests := []struct {
		tokenType  tokeniser.ItemType
		tokenValue string
	}{
		{FUNCTION, "Hello"},
		{FUNCTION, "World"},
		{NUMBER, "1"},
		{NUMBER, "2"},
		{NUMBER, "3"},
		{COMMENT, "Comment 3 5 7 1"},
		{FUNCTION, "Beep"},
		{STRING, "This is a String"},
		{NUMBER, "5.4"},
		{COMMENT, "Comment 1"},
	}
	for n, test := range tests {
		token, err := tokens.Next()
		if err != nil {
			t.Errorf("test %d: received unexpected error: %s", n+1, err)
		} else if token == nil {
			t.Errorf("test %d: received unexpected nil token", n+1)
		} else if test.tokenType != token.Typ {
			t.Errorf("test %d: expecting token type %d, got %d", n+1, test.tokenType, token.Typ)
		} else if test.tokenValue != token.Val {
			t.Errorf("test %d: expecting token value %q, got %q", n+1, test.tokenValue, token.Val)
		}
	}
	if token, err := tokens.Next(); err == nil {
		t.Errorf("test 9: expecting EOF error, nil received")
	} else if err != io.EOF {
		t.Errorf("test 9: expecting EOF error, received %q", err)
	} else if token != nil {
		t.Errorf("test 9: expecting nil token, received %v", token)
	}

}
