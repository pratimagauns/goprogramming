package main

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"text/scanner"
)

type ruleList struct {
	lineRules, textRules []func([]byte) []byte
}

// filter rule definition reader
func (maker *ruleList) add(input io.Reader, name string) (err error) {
	// tokeniser
	tok := makeTokeniser(input, func(t *scanner.Scanner, msg string) {
		err = fmt.Errorf("rule definition in \"%s\", line %d: %s", name, t.Line, msg)
	})

	fail := func(msg string) {
		tok.Error(tok, msg)
	}

	failInvalidToken := func(msg string, t rune) {
		fail(fmt.Sprintf("Expected %s, but found %s", msg, strconv.Quote(tok.TokenText())))
	}

	// Rule spec format
	//	scope type `match` `replacement`
	// where
	//	scope: 'line' | 'text'
	//	type:	'word' | 'regex'

	// parser
	for t := skipNewLines(tok); t != scanner.EOF; {
		// scope
		if t != scanner.Ident {
			failInvalidToken("rule scope", t)
			return
		}

		ruleScope := tok.TokenText()

		// rule type
		if t = tok.Scan(); t != scanner.Ident {
			failInvalidToken("rule type", t)
			return
		}

		ruleType := tok.TokenText()

		// regex or word
		if t = tok.Scan(); t != scanner.String {
			failInvalidToken("match string", t)
			return
		}

		var match string

		if match, err = strconv.Unquote(tok.TokenText()); err != nil {
			fail("Match string: " + err.Error())
			return
		}

		if len(match) == 0 {
			fail("Match string cannot be empty")
			return
		}

		// substitution
		if t = tok.Scan(); t != scanner.String {
			failInvalidToken("substitution string", t)
			return
		}

		var subst string

		if subst, err = strconv.Unquote(tok.TokenText()); err != nil {
			fail("Invalid substitution string: " + err.Error())
			return
		}

		// create filter function
		var ruleFunc func([]byte) []byte

		switch ruleType {
		case "word":
			ruleFunc = makeWordRule([]byte(match), []byte(subst))
		case "regex":
			re, e := regexp.Compile(match)
			if e != nil {
				fail(e.Error())
				return
			}
			ruleFunc = makeRegexRule(re, []byte(subst))
		default:
			fail("Unknown rule type: " + ruleType)
			return
		}

		switch ruleScope {
		case "line":
			maker.lineRules = append(maker.lineRules, ruleFunc)
		case "text":
			maker.textRules = append(maker.textRules, ruleFunc)
		default:
			fail("Unknown rule scope: " + ruleScope)
			return
		}

		// newline or EOF
		switch t = tok.Scan(); t {
		case scanner.EOF:
			// nothing to do
		case '\n':
			t = skipNewLines(tok)
		default:
			failInvalidToken("newline", t)
			return
		}
	}

	return
}

func makeTokeniser(input io.Reader, errFunc func(*scanner.Scanner, string)) *scanner.Scanner {
	tok := new(scanner.Scanner).Init(input)

	tok.Mode = scanner.SkipComments | scanner.ScanComments | scanner.ScanIdents | scanner.ScanStrings | scanner.ScanRawStrings
	tok.Whitespace = 1<<'\t' | 1<<'\r' | 1<<' '
	tok.Error = errFunc
	return tok
}

func skipNewLines(tok *scanner.Scanner) (t rune) {
	for t = tok.Scan(); t == '\n'; t = tok.Scan() {
		// empty
	}

	return
}

func makeWordRule(match, subst []byte) func([]byte) []byte {
	return func(s []byte) []byte {
		return bytes.Replace(s, match, subst, -1)
	}
}

func makeRegexRule(re *regexp.Regexp, subst []byte) func([]byte) []byte {
	return func(s []byte) []byte {
		return re.ReplaceAll(s, subst)
	}
}
