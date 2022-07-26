// Code generated by "stringer -type=TokenType"; DO NOT EDIT.

package lexer

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[TokenNone-0]
	_ = x[LEFT_PAREN-1]
	_ = x[RIGHT_PAREN-2]
	_ = x[LEFT_BRACE-3]
	_ = x[RIGHT_BRACE-4]
	_ = x[LEFT_BRACKET-5]
	_ = x[RIGHT_BRACKET-6]
	_ = x[COMMA-7]
	_ = x[DOT-8]
	_ = x[QUESTION-9]
	_ = x[COLON-10]
	_ = x[MINUS-11]
	_ = x[PLUS-12]
	_ = x[SEMICOLON-13]
	_ = x[SLASH-14]
	_ = x[STAR-15]
	_ = x[BANG-16]
	_ = x[BANG_EQUAL-17]
	_ = x[EQUAL-18]
	_ = x[EQUAL_EQUAL-19]
	_ = x[GREATER-20]
	_ = x[GREATER_EQUAL-21]
	_ = x[LESS-22]
	_ = x[LESS_EQUAL-23]
	_ = x[PLUSPLUS-24]
	_ = x[MINUSMINUS-25]
	_ = x[IDENTIFIER-26]
	_ = x[STRING-27]
	_ = x[NUMBER-28]
	_ = x[AND-29]
	_ = x[CLASS-30]
	_ = x[ELSE-31]
	_ = x[FALSE-32]
	_ = x[FUN-33]
	_ = x[FOR-34]
	_ = x[IF-35]
	_ = x[NIL-36]
	_ = x[OR-37]
	_ = x[PRINT-38]
	_ = x[RETURN-39]
	_ = x[SUPER-40]
	_ = x[THIS-41]
	_ = x[TRUE-42]
	_ = x[VAR-43]
	_ = x[WHILE-44]
	_ = x[BREAK-45]
	_ = x[CONTINUE-46]
	_ = x[EOF-47]
}

const _TokenType_name = "TokenNoneLEFT_PARENRIGHT_PARENLEFT_BRACERIGHT_BRACELEFT_BRACKETRIGHT_BRACKETCOMMADOTQUESTIONCOLONMINUSPLUSSEMICOLONSLASHSTARBANGBANG_EQUALEQUALEQUAL_EQUALGREATERGREATER_EQUALLESSLESS_EQUALPLUSPLUSMINUSMINUSIDENTIFIERSTRINGNUMBERANDCLASSELSEFALSEFUNFORIFNILORPRINTRETURNSUPERTHISTRUEVARWHILEBREAKCONTINUEEOF"

var _TokenType_index = [...]uint16{0, 9, 19, 30, 40, 51, 63, 76, 81, 84, 92, 97, 102, 106, 115, 120, 124, 128, 138, 143, 154, 161, 174, 178, 188, 196, 206, 216, 222, 228, 231, 236, 240, 245, 248, 251, 253, 256, 258, 263, 269, 274, 278, 282, 285, 290, 295, 303, 306}

func (i TokenType) String() string {
	if i >= TokenType(len(_TokenType_index)-1) {
		return "TokenType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _TokenType_name[_TokenType_index[i]:_TokenType_index[i+1]]
}
