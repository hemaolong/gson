/*
 * @Author: maolong.he@gmail.com
 * @Date: 2019-11-20 09:00:20
 * @Last Modified by: maolong.he@gmail.com
 * @Last Modified time: 2019-11-21 10:20:37
 */

package laxer

import (
	"fmt"
	"strings"
)

type TokenType byte

const (
	TokenString TokenType = 's'
	TokenKey    TokenType = 'k'
	// TokenEOF    TokenType = 'e'
	// TokenValue  TokenType = 'v'

	TokenArrayBegin TokenType = '['
	TokenArrayEnd   TokenType = ']'

	TokenMapBegin TokenType = '{'
	TokenMapEnd   TokenType = '}'

	TokenSplitor TokenType = ','
	TokenColon   TokenType = ':'

	TokenError TokenType = 'E'
)

type Token struct {
	Type  TokenType
	Value string
}

type Laxer struct {
	input       string
	start       int
	pos         int
	preIsEscape bool

	stackPos int
	tokens   []*Token
}

type LexFn func(*Laxer) LexFn

func (self *TokenType) String() string {
	return fmt.Sprintf("%c", *self)
}

func (self *Token) String() string {
	if self.Type == TokenString {
		return fmt.Sprintf("{'%c' '%s'}", self.Type, self.Value)
	}
	return fmt.Sprintf("{'%c'}", self.Type)
}

func Lax(s string) *Laxer {
	result := &Laxer{input: s, tokens: make([]*Token, 0, 200)}
	result.run()
	return result
}

func (self *Laxer) String() string {
	return fmt.Sprintf("%v", self.tokens)
}

func (self *Laxer) LastError() *Token {
	l := len(self.tokens)
	if l == 0 {
		return &Token{TokenError, "empty token list"}
	}

	t := self.tokens[l-1]
	if t.Type == TokenError {
		return t
	}
	return nil
}

// func (self *Laxer) PeekFirst() *Token {
// 	l := len(self.tokens)
// 	if l == 0 {
// 		return nil
// 	}

// 	t := self.tokens[0]
// 	return t
// }

func (self *Laxer) PopFirst() *Token {
	l := len(self.tokens)
	if self.stackPos >= l || self.stackPos < 0 {
		return nil
	}

	self.stackPos++
	return self.tokens[self.stackPos-1]
}

func (self *Laxer) PeekToken() *Token {
	l := len(self.tokens)
	if self.stackPos >= l || self.stackPos < 0 {
		return nil
	}

	return self.tokens[self.stackPos]
}

func (self *Laxer) IncrTokenPos(v int) {
	self.stackPos += v
}

func (self *Laxer) push(t TokenType) {
	if t == TokenString {
		l, ok := self.left()
		if ok {
			trimed := strings.TrimSpace(l)
			if len(trimed) > 0 {
				v := &Token{Type: t, Value: strings.TrimSpace(l)}
				// fmt.Println("push token|", v)
				self.tokens = append(self.tokens, v)
			}
		}
	} else {
		self.tokens = append(self.tokens, &Token{Type: t})
	}
}

func (self *Laxer) pushWithValue(t TokenType, v string) {
	self.tokens = append(self.tokens, &Token{Type: t, Value: v})
}

// 有待解析的字符串是否以s开头
func (self *Laxer) leftHasPrefix(prefix string) bool {
	left := self.input[self.pos:]
	return strings.HasPrefix(left, prefix)
}

// 去掉前导空白字符
func (self *Laxer) incrPos(v int) {
	self.pos += v
}

func (self *Laxer) ignore() {
	self.start = self.pos
}

func (self *Laxer) next() (byte, bool) {
	if self.pos >= len(self.input) {
		return 0, false
	}

	self.pos++
	return self.input[self.pos-1], true
}

func (self *Laxer) nextSkipEscape() (byte, bool) {
	c, ok := self.next()
	if self.preIsEscape {
		self.preIsEscape = false
		c, ok = self.next()
	} else {
		if isEscape(c) {
			self.preIsEscape = true
			c, ok = self.next()
		}
	}

	return c, ok
}

func (self *Laxer) left() (string, bool) {
	if self.pos > self.start+1 && self.pos < len(self.input) {
		s := self.input[self.start : self.pos-1]
		self.start = self.pos
		return s, true
	}
	return "", false
}

func (self *Laxer) peekToLast() string {
	return self.input[self.pos:]
}

// func (self *Laxer) peek() (byte, bool) {
// 	if self.pos >= len(self.input) {
// 		return 0, false
// 	}

// 	return self.input[self.pos], true
// }

func isEscape(c byte) bool {
	return c == '\\'
}

func isPlain(c TokenType) bool {
	return c != TokenArrayBegin && c != TokenArrayEnd &&
		c != TokenMapBegin && c != TokenMapEnd &&
		c != TokenSplitor && c != TokenColon
}

func (self *Laxer) run() {
	self.laxBegin()
	// l.pushWithValue(TokenEOF, "EOF")
}

func (self *Laxer) laxBegin() {
	c, ok := self.nextSkipEscape()
	if !ok {
		return
	}
	self.incrPos(-1)

	switch TokenType(c) {
	case TokenArrayBegin:
		self.laxArray()

	case TokenMapBegin:
		self.laxMap()

	case TokenSplitor:
		self.incrPos(1)
		self.laxToken()
	default:
		// self.incrPos(-1)
		self.laxToken()
	}
}

func (self *Laxer) laxArray() {
	self.next() // skip '['
	self.ignore()
	self.push(TokenArrayBegin)
	for {
		self.laxBegin()

		c, ok := self.nextSkipEscape()
		if !ok {
			break
		}
		if TokenType(c) == TokenArrayEnd {
			// fmt.Println("->>>>>>", self.tokens)
			break
		}

	}
	self.push(TokenArrayEnd)

}

func (self *Laxer) laxMap() {
	self.next() // skip '{'
	self.ignore()
	self.push(TokenMapBegin)
	for {
		//
		self.laxBegin()

		c, ok := self.nextSkipEscape()
		if !ok {
			break
		}
		if TokenType(c) == TokenMapEnd {
			break
		}
		self.incrPos(-1)
	}
	self.push(TokenMapEnd)

	// c, ok := l.next()
	// if !ok {
	// 	l.pushWithValue(TokenError, "map value expected but found EOF")
	// }

	// if TokenType(c) == TokenMapEnd {
	// 	l.push(TokenMapEnd)
	// 	break
	// }
}

// func laxToken(l *Laxer) {
// 	for {
// 		c, ok := l.nextSkipEscape()
// 		if !ok {
// 			l.pushWithValue(TokenEOF, "expect value but found EOF")
// 			return
// 		}

// 		if !isPlain(TokenType(c)) {
// 			l.incrPos(-1)
// 			fmt.Println("to end|", l.peekToLast())

// 			l.push(TokenValue)
// 			l.incrPos(1)
// 			laxBegin(l)
// 			return
// 		}
// 	}
// }

// 取出一个单词
func (self *Laxer) laxToken() {
	for {
		c, ok := self.nextSkipEscape()
		if !ok {
			break
		}

		tc := TokenType(c)
		if !isPlain(tc) {
			self.push(TokenString)
			self.ignore()

			if tc == TokenColon {
				continue
			}

			self.incrPos(-1)
			break
		}
	}
}
