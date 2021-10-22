package icolor

import "fmt"

const (
	RedColor = iota + 31
	GreenColor
	YellowColor
	BlueColor
)

// arrToTransform 数组转变为空格分隔的字符串
func arrToTransform(arg []interface{}) interface{} {
	var res interface{}

	for _, v := range arg {
		if res != nil {
			res = fmt.Sprintf("%v %v", res, v)
		} else {
			res = v
		}
	}

	return res
}
