package main

import (
	"github.com/vulppine/fotoDen/tool"
)

func main() {
	err := tool.ParseCmd()
	if err != nil {
		panic(err)
	}
}
