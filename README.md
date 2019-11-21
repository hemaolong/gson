
## gson game-json
*格式-内容分开的json*

json非常适合游戏配置，但是面对大量配置的情况，自描述的结构非常啰嗦。
特别是大量表格类配置，格式通常一致。格式-内容分开非常必要。


### 例子
---
```
package main

import "fmt"
import "github.com/hemaolong/gson"

func main() {
	formatStr := "{k_str:string,array_int:[int]}"
	contentStr := "{hemaolong, [1,2,3]}"

	encoder, _ := gson.NewEncoder([]byte(formatStr))
	output, _ := encoder.Marshal([]byte(contentStr))
	fmt.Println("output|", string(output))

	// {"k_str":"hemaolong","array_int":[1,2,3]}
}

```

---
* format: `{k_int:int}`,
* content: `{999}`,
* expect:  `{"k_int":999}`,
---
* format: `{array_int:[int]}`,
* content: `{  [1,2,3]}`,
* expect:  `{"array_int":[1,2,3]}`,
---
* format: `{k_str:string,array_int:[int]}`,
* content: `{hemaolong,   [1,2,3]}`,
* expect:  `{"k_str":"hemaolong","array_int":[1,2,3]}`,
---
* format: `{map_array:[{x:int,y:float, z:double}]}`,
* content: `{[{11,22,33},{55,66,77}]}`,
* expect:  `{"map_array":[{"x":11,"y":22,"z":33},{"x":55,"y":66,"z":77}]}`,
---
* format: `{array_str:[string]}`,
* content: `{  [1,2,3]}`,
* expect:  `{"array_str":["1","2","3"]}`,
---
* format: `{k_str:string,coopCardID:int, array_int:[int], map_array:[{x:int,y:float, z:double}], k_int:int}`,
* content: `{hemaolong, 1024, [1,2,3], [{11,22,33},{55,66,77}], 999}`,
* expect:  `{"k_int":999,"coopCardID":1024,"k_str":"hemaolong","array_int":[1,2,3],"map_array":[{"x":11,"y":22,"z":33},{"x":55,"y":66,"z":77}]}`,

