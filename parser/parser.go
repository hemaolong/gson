/*
 * @Author: maolong.he@gmail.com
 * @Date: 2019-11-20 09:08:17
 * @Last Modified by: maolong.he@gmail.com
 * @Last Modified time: 2019-11-21 18:50:40
 */

package parser

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/hemaolong/gson/laxer"
)

type GsonType byte

// const (
// 	WordMap    GsonType = 'm'
// 	WordArray  GsonType = 'a'
// 	WordString GsonType = 's'
// )

func (self *GsonType) String() string {
	return fmt.Sprintf("%c", *self)
}

type gsonMap map[string]interface{}
type gsonArray []interface{}
type gsonPrimitive string

type GsonObject struct {
	Type  GsonType
	Value interface{}
}

type Parser struct {
	format  *laxer.Laxer
	content *laxer.Laxer

	lastError error
}

func Parse(format *laxer.Laxer, content *laxer.Laxer) *Parser {
	return &Parser{format: format, content: content}
}

func (self *Parser) Parse(buf *bytes.Buffer) {
	self.format.InitFormat()
	first := self.content.PeekToken()
	if first.Type == laxer.TokenString {
		self.lastError = self.parseMapPair(buf)
	} else {
		self.doParse(buf)
	}
}

func (self *Parser) LastError() error {
	return self.lastError
}

func (self *Parser) Finished() bool {
	return self.lastError != nil || self.content.PeekToken() == nil
}

func (self *Parser) doParse(buf *bytes.Buffer) {
	if self.Finished() {
		return
	}

	v := self.content.PopFirst()
	if v == nil {
		self.lastError = fmt.Errorf("not vailid end")
		return
	}

	switch v.Type {
	// case laxer.TokenEOF:
	// 	break

	case laxer.TokenArrayBegin:
		self.content.IncrTokenPos(-1)
		self.parseArray(buf)
	case laxer.TokenMapBegin:
		self.content.IncrTokenPos(-1)
		self.parseMap(buf)

	default:
		self.lastError = fmt.Errorf("unsupoorted token:%v", v)
	}
}

func (self *Parser) parseArray(buf *bytes.Buffer) {
	self.content.PopFirst() // Pop '['

	formatT := self.format.PopFirst() // Pop '['
	initFormatPos := self.format.GetStackPos()

	needComma := false
	buf.WriteByte('[')
	for {
		if self.Finished() {
			return
		}
		contentT := self.content.PeekToken()
		if contentT.Type == laxer.TokenArrayEnd {
			self.content.PopFirst()
			break
		}
		if needComma {
			buf.WriteByte(',')
		}

		// 数组类型的格式配置，只会配置一个元素的格式，默认所有格式相同。所以需要复位
		self.format.SetStackPos(initFormatPos)
		if contentT.Type == laxer.TokenString {
			formatT = self.format.PopFirst() // Get value type
			self.parsePrimitiveField(buf, formatT)
		} else {
			self.doParse(buf)
		}

		needComma = true
	}
	self.format.PopFirst() // Pop ']'
	// fmt.Println("pop array", initFormatPos, ttt)
	buf.WriteByte(']')
}

func (self *Parser) parseMap(buf *bytes.Buffer) {
	formatT := self.format.PopFirst()
	contentT := self.content.PopFirst()
	if formatT.Type != contentT.Type || formatT.Type != laxer.TokenMapBegin {
		self.lastError = fmt.Errorf("invalid map field. format:%v content:%v", formatT, contentT)
		return
	}

	needComma := false
	buf.WriteByte('{')
	for {
		if self.Finished() {
			return
		}
		contentT := self.content.PeekToken()
		if contentT.Type == laxer.TokenMapEnd {
			break
		}

		if needComma {
			buf.WriteByte(',')
		}
		self.lastError = self.parseMapPair(buf)
		needComma = true
	}
	self.format.PopFirst()
	// fmt.Println("----map pop format--", ttt)
	self.content.PopFirst()
	buf.WriteByte('}')

}

func (self *Parser) parseMapPair(buf *bytes.Buffer) error {
	if self.lastError != nil {
		return self.lastError
	}

	keyToken := self.format.PopFirst()
	if keyToken == nil {
		// self.format.IncrTokenPos(-1)
		return fmt.Errorf("expect key, but found EOF")
	}
	if keyToken.Type == laxer.TokenString {
		buf.WriteString(strconv.Quote(keyToken.Value))
		buf.WriteByte(':')
		self.parseMapValue(buf, keyToken)
	}
	return self.lastError // fmt.Errorf("expect key, map field expect string, but found:%v", keyToken)
}

func (self *Parser) parsePrimitiveField(buf *bytes.Buffer, typeToken *laxer.Token) {
	if self.lastError != nil {
		return
	}

	// 以value为准，允许value省略后续字段
	contentToken := self.content.PopFirst()
	if typeToken.Type != laxer.TokenString || contentToken.Type != laxer.TokenString {
		if typeToken.Type != laxer.TokenString {
			self.lastError = fmt.Errorf("miss array value type :%v content:%v", typeToken, contentToken)
			return
		}
	}

	if typeToken.Ultra == "string" {
		buf.WriteString(strconv.Quote(contentToken.Value))
	} else {
		buf.WriteString(contentToken.Value)
	}
}

// map-value 字段类型跟值一一对应
func (self *Parser) parseMapValue(buf *bytes.Buffer, t *laxer.Token) {
	if self.lastError != nil {
		return
	}

	if len(t.Ultra) != 0 {
		self.parsePrimitiveField(buf, t)
		return
	}
	self.doParse(buf)
}
