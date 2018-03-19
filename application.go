package bingo

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"github.com/aosfather/bingo/utils"
	"github.com/aosfather/bingo/sql"
	"github.com/aosfather/bingo/mvc"

	"syscall"
	"os/signal"
)


type Application struct {
	Name    string
	config  map[string]string
	router  mvc.DefaultRouter
	factory *sql.SessionFactory
	port    int
	logfactory *utils.LogFactory
}

func (this *Application) Validate(obj interface{}) []mvc.BingoError {

	return this.router.Validate(obj)

}

func (this *Application) AddInterceptor(h mvc.CustomHandlerInterceptor) {
	if h != nil {
		this.router.AddInterceptor(h)
	}

}

func (this *Application) GetSession() *sql.TxSession {
	if this.factory != nil {
		return this.factory.GetSession()
	}
	return nil
}

func (this *Application)GetLog(module string)utils.Log {
	return this.logfactory.GetLog(module)
}

//不能获取bingo自身的属性，只能获取应用自身的扩展属性
func (this *Application) GetPropertyFromConfig(key string) string {
	if strings.HasPrefix(key, "bingo.") {
		return ""
	}
	return this.getProperty(key)
}

func (this *Application) getProperty(key string) string {
	if this.config == nil {
		return ""
	}
	return this.config[key]
}

func (this *Application) AddHandler(url string, handler mvc.HttpMethodHandler) {
	var rule mvc.RouterRule
	rule.Init(url,handler)
	//rule.url = url
	//rule.methodHandler = handler
	this.router.AddRouter(&rule)
}

func (this *Application) init() {
	if this.config == nil {
		this.config = make(map[string]string)
	}
	if this.config["bingo.system.usedb"] == "true" {
		this.logfactory.Write(utils.LEVEL_INFO,"bingo","init db")

		var sqlfactory sql.SessionFactory
		sqlfactory.DBtype = this.config["bingo.db.type"]
		sqlfactory.DBname = this.config["bingo.db.name"]
		sqlfactory.DBurl = this.config["bingo.db.url"]
		sqlfactory.DBuser = this.config["bingo.db.user"]
		sqlfactory.DBpassword = this.config["bingo.db.password"]
		sqlfactory.Init()

		this.factory = &sqlfactory
	}
	this.router.Init(this.factory)
	this.port = 8990
	if this.config["bingo.system.port"] != "" {
		port, err := strconv.Atoi(this.config["bingo.system.port"])
		if err == nil {
			this.port = port
		}

	}
	this.router.SetTemplateRoot(this.config["bingo.mvc.template"])
	//设置静态处理
	this.router.SetStaticControl(this.config["bingo.mvc.static"],this.logfactory.GetLog("bingo.static"))

}

func (this *Application) Run() {
    defer this.logfactory.Close()
	this.init()
	http.ListenAndServe(":"+strconv.Itoa(this.port), &this.router)
}

func (this *Application) Load(file string) {
	if file != "" && utils.IsFileExist(file) {
		f, err := os.Open(file)
		if err == nil {
			txt, _ := ioutil.ReadAll(f)
			json.Unmarshal(txt, &this.config)
		}

	}
	this.logfactory=&utils.LogFactory{}
	this.logfactory.SetConfig(utils.LogConfig{true,this.config["bingo.log.file"]})
	this.router.SetLog(this.logfactory.GetLog("bingo.router"))
}


type TApplication struct {
	Name    string
	context ApplicationContext
	mvc MvcEngine
	onload OnLoad
	loadHandler OnLoadHandler
	onShutdown OnDestoryHandler
}
func (this *TApplication)SetHandler(load OnLoad,handler OnLoadHandler){
	this.onload=load
	this.loadHandler=handler
}

func (this *TApplication)SetOnDestoryHandler(h OnDestoryHandler){
	this.onShutdown=h
}

func (this *TApplication)Run(file string){
	go this.signalListen()
	this.context.init(file)
	//加载factory
    if this.onload!=nil {
    	if !this.onload(&this.context){
    		panic("load service error! please check onload function")

		}
	}
	this.context.services.Inject()
	this.mvc.Init(&this.context)
	//加载controller
	this.mvc.AddController(&rootController{})
	if this.loadHandler!=nil {
		if !this.loadHandler(&this.mvc,&this.context){
			panic("load http handler error! please check OnLoadHandler function")
		}
	}

	this.mvc.run()

}

//监听被kill的信号，当被kill的时候执行处理
func (this *TApplication) signalListen() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGUSR1, syscall.SIGUSR2)
	//for {
		s := <-c
		//收到信号后的处理，这里只是输出信号内容，可以做一些更有意思的事
		this.context.GetLog("bingo").Info("get signal:%s", s)
		this.processShutdown()

	//}
}

func (this *TApplication) processShutdown(){
    //处理关闭操作
    this.mvc.shutdown()
    //关闭service
    this.context.shutdown()
    //处理自定义关闭操作
	if this.onShutdown!=nil {
		this.onShutdown()
	}

	os.Exit(0)

}