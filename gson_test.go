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
			expect:  `{"k_str":"hemaolong","coopCardID":1024,"array_int":[1,2,3],"map_array":[{"x":11,"y":22,"z":33},{"x":55,"y":66,"z":77}],"k_int":999}`,
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

		// escapes
		testInput{format: `[[string]]`,
			content: `[["heml1,-,",6,8],[2]]`,
			expect:  `[["heml1,-,","6","8"],["2"]]`,
		},
		testInput{format: `[[string]]`,
			content: `[["heml1,-\n,\"",6,8],[2]]`,
			expect:  `[["heml1,-\n,\"","6","8"],["2"]]`,
		},

		// ellipsis empty fields
		testInput{format: `{array_str:[string],born_date:int}`,
			content: `{  [1024]}`,
			expect:  `{"array_str":["1024"]}`,
		},

		// bool - false(default)
		testInput{format: `{array_str:[string],is_ok:bool}`,
			content: `{  [1024], 0}`,
			expect:  `{"array_str":["1024"]}`,
		},
		testInput{format: `{array_str:[string],is_ok:bool}`,
			content: `{  [1024], false}`,
			expect:  `{"array_str":["1024"]}`,
		},
		// bool-true
		testInput{format: `{array_str:[string],is_ok:bool}`,
			content: `{  [1024], true}`,
			expect:  `{"array_str":["1024"],"is_ok":true}`,
		},
		testInput{format: `{array_str:[string],is_ok:bool}`,
			content: `{  [1024], 121}`,
			expect:  `{"array_str":["1024"],"is_ok":true}`,
		},

		testInput{format: `{map:{{group:int,id:int,count:int}}}`,
			content: `{1,{101,102,103}}`,
			expect:  `{"1":{"group":101,"id":102,"count":103}}`,
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
			fmt.Println("output|", string(out))
			panic(fmt.Sprintf("marshal content error:%v", err))
		}

		_, err = simplejson.NewJson(out)
		if err != nil {
			fmt.Println("output|", string(out))
			panic(fmt.Sprintf("output not valid json, unmarshal error:%v", err))
		}

		fmt.Println("  real|", string(out))
		fmt.Println("  expect|", v.expect)

		if string(out) != v.expect {
			fmt.Println("len real|", len(out), int(out[0]))
			fmt.Println("len expect|", len(v.expect), int(v.expect[0]))

			panic("--------")
		}

		assert.Equal(t, out, []byte(v.expect))
	}
}
