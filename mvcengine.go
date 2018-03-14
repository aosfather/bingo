package bingo

import (
	"github.com/aosfather/bingo/mvc"
	"strconv"
	"net/http"
	"github.com/aosfather/bingo/utils"
)

type MvcEngine struct {
	router  mvc.DefaultRouter
	port    int
	context *ApplicationContext
	server *http.Server
}

func (this *MvcEngine) Init(context *ApplicationContext) {
	this.context=context
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

func (this *MvcEngine)AddController(c mvc.HttpController){
    if c!=nil {
    	c.SetBeanFactory(this.context)
    	this.context.holder.ProcessValueTag(c)
    	this.context.services.InjectObject(c)
    	c.Init()
    	url:=c.GetUrl()
    	if url=="" {
    		url="/"+utils.GetRealType(c).Name()
		}
    	this.AddHandler(url,c.(mvc.HttpMethodHandler))
	}

}

func (this *MvcEngine) AddInterceptor(h mvc.CustomHandlerInterceptor) {
	if h != nil {
		this.router.AddInterceptor(h)
	}

}

func (this *MvcEngine) run(){
	if this.server!=nil {
		return
	}
	this.server= &http.Server{Addr: ":"+strconv.Itoa(this.port), Handler: &this.router}
	this.server.ListenAndServe()
}

func (this *MvcEngine)shutdown(){
	if this.server!=nil {
		this.server.Shutdown(nil)
	}
}


//加载handler
type OnLoadHandler func(mvc *MvcEngine,context *ApplicationContext) bool

type rootController struct {
	mvc.SimpleController
}

func (this *rootController) GetUrl()string {
	return "/"
}
func (this *rootController) Get(c mvc.Context, p interface{}) (interface{}, mvc.BingoError){
	return "<h1>hello bingo!</h1>",nil
}