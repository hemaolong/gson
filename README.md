
## gson game-json
*格式-内容分开的json*

json是一种非常好的配置方式，有大量的使用场景。
游戏中广为采纳。但是如果配置量非常大，自描述的结构显不够简洁。
例如大量表格类配置，格式通常一致。格式-内容分开就显得非常必要。
既便于阅读，又方便“手工填写”。


跟普通json的差别
- 格式、内容分离
- 字符串类型的引号省略
- 支持特殊字符，用'\'进行转义



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

