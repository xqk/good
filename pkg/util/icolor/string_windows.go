//go:build windows
// +build windows

package icolor

import (
	"fmt"
	"math/rand"
	"strconv"
)

var _ = RandomColor()

// RandomColor 生成随机颜色
func RandomColor() string {
	return fmt.Sprintf("#%s", strconv.FormatInt(int64(rand.Intn(16777216)), 16))
}

// Yellow 黄色
func Yellow(msg string, arg ...interface{}) string {
	return sprint(msg, arg...)
}

// Red 红色
func Red(msg string, arg ...interface{}) string {
	return sprint(msg, arg...)
}

// Blue 蓝色
func Blue(msg string, arg ...interface{}) string {
	return sprint(msg, arg...)
}

// Green 绿色
func Green(msg string, arg ...interface{}) string {
	return sprint(msg, arg...)
}

// sprint 格式化
func sprint(msg string, arg ...interface{}) string {
	if arg != nil {
		return fmt.Sprintf("%s %+v\n", msg, arrToTransform(arg))
	} else {
		return fmt.Sprintf("%s", msg)
	}
}
