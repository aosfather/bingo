package main

import (
	"github.com/aosfather/bingo_mvc/context"
	"github.com/aosfather/bingo_mvc/fasthttp"
)

func main() {
	boot := context.Boot{}
	boot.Init(&fasthttp.FastHTTPDispatcher{}, load)
	boot.Start()
}
func load() []interface{} {

	return []interface{}{&System{}}
}
