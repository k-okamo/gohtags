package main

import (
	"strings"
	"unicode"
)

type scanStr struct {
	in             string
	out            string
	isString       bool
	isCharacter    bool
	isRawstring    bool
	isBlockcomment bool
	isLinecomment  bool
	isSeparator    bool
}

var reserved = []string{
	"break", "case", "chan", "const", "continue",
	"default", "defer", "else", "fallthrough", "for",
	"func", "go", "goto", "if", "import",
	"interface", "map", "package", "range", "return",
	"select", "struct", "switch", "type", "var",
}
var primtypes = []string{
	"bool", "byte", "complex64", "complex128", "error", "float32", "float64",
	"int", "int8", "int16", "int32", "int64", "rune", "string",
	"uint", "uint8", "uint16", "uint32", "uint64", "uintptr",
}
var constants = []string{
	"true", "false", "iota", "nil",
}

var builtinfuncs = []string{
	"append", "cap", "close", "complex", "copy", "delete", "imag",
	"len", "make", "new", "panic", "print", "println", "real", "recover",
}

var delimiters = []string{
	"+", "&", "+=", "&=", "&&", "==", "!=", "(", ")",
	"-", "|", "-=", "|=", "||", "<", "<=", "[", "]",
	"*", "^", "*=", "^=", "<-", ">", ">=", "{", "}",
	"/", "<<", "/=", "<<=", "++", "=", ":=", ",", ";",
	"%", ">>", "%=", ">>=", "--", "!", "...", ".", ":",
	"&^", "&^=",
}

func scan(s *scanStr) {

	if s.isBlockcomment {
		s.out += "<em class='comment'>"
	}
	if s.isRawstring {
		s.out += "<em class='string'>"
	}

ROOT:
	for len(s.in) > 0 {
		c := ([]rune(s.in))[0]
		switch c {
		case '&': // encode
			encodeAmp(s)
			continue
		case '<': // encode
			encodeLt(s)
			continue
		case '>': // encode
			encodeGt(s)
			continue
		case '\\': // encode
			encodeBslash(s)
			continue
		case '\t': // encode
			encodeTab(s)
			continue
		case '\'': // encode, character literal
			if s.isString || s.isRawstring || s.isBlockcomment || s.isLinecomment {
				encodeSquot(s)
			} else if s.isCharacter {
				characterLiteral_end(s)
			} else {
				characterLiteral_start(s)
			}
			continue
		case '`': // encode, raw strubg
			if s.isString || s.isBlockcomment || s.isLinecomment {
				encodeBquot(s)
			} else if s.isRawstring {
				rawstringLiteral_end(s)
			} else {
				rawstringLiteral_start(s)
			}
			continue
		case '"': // encode, string literal
			if s.isCharacter || s.isRawstring || s.isBlockcomment || s.isLinecomment {
				encodeDquot(s)
			} else if s.isString {
				stringLiteral_end(s)
			} else {
				stringLiteral_start(s)
			}
			continue

		case '/': // block comment(start), line comment
			if len(s.in) > 1 {
				if ([]rune(s.in))[1] == '*' && !s.isLinecomment && !s.isString {
					blockComment_start(s)
					continue
				}
				if ([]rune(s.in))[1] == '/' && !s.isLinecomment && !s.isBlockcomment && !s.isString {
					lineComment_start(s)
					continue
				}
			}
			break

		case '*': // block comment(end)
			if len(s.in) > 1 {
				if ([]rune(s.in))[1] == '/' && !s.isLinecomment && !s.isString {
					blockComment_end(s)
					continue
				}
			}
			break

		default:
		}

		isPlaincode := !s.isBlockcomment && !s.isLinecomment && !s.isString && !s.isRawstring
		// reserved word, primary types
		for _, key := range append(reserved, primtypes...) {
			if strings.HasPrefix(s.in, key) && isPlaincode && s.isSeparator {
				l := len(key)
				if s.in == key || unicode.IsSpace(([]rune(s.in))[l]) || isDelimiter(string(([]rune(s.in))[l:])) {
					reserved_start(s, key)
					continue ROOT
				}
			}
		}
		// builtin functions
		for _, key := range builtinfuncs {
			if strings.HasPrefix(s.in, key) && isPlaincode && s.isSeparator {
				l := len(key)
				if s.in == key || rune(s.in[l]) == '(' {
					reserved_start(s, key)
					continue ROOT
				}
			}
		}
		// constants (true, false, nil, iota)
		for _, key := range constants {
			if strings.HasPrefix(s.in, key) && isPlaincode && s.isSeparator {
				l := len(key)
				if s.in == key || unicode.IsSpace(([]rune(s.in))[l]) || isDelimiter(string(([]rune(s.in))[l:])) {
					constants_start(s, key)
					continue ROOT
				}
			}
		}

		s.isSeparator = unicode.IsSpace(c) || isDelimiter(s.in)
		s.out += string(c)
		s.in = string(([]rune(s.in))[1:])
	}

	if s.isBlockcomment || s.isLinecomment || s.isRawstring {
		s.out += "</em>"
		s.isLinecomment = false
	}
}

func encodeAmp(s *scanStr) {
	s.out += "&amp;"
	s.in = string(([]rune(s.in))[1:])
}

func encodeLt(s *scanStr) {
	s.out += "&lt;"
	s.in = string(([]rune(s.in))[1:])
}

func encodeGt(s *scanStr) {
	s.out += "&gt;"
	s.in = string(([]rune(s.in))[1:])
}

func encodeSquot(s *scanStr) {
	s.out += "&#039;"
	s.in = string(([]rune(s.in))[1:])
}

func encodeDquot(s *scanStr) {
	s.out += "&quot;"
	s.in = string(([]rune(s.in))[1:])
}

func encodeBquot(s *scanStr) {
	s.out += "&#906;"
	s.in = string(([]rune(s.in))[1:])
}

func encodeBslash(s *scanStr) {
	if len(s.in) > 1 {
		s.out += string(([]rune(s.in))[0:2])
		s.in = string(([]rune(s.in))[2:])
	} else if len(s.in) == 1 { // "\"
		s.out += string(([]rune(s.in))[0:1])
		s.in = string(([]rune(s.in))[1:])
	}
}

func encodeTab(s *scanStr) {
	s.out += "    "
	s.in = string(([]rune(s.in))[1:])
}

func characterLiteral_start(s *scanStr) {
	s.isCharacter = true
	s.out += "<em class='string'>&#039;"
	s.in = string(([]rune(s.in))[1:])
}

func characterLiteral_end(s *scanStr) {
	s.isCharacter = false
	s.out += "&#039;</em>"
	s.in = string(([]rune(s.in))[1:])
}

func stringLiteral_start(s *scanStr) {
	s.isString = true
	s.out += "<em class='string'>&quot;"
	s.in = string(([]rune(s.in))[1:])
}

func stringLiteral_end(s *scanStr) {
	s.isString = false
	s.out += "&quot;</em>"
	s.in = string(([]rune(s.in))[1:])
}

func rawstringLiteral_start(s *scanStr) {
	s.isRawstring = true
	s.out += "<em class='string'>&#096;"
	s.in = string(([]rune(s.in))[1:])
}

func rawstringLiteral_end(s *scanStr) {
	s.isRawstring = false
	s.out += "&#096;</em>"
	s.in = string(([]rune(s.in))[1:])
}

func lineComment_start(s *scanStr) {
	s.isLinecomment = true
	s.out += "<em class='comment'>" + string(([]rune(s.in))[:2])
	s.in = string(([]rune(s.in))[2:])
}

func blockComment_start(s *scanStr) {
	s.isBlockcomment = true
	s.out += "<em class='comment'>" + string(([]rune(s.in))[:2])
	s.in = string(([]rune(s.in))[2:])
}

func blockComment_end(s *scanStr) {
	s.isBlockcomment = false
	s.out += string(([]rune(s.in))[:2]) + "</em>"
	s.in = string(([]rune(s.in))[2:])
}

func reserved_start(s *scanStr, key string) {
	s.out += "<strong class='reserved'>" + key + "</strong>"
	s.in = string(([]rune(s.in))[len(key):])
}

func constants_start(s *scanStr, key string) {
	s.out += "<strong class='constants'>" + key + "</strong>"
	s.in = string(([]rune(s.in))[len(key):])
}

func isDelimiter(s string) bool {
	for _, d := range delimiters {
		if len(s) >= len(s) && strings.HasPrefix(s, d) {
			return true
		}
	}
	return false
}
