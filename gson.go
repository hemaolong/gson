/*
 * @Author: maolong.he@gmail.com
 * @Date: 2019-11-21 20:04:18
 * @Last Modified by: maolong.he@gmail.com
 * @Last Modified time: 2019-11-21 20:37:36
 */
package gson

import (
	"bytes"
	"fmt"

	"github.com/hemaolong/gson/laxer"
	"github.com/hemaolong/gson/parser"
)

type Encoder struct {
	fl *laxer.Laxer
	p  *parser.Parser
}

func NewEncoder(formatStr []byte) (*Encoder, error) {
	fl := laxer.Lax(string(formatStr))
	if fl.LastError() != nil {
		return nil, fmt.Errorf("gson parse format error:%v", fl.LastError())
	}
	result := &Encoder{}
	result.fl = fl
	result.p = parser.Parse(fl)

	return result, nil
}

func (self *Encoder) Marshal(content []byte) ([]byte, error) {
	cl := laxer.Lax(string(content))
	if cl.LastError() != nil {
		return nil, fmt.Errorf("gson parse content error:%v", cl.LastError())
	}
	buf := bytes.Buffer{}
	self.p.Parse(cl, &buf)
	err := self.p.LastError()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
