package instructions

import "github.com/MJKWoolnough/tokeniser"

// Token Types
const (
	FUNCTION tokeniser.ItemType = iota
	COMMENT
	STRING
	NUMBER
)

func lexFunction(l *tokeniser.Lexer) tokeniser.StateFn {
	if l.Accept("#") {
		return lexComment
	}
	switch l.AcceptUntil(" \n") {
	case ' ':
		l.Backup()
		l.Emit(FUNCTION)
		l.Next()
		l.Clear()
		return lexArgument
	case '\n':
		l.Backup()
		l.Emit(FUNCTION)
		l.Next()
		l.Clear()
		return lexFunction
	default:
		l.Errorf("unexpected eof")
	}
	return nil
}

func lexComment(l *tokeniser.Lexer) tokeniser.StateFn {
	l.Clear()
	l.AcceptUntil("\n")
	l.Backup()
	l.Emit(COMMENT)
	l.Next()
	l.Clear()
	return lexFunction
}

func lexArgument(l *tokeniser.Lexer) tokeniser.StateFn {
	l.AcceptRun(" ")
	if l.Peek() == '"' {
		l.Next()
		l.Clear()
		lexString(l)
	} else {
		l.Clear()
		lexNumber(l)
	}
	l.Clear()
	l.AcceptRun(" ")
	if l.Accept(",") {
		return lexArgument
	} else if l.Accept("\n") {
		l.Clear()
		return lexFunction
	}
	return nil
}

func lexNumber(l *tokeniser.Lexer) tokeniser.StateFn {
	l.Accept("+-")
	digits := "0123456789"
	if l.Accept("0") && l.Accept("xX") {
		digits = "0123456789abcdefABCDEF"
	}
	l.AcceptRun(digits)
	if l.Accept(".") {
		l.AcceptRun(digits)
	}
	if l.Accept("eE") {
		l.Accept("+-")
		l.AcceptRun(digits)
	}
	l.Accept("i")
	l.Emit(NUMBER)
	return nil
}

func lexString(l *tokeniser.Lexer) tokeniser.StateFn {
	for l.AcceptUntil("\\\"") != '"' {
		if !l.Accept("\\\"abtnfr") {
			l.Errorf("invalid escape character")
			return nil
		}
	}
	l.Backup()
	l.Emit(STRING)
	l.Next()
	return nil
}
