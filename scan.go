package main

import (
	"unicode"
)

var (
	out       string
	isComment bool
)

var keywords = []string{
	"break", "case", "chan", "const", "continue",
	"default", "defer", "else", "fallthrough", "for",
	"func", "go", "goto", "if", "import",
	"interface", "map", "package", "range", "return",
	"select", "struct", "switch", "type", "var",
}
var primtypes = []string{
	"uint8", "uint16", "uint32", "uint64", "uint", "uintptr", "byte",
	"int8", "int16", "int32", "int64", "int", "rune",
	"float32", "float64", "string",
}

// scanner
func scan(s string) {
	p := s

ROOT:
	for len(p) != 0 {
		c := rune(p[0])

		// Whitespace
		if unicode.IsSpace(c) {
			p = whitespace(p)
			continue
		}

		// block comment
		if isComment {
			p = block_comment(p)
			continue
		}
		if c == '/' {
			if len(p) > 1 && rune(p[1]) == '*' {
				isComment = true
				p = block_comment(p)
				continue
			}
		}

		// line comment
		if c == '/' {
			if len(p) > 1 && rune(p[1]) == '/' {
				p = line_comment(p)
				continue
			}
		}

		// Character literal
		if c == '\'' {
			if len(p) > 2 {
				p = character_literal(p)
				continue
			}
		}

		// String literal
		if c == '"' {
			p = string_literal(p)
			continue
		}

		// \"
		if c == '\\' {
			if len(p) < 2 {
				goto err
			}
			if rune(p[1]) == '"' {
				out += string(p[:1])
				p = p[2:]
				continue
			}
		}

		// keywords and primitive types
		for _, key := range append(keywords, primtypes...) {
			if len(p) < len(key) {
				continue
			}
			if p == key {
				p = keyword(p, key)
				break ROOT
			}
			l := len(key)
			if p[:l] == key && unicode.IsSpace(rune(p[l])) {
				p = keyword(p, key)
				break
			}
		}

		out += string(p[0])
		p = p[1:]
	}
err:
}

func keyword(s, key string) string {
	out += "TK_KEYWORD_S"
	out += s[:len(key)]
	out += "TK_KEYWORD_E"
	if len(s) == len(key) {
		return ""
	}
	return s[len(key):]
}

func whitespace(s string) string {
	out += s[:1]
	return s[1:]
}

func character_literal(s string) string {

	out += "TK_STRING_S"
	out += s[:3]
	out += "TK_STRING_E"
	s = s[3:]
	return s
}

func string_literal(s string) string {

	// s starts from double quotation(").
	out += "TK_STRING_S"
	out += s[:1]
	s = s[1:]

	for rune(s[0]) != '"' {
		if len(s) == 0 {
			goto err
		}
		if rune(s[0]) == '\\' {
			out += string(s[0:2])
			s = s[2:]
			continue
		}
		out += string(s[0])
		s = s[1:]
		if len(s) == 0 {
			goto err
		}
	}
	out += string(s[0])
	out += "TK_STRING_E"
err:
	return s[1:]
}

func line_comment(s string) string {
	out += "TK_COMMENT_S"
	out += s
	out += "TK_COMMENT_E"
	return ""
}

func block_comment(s string) string {
	p := s
	out += "TK_COMMENT_S"

	for len(p) > 0 {
		c := rune(p[0])
		if c == '*' {
			if len(p) > 1 && rune(p[1]) == '/' {
				out += string(p[0:2])
				p = p[2:]
				out += "TK_COMMENT_E"
				isComment = false
				break
			}
		}
		out += string(c)
		p = p[1:]
	}
	if isComment {
		out += "TK_COMMENT_E"
	}
	return string(p)
}
