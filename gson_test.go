/*
 * @Author: maolong.he@gmail.com
 * @Date: 2019-11-20 14:29:57
 * @Last Modified by: maolong.he@gmail.com
 * @Last Modified time: 2019-11-21 20:22:56
 */

package gson

import (
	"fmt"
	"testing"

	"github.com/bitly/go-simplejson"
	"github.com/stretchr/testify/assert"
)

type testInput struct {
	format  string
	content string

	expect string
}

var (
	inputs = []testInput{
		testInput{format: `{k_int:int}`,
			content: `{999}`,
			expect:  `{"k_int":999}`,
		},
		testInput{format: `{array_int:[int]}`,
			content: `{  [1,2,3]}`,
			expect:  `{"array_int":[1,2,3]}`,
		},
		testInput{format: `{k_str:string,array_int:[int]}`,
			content: `{hemaolong,   [1,2,3]}`,
			expect:  `{"k_str":"hemaolong","array_int":[1,2,3]}`,
		},
		testInput{format: `{map_array:[{x:int,y:float, z:double}]}`,
			content: `{[{11,22,33},{55,66,77}]}`,
			expect:  `{"map_array":[{"x":11,"y":22,"z":33},{"x":55,"y":66,"z":77}]}`,
		},

		testInput{format: `{array_str:[string]}`,
			content: `{  [1,2,3]}`,
			expect:  `{"array_str":["1","2","3"]}`,
		},

		testInput{format: `{k_str:string,coopCardID:int, array_int:[int], map_array:[{x:int,y:float, z:double}], k_int:int}`,
			content: `{hemaolong, 1024, [1,2,3], [{11,22,33},{55,66,77}], 999}`,
			expect:  `{"k_int":999,"coopCardID":1024,"k_str":"hemaolong","array_int":[1,2,3],"map_array":[{"x":11,"y":22,"z":33},{"x":55,"y":66,"z":77}]}`,
		},

		// arrays
		testInput{format: `[[int]]`,
			content: `[[1,6,8],[2]]`,
			expect:  `[[1,6,8],[2]]`,
		},
		testInput{format: `[[string]]`,
			content: `[[1,6,8],[2]]`,
			expect:  `[["1","6","8"],["2"]]`,
		},
	}
)

func TestParser(t *testing.T) {

	for _, v := range inputs {
		encoder, err := NewEncoder([]byte(v.format))
		if err != nil {
			panic(fmt.Sprintf("lax format error:%v", err))
		}
		out, err := encoder.Marshal([]byte(v.content))
		if err != nil {
			panic(fmt.Sprintf("marshal content error:%v", err))
		}

		realJson, err := simplejson.NewJson(out)
		if err != nil {
			panic(fmt.Sprintf("output js not valid json, unmarshal error:%v", err))
		}
		realStr, err := realJson.MarshalJSON()
		if err != nil {
			panic(fmt.Sprintf("output js not valid json, marshal error:%v", err))
		}

		expectJson, err := simplejson.NewJson([]byte(v.expect))
		if err != nil {
			panic(fmt.Sprintf("expect js not valid json, unmarshal error:%v", err))
		}
		expectStr, err := expectJson.MarshalJSON()
		if err != nil {
			panic(fmt.Sprintf("expect js not valid json, marshal error:%v", err))
		}
		fmt.Println("expect|", string(expectStr))
		fmt.Println("  real|", string(realStr))

		assert.Equal(t, realStr, expectStr)
	}
}
