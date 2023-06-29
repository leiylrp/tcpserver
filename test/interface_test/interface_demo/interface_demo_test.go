package interface_demo

import (
	"fmt"
	"testing"
)

func TestInterfaceDemo(t *testing.T) {
	var people People = &Stu{}
	think := "love"
	fmt.Println(people.Speak(think))
}

func TestLive(t *testing.T) {
	if live() == nil {
		fmt.Println("AAAAAAAAA")
	} else {
		fmt.Println("BBBBBBBBB")
	}
}
