
## gson game-json
*格式-内容分开的json*

json是一种非常好的配置方式，游戏中广为采纳。但在配置量非常大，格式相对单一的场景，自描述的结构显不够简洁，这也是很多游戏采用csv/excel配置的原因。
gson就是格式-内容的json，格式可以理解为csv的表头。同一张表中所有数据格式都根据表头来。既保留json强大的表达能力，又省略了重复、冗余的格式配置。


跟普通json的差别
- 格式、内容分离
- 省略“键、字符串类型值”的引号
- 支持特殊字符，显式加首尾的'"'允许包含特殊字符

缺点
- 不是非常灵活，内容要跟格式完全一致（后面内容可以省略，但是中间部分不行）


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
Test samples
---
* format: `{k_int:int}`
* content: `{999}`
* expect:  `{"k_int":999}`
---
* format: `{array_int:[int]}`
* content: `{  [1,2,3]}`
* expect:  `{"array_int":[1,2,3]}`
---
* format: `{k_str:string,array_int:[int]}`
* content: `{hemaolong,   [1,2,3]}`
* expect:  `{"k_str":"hemaolong","array_int":[1,2,3]}`
---
* format: `{map_array:[{x:int,y:float, z:double}]}`
* content: `{[{11,22,33},{55,66,77}]}`
* expect:  `{"map_array":[{"x":11,"y":22,"z":33},{"x":55,"y":66,"z":77}]}`
---
* format: `{array_str:[string]}`
* content: `{  [1,2,3]}`
* expect:  `{"array_str":["1","2","3"]}`
---
* format: `{k_str:string,coopCardID:int, array_int:[int], map_array:[{x:int,y:float, z:double}], k_int:int}`
* content: `{hemaolong, 1024, [1,2,3], [{11,22,33},{55,66,77}], 999}`
* expect:  `{"k_int":999,"coopCardID":1024,"k_str":"hemaolong","array_int":[1,2,3],"map_array":[{"x":11,"y":22,"z":33},{"x":55,"y":66,"z":77}]}`


---  
* format: `[[string]]`
* content: `[[1,6,8],[2]]`
* expect:  `[["1","6","8"],["2"]]`

