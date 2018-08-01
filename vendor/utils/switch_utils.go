package utils

import (
	"fmt"
	"strconv"
)

// create by cwj
// 自动展开多个exp判断
// 当我们有多个判空条件做排列组合时，往往代码很难组织，多个if-else代码很难维护
// 本方法将以 2^n的方式对入参进行权重设置，最终返回权重值
// eg: (nil, nil) = 0
// eg: (nonil, nil) = 1
// eg: (nil, nonil) = 2
// eg: (nonil, nonil) = 3
// eg: (nil, nil, nonil) = 4
func NilCase(exp ...interface{}) (ret int){
	var l string = "正确条件：第"
	for i, e := range exp{
		if !IsEmpty(e){
			ret += 1 << uint(i)
			l += strconv.Itoa(i + 1) + ","
		}
	}
	if len(l) > 18 {
		fmt.Println(l[:len(l)-1] + "参")
	}
	return
}

func BoolCase(exp ...bool) (ret int){
	var l string = "正确条件：第"
	for i, e := range exp{
		if !e{
			ret += 1 << uint(i)
			l += strconv.Itoa(i + 1) + ","
		}
	}
	if len(l) > 18 {
		fmt.Println(l[:len(l)-1] + "参")
	}
	return
}
