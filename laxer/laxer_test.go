/*
 * @Author: maolong.he@gmail.com
 * @Date: 2019-11-20 11:21:09
 * @Last Modified by: maolong.he@gmail.com
 * @Last Modified time: 2019-11-21 18:43:06
 */

package laxer

import (
	"fmt"
	"testing"
)

func TestLaxerRun(t *testing.T) {
	{
		//
		formatS := "{coopName:string,coopCardID:int, intArray:[int],multi_array:[{x:int,y:float, z:double}],bornDate:int}"
		l := Lax(formatS)
		if l.LastError() != nil {
			panic(l.LastError())
		}
		l.InitFormat()
		fmt.Println(l)
	}
	{
		contentS := "{hemaolong, 1024, [1,2,3], 999}"
		l := Lax(contentS)
		if l.LastError() != nil {
			panic(l.LastError())
		}
		fmt.Println(l)
	}

}
