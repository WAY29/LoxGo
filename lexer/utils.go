package lexer

func isDigit(r rune) bool {
	return ('0' <= r && r <= '9')
}

func isAlpha(r rune) bool {
	return ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z')
}

func isAlphaNumeric(r rune) bool {
	return isAlpha(r) || isDigit(r)
}
