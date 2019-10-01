package main

import (
//	"unicode"
)

var (
	out       string
	isComment bool
)

// scanner
func scan(in string) {
	p := in

	for len(p) != 0 {
		c := rune(p[0])

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
				line_comment(p)
				p = ""
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
				out += string(p[0:1])
				p = p[2:]
				continue
			}
		}

		// newline
		// white space
		// block comment
		// character literal

		out += string(p[0])
		p = p[1:]
	}
err:
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

func line_comment(s string) {
	out += "TK_COMMENT_S"
	out += s
	out += "TK_COMMENT_E"
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
