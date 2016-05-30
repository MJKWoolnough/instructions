package instructions

import (
	"errors"
	"io"

	"github.com/MJKWoolnough/parser"
)

type stateFn func() (token, stateFn)

type tokenType uint8

const (
	tokenError tokenType = iota
	tokenFunction
	tokenNumber
	tokenString
	tokenComment
	tokenDone
)

func (t tokenType) String() string {
	switch t {
	case tokenError:
		return "error"
	case tokenFunction:
		return "function"
	case tokenNumber:
		return "number"
	case tokenString:
		return "string"
	case tokenComment:
		return "comment"
	case tokenDone:
		return "done"
	default:
		return "unknown"
	}
}

type token struct {
	typ  tokenType
	data string
}

type lexer struct {
	p     parser.Parser
	state stateFn
	err   error
}

func newLexer(r io.Reader) *lexer {
	l := &lexer{
		p: parser.NewReaderParser(r),
	}
	l.state = l.lexFunction
	return l
}

func (l *lexer) GetToken() (token, error) {
	if l.err == io.EOF {
		return token{tokenDone, ""}, l.err
	}
	var t token
	t, l.state = l.state()
	//l.p.Get()
	if l.err == io.EOF {
		if t.typ == tokenError {
			l.err = io.ErrUnexpectedEOF
		}
	}
	return t, l.err
}

func (l *lexer) lexFunction() (token, stateFn) {
	if l.p.Peek() == -1 {
		l.err = io.EOF
		return token{tokenDone, "EOF"}, l.errorFn
	}
	l.p.Get()
	if l.p.Accept("#") {
		return l.lexComment()
	}
	l.p.ExceptRun(" \n")
	return token{tokenFunction, l.p.Get()}, l.lexArgument
}

func (l *lexer) lexComment() (token, stateFn) {
	l.p.Get()
	l.p.ExceptRun("\n")
	c := l.p.Get()
	l.p.Accept("\n")
	l.p.Get()
	return token{tokenComment, c}, l.lexFunction
}

func (l *lexer) lexArgument() (token, stateFn) {
	l.p.AcceptRun(" ")
	char := l.p.Peek()
	switch char {
	case '"':
		l.p.Accept(string(char))
		return l.lexString()
	case '\n':
		l.p.Accept(string(char))
		return l.lexFunction()
	default:
		return l.lexNumber()
	}
}

func (l *lexer) lexNumber() (token, stateFn) {
	l.p.Get()
	l.p.Accept("+-")
	digits := "0123456789"
	if l.p.Accept("0") && l.p.Accept("xX") {
		digits = "0123456789abcdefABCDEF"
	}
	l.p.AcceptRun(digits)
	if l.p.Accept(".") {
		l.p.AcceptRun(digits)
	}
	if l.p.Accept("eE") {
		l.p.Accept("+-")
		l.p.AcceptRun(digits)
	}
	l.p.Accept("i")
	n := l.p.Get()
	if len(n) == 0 {
		l.err = ErrInvalidNumber
		return l.errorFn()
	}
	return token{tokenNumber, n}, l.lexArgument
}

func (l *lexer) lexString() (token, stateFn) {
	l.p.Get()
Loop:
	for {
		l.p.ExceptRun("\\\"")
		switch l.p.Peek() {
		case '\\':
			if !l.p.Accept("\\\"abtnfr") {
				l.err = ErrInvalidEscape
				return l.errorFn()
			}
		case '"':
			break Loop
		default:
			l.err = io.ErrUnexpectedEOF
			return token{tokenError, l.err.Error()}, l.errorFn
		}
	}
	s := l.p.Get()
	l.p.Accept("\"")
	l.p.Get()
	return token{tokenString, s}, l.lexArgument
}

func (l *lexer) errorFn() (token, stateFn) {
	return token{tokenError, l.err.Error()}, l.errorFn
}

// Errors
var (
	ErrInvalidEscape = errors.New("invalid escape character")
	ErrInvalidNumber = errors.New("invalid number format")
)
