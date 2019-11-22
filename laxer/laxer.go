/*
 * @Author: maolong.he@gmail.com
 * @Date: 2019-11-20 09:00:20
 * @Last Modified by: maolong.he@gmail.com
 * @Last Modified time: 2019-11-21 18:52:06
 */

package laxer

import (
	"fmt"
	"strings"
)

type TokenType byte

const (
	TokenString TokenType = 's'
	// TokenKey    TokenType = 'k'
	// TokenEOF    TokenType = 'e'
	// TokenValue  TokenType = 'v'

	TokenArrayBegin TokenType = '['
	TokenArrayEnd   TokenType = ']'

	TokenMapBegin TokenType = '{'
	TokenMapEnd   TokenType = '}'

	TokenComma TokenType = ','
	// TokenColon TokenType = ':'

	TokenError TokenType = 'E'
)

type Token struct {
	Type  TokenType
	Value string

	// 作为一个类型Token，表示类型的Type：string int？
	Ultra string
}

type Laxer struct {
	input   string
	nextPos int

	stackPos int
	tokens   []*Token
}

type LexFn func(*Laxer) LexFn

func (self *TokenType) String() string {
	return fmt.Sprintf("%c", *self)
}

func (self *Token) String() string {
	if self.Type == TokenString {
		if len(self.Ultra) > 0 {
			return fmt.Sprintf("{'%c' '%s' '%s'}", self.Type, self.Value, self.Ultra)
		} else {
			return fmt.Sprintf("{'%c' '%s'}", self.Type, self.Value)
		}
	}
	if self.Type == TokenError {
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

// 只有格式 laxer需要调用，将带':'的字符串类型拆分成两个Token
func (self *Laxer) InitFormat() {
	self.stackPos = 0

	for _, v := range self.tokens {
		if v.Type == TokenString {
			tmp := strings.Split(v.Value, ":")
			v.Value = tmp[0]
			if len(tmp) >= 2 {
				v.Ultra = tmp[1]
			} else {
				// fmt.Println(tmp)
				// 数组类型才会缺失，补齐
				v.Ultra = tmp[0]
			}
		}
	}
}

func (self *Laxer) LastError() *Token {
	l := len(self.tokens)
	if l == 0 {
		return nil
	}

	t := self.tokens[l-1]
	if t.Type == TokenError {
		return t
	}
	return nil
}

func (self *Laxer) peekLaxStrPos() string {
	start := self.nextPos
	if start < 0 {
		start = 0
	}
	end := start + 20
	if end > len(self.input) {
		end = len(self.input)
	}
	return self.input[start:end]
}

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

func (self *Laxer) GetStackPos() int {
	return self.stackPos
}

func (self *Laxer) SetStackPos(i int) {
	self.stackPos = i
}

func (self *Laxer) push(t *Token) {
	// fmt.Println("push|", t)
	self.tokens = append(self.tokens, t)
}

func (self *Laxer) pushError(err error) {
	self.tokens = append(self.tokens, &Token{Type: TokenError, Value: err.Error()})
}

// func unscape(s string) string {
// 	if len(s) == 0 {
// 		return s
// 	}
// 	buf := make([]byte, 0, len(s))
// 	isEscape := false
// 	for i := 0; i < len(s); i++ {
// 		c := s[i]
// 		if c == '\\' {
// 			isEscape = true
// 			continue
// 		}
// 		if isEscape {
// 			isEscape = false
// 		}
// 		buf = append(buf, c)
// 	}
// 	s = string(buf)
// 	return s
// }

func (self *Laxer) genToken() *Token {
	ls := len(self.input)
	if self.nextPos >= ls {
		return nil
	}

	preCharIsEscape := false
	duringQuotedString := false
	for k := self.nextPos; k < ls; k++ {
		if preCharIsEscape {
			preCharIsEscape = false
			continue
		}
		char := byte(self.input[k])
		isQuote := char == '"'
		if !duringQuotedString {
			duringQuotedString = isQuote
		} else {
			if isQuote {
				duringQuotedString = false
			}
		}

		preCharIsEscape = isEscape(char)
		if duringQuotedString {
			continue
		}

		if !isPlain(TokenType(char)) {
			if k > self.nextPos {
				// 前面有剩余未解析字符，需要弹出内容Token
				scrToken := strings.TrimSpace(self.input[self.nextPos:k])
				if len(scrToken) > 0 {
					t := &Token{Type: TokenString, Value: scrToken}
					self.nextPos = k
					return t
				}
			}
			// fmt.Println("i am char ", fmt.Sprintf("%c", char))
			t := &Token{Type: TokenType(char)}
			self.nextPos = k + 1
			return t
		}
	}
	return nil
}

func isEscape(c byte) bool {
	return c == '\\'
}

func isPlain(c TokenType) bool {
	return c != TokenArrayBegin && c != TokenArrayEnd &&
		c != TokenMapBegin && c != TokenMapEnd &&
		c != TokenComma
}

func (self *Laxer) run() {
	self.laxBegin(self.genToken())
}

func (self *Laxer) laxBegin(cur *Token) {
	self.push(cur)
	switch cur.Type {
	case TokenArrayBegin:
		self.laxArray()

	case TokenMapBegin:
		self.laxMap()
	case TokenString:
		// fmt.Println("nothing todo")
	default:
		self.pushError(fmt.Errorf(`invalid json format '%s'...`, self.peekLaxStrPos()))
	}
	// fmt.Println("push token|", t)
}

func (self *Laxer) laxArray() {
	for {
		nextEle := self.genToken()
		if nextEle == nil {
			self.pushError(fmt.Errorf(`expect array element but found EOF '%s'...`, self.peekLaxStrPos()))
			return
		}

		if nextEle.Type == TokenArrayEnd {
			self.push(nextEle)
			return
		}

		self.laxBegin(nextEle)
		nextEle = self.genToken()
		if self.LastError() != nil {
			return
		}
		if nextEle == nil {
			self.pushError(fmt.Errorf(`expect array end but found EOF '%s'...`, self.peekLaxStrPos()))
			return
		}
		// self.push(nextEle)
		if nextEle.Type == TokenArrayEnd {
			self.push(nextEle)
			return
		}
		if nextEle.Type != TokenComma {
			self.pushError(fmt.Errorf(`expect array splititor:"%c", but found "%v"`, TokenComma, nextEle))
			return
		}
	}
}

func (self *Laxer) laxMap() {
	nextEle := self.genToken()
	for {
		if nextEle == nil || self.LastError() != nil {
			return
		}
		// 剩余字符应该是',' or '?'
		if nextEle.Type == TokenMapEnd {
			self.push(nextEle)
			return
		}

		self.laxBegin(nextEle)
		nextEle = self.genToken()
		if self.LastError() != nil {
			return
		}
		if nextEle == nil {
			self.pushError(fmt.Errorf(`expect map end but found EOF '%s'...`, self.peekLaxStrPos()))
			return
		}
		if nextEle.Type == TokenMapEnd {
			self.push(nextEle)
			return
		}
		if nextEle.Type == TokenComma {
			nextEle = self.genToken()
			continue
		}
	}
}

// 取出一个单词
// func (self *Laxer) laxToken() {
// 	for {
// 		c, ok := self.nextSkipEscape()
// 		if !ok {
// 			break
// 		}

// 		tc := TokenType(c)
// 		if !isPlain(tc) {
// 			self.push(TokenString)
// 			self.ignore()

// 			if tc == TokenColon {
// 				continue
// 			}

// 			self.incrPos(-1)
// 			break
// 		}
// 	}
// }
