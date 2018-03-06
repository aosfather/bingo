package bingo

import (
	"github.com/aosfather/bingo/mvc"
	"strconv"
)

type MvcEngine struct {
	router  mvc.DefaultRouter
	port    int
}

func (this *MvcEngine) Init(context *ApplicationContext) {
	this.router.SetLog(context.GetLog("bingo.router"))
	this.router.Init(context.factory)
	this.port = 8990
	if context.getProperty("bingo.system.port") != "" {
		port, err := strconv.Atoi(context.getProperty("bingo.system.port"))
		if err == nil {
			this.port = port
		}

	}
	this.router.SetTemplateRoot(context.getProperty("bingo.mvc.template"))
	//设置静态处理
	this.router.SetStaticControl(context.getProperty("bingo.mvc.static"),context.GetLog("bingo.static"))

}

func (this *MvcEngine) AddHandler(url string, handler mvc.HttpMethodHandler){
	var rule mvc.RouterRule
	rule.Init(url,handler)
	//rule.url = url
	//rule.methodHandler = handler
	this.router.AddRouter(&rule)
}


//加载handler
type OnLoadHandler func(mvc *MvcEngine,context *ApplicationContext) bool