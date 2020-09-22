package main

import (
	"flag"
	"github.com/aosfather/bingo_mvc/context"
	"github.com/aosfather/bingo_mvc/fasthttp"
	"github.com/aosfather/bingo_utils/files"
)

func main() {
	//应用目录
	var appdirctory = flag.String("run", "app", "Input Your application dirctory path")
	flag.Parse()
	if appdirctory != nil {
		if files.IsFileExist(*appdirctory) {
			app := Application{Root: *appdirctory}
			start(&app)
		}
	}

}

func start(app *Application) {
	boot := context.Boot{}
	boot.Init(&fasthttp.FastHTTPDispatcher{}, load)
	//boot.Init(&http.HttpDispatcher{}, load)
	boot.StartByConfigFile(app.GetFilePath("bingo.yaml"))
}

func load() []interface{} {
	return []interface{}{&System{}}
}
