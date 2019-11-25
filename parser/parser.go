/*
 * @Author: maolong.he@gmail.com
 * @Date: 2019-11-20 09:08:17
 * @Last Modified by: maolong.he@gmail.com
 * @Last Modified time: 2019-11-21 20:14:18
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

func Parse(format *laxer.Laxer) *Parser {
	format.InitFormat()
	return &Parser{format: format}
}

func (self *Parser) Parse(content *laxer.Laxer, buf *bytes.Buffer) {
	self.format.SetStackPos(0)
	self.content = content

	first := self.content.PeekToken()
	if first.Type == laxer.TokenString {
		_, self.lastError = self.parseMapPair(buf)
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

		// 数组类型的格式配置，只会配置一个元素的格式，默认所有格式相同。所以需要复位
		self.format.SetStackPos(initFormatPos)
		if contentT.Type == laxer.TokenString {
			formatT = self.format.PopFirst() // Get value type
			self.parsePrimitiveField(buf, formatT)
		} else {
			self.doParse(buf)
		}
		buf.WriteByte(',')
	}
	buf.Truncate(buf.Len() - 1)

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

	buf.WriteByte('{')
	for {
		if self.Finished() {
			return
		}
		contentT := self.content.PeekToken()
		if contentT.Type == laxer.TokenMapEnd {
			break
		}

		isDefaultValue := false
		isDefaultValue, self.lastError = self.parseMapPair(buf)
		if !isDefaultValue {
			buf.WriteByte(',')
		}
	}
	// truncate the last ','
	buf.Truncate(buf.Len() - 1)
	self.format.PopFirst()
	// fmt.Println("----map pop format--", ttt)
	self.content.PopFirst()
	buf.WriteByte('}')

}

// return hasData, error
func (self *Parser) parseMapPair(buf *bytes.Buffer) (bool, error) {
	if self.lastError != nil {
		return false, self.lastError
	}

	keyToken := self.format.PopFirst()
	if keyToken == nil {
		// self.format.IncrTokenPos(-1)
		return false, fmt.Errorf("expect key, but found EOF")
	}
	if keyToken.Type == laxer.TokenString {
		prePos := buf.Len()
		buf.WriteString(strconv.Quote(keyToken.Value))
		buf.WriteByte(':')
		// self.parseMapValue(buf, keyToken)
		if self.lastError != nil {
			return false, self.lastError
		}

		if len(keyToken.Ultra) != 0 {
			isDefault := self.parsePrimitiveField(buf, keyToken)
			if isDefault {
				buf.Truncate(prePos)
			}
			return isDefault, nil
		}
		self.doParse(buf)
	}
	return false, self.lastError // fmt.Errorf("expect key, map field expect string, but found:%v", keyToken)
}

// return if is default value
// 返回是否是默认值，调用者可以根据返回值确定是否省略字段
func (self *Parser) parsePrimitiveField(buf *bytes.Buffer, typeToken *laxer.Token) (isDefaultValue bool) {
	isDefaultValue = false
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
		lv := len(contentToken.Value)
		if lv > 0 {
			if contentToken.Value[0] == '"' && contentToken.Value[lv-1] == '"' {
				// 自带引号分割的字符串
				buf.WriteString(contentToken.Value)
			} else {
				buf.WriteString(strconv.Quote(contentToken.Value))
			}
		}
		// buf.WriteString(strconv.Quote(contentToken.Value))
	} else if typeToken.Ultra == "bool" {
		if len(contentToken.Value) == 0 || contentToken.Value == "0" || contentToken.Value == "false" {
			buf.WriteString("false")
			isDefaultValue = true
		} else {
			buf.WriteString("true")
		}
	} else {
		buf.WriteString(contentToken.Value)
	}
	return
}
