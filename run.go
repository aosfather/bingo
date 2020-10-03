package main

import (
	"flag"
	"github.com/aosfather/bingo_mvc/context"
	"github.com/aosfather/bingo_mvc/fasthttp"
	"github.com/aosfather/bingo_utils/files"
)

var _app *Application
var _login *LoginAccess = &LoginAccess{}

func main() {
	//应用目录
	var appdirctory = flag.String("run", "app", "Input Your application dirctory path")
	flag.Parse()
	if appdirctory != nil {
		if files.IsFileExist(*appdirctory) {
			_app = &Application{Root: *appdirctory}
			start()
		}
	}

}

func start() {
	boot := context.Boot{}
	dispatch := &fasthttp.FastHTTPDispatcher{}
	dispatch.AddInterceptor(_login)
	boot.Init(dispatch, load)
	//boot.Init(&http.HttpDispatcher{}, load)
	boot.StartByConfigFile(_app.GetFilePath("bingo.yaml"))
}

func load() []interface{} {
	return []interface{}{&System{}, &MenuTree{}, _app, _login, &FormActions{}, &DefaultLogin{}, &Desktop{}}
}
