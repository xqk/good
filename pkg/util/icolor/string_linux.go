//go:build linux
// +build linux

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
	return sprint(YellowColor, msg, arg...)
}

// Red 红色
func Red(msg string, arg ...interface{}) string {
	return sprint(RedColor, msg, arg...)
}

// Blue 蓝色
func Blue(msg string, arg ...interface{}) string {
	return sprint(BlueColor, msg, arg...)
}

// Green 绿色
func Green(msg string, arg ...interface{}) string {
	return sprint(GreenColor, msg, arg...)
}

// sprint 格式化
func sprint(colorValue int, msg string, arg ...interface{}) string {
	if arg != nil {
		return fmt.Sprintf("\x1b[%dm%s\x1b[0m %+v", colorValue, msg, arrToTransform(arg))
	} else {
		return fmt.Sprintf("\x1b[%dm%s\x1b[0m", colorValue, msg)
	}
}
