package bingo

import (
	utils "github.com/aosfather/bingo_utils"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

type Application struct {
	Name       string
	context    ApplicationContext
	onload     OnLoad
	onShutdown OnDestoryHandler
}

func (this *Application) SetHandler(load OnLoad) {
	this.onload = load
}

func (this *Application) SetOnDestoryHandler(h OnDestoryHandler) {
	this.onShutdown = h
}

func (this *Application) RunApp() {
	var configfile string
	if len(os.Args) > 1 {
		configfile = os.Args[1]
	} else {
		configfile = "config.conf"
	}

	if !utils.IsFileExist(configfile) {
		dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		configfile = dir + "/" + configfile
	}

	this.Run(configfile)

}

func (this *Application) Run(file string) {
	go this.signalListen()
	this.context.init(file)
	//加载factory
	if this.onload != nil {
		if !this.onload(&this.context) {
			panic("load service error! please check onload function")

		}
	}
	this.context.services.Inject()
	this.mvc.Init(&this.context)
	//加载controller
	this.mvc.AddController(&rootController{})
	if this.loadHandler != nil {
		if !this.loadHandler(&this.mvc, &this.context) {
			panic("load http handler error! please check OnLoadHandler function")
		}
	}

	this.mvc.run()

}

const (
	SIGUSR1 = syscall.Signal(0x1e)
	SIGUSR2 = syscall.Signal(0x1f)
)

//监听被kill的信号，当被kill的时候执行处理
func (this *Application) signalListen() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM, SIGUSR1, SIGUSR2)
	//for {
	s := <-c
	//收到信号后的处理，这里只是输出信号内容，可以做一些更有意思的事
	this.context.GetLog("bingo").Info("get signal:%s", s)
	this.processShutdown()

	//}
}

func (this *Application) processShutdown() {
	//处理关闭操作
	this.mvc.shutdown()
	//关闭service
	this.context.shutdown()
	//处理自定义关闭操作
	if this.onShutdown != nil {
		this.onShutdown()
	}

	os.Exit(0)

}

//load函数，如果加载成功返回true，否则返回FALSE
type OnLoad func(context *ApplicationContext) bool
type OnDestoryHandler func() bool //shutdown的handler，用于处理关闭服务的自定义动作
