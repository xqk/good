package icolor

import (
	"fmt"
	"testing"
)

func TestRandomColor(t *testing.T) {
	color := RandomColor()
	fmt.Println(color)
}

func TestBlue(t *testing.T) {
	fmt.Println(Blue("hello "))
	fmt.Println(Blue("hello ", "world"))
}

func TestGreen(t *testing.T) {
	fmt.Println(Green("hello "))
	fmt.Println(Green("hello ", "world"))
}

func TestRed(t *testing.T) {
	fmt.Println(Red("hello "))
	fmt.Println(Red("hello ", "world"))
}

func TestYellow(t *testing.T) {
	fmt.Println(Yellow("hello "))
	fmt.Println(Yellow("hello ", "world"))
}
