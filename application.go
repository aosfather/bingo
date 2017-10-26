package bingo

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Application struct {
	config  map[string]string
	router  defaultRouter
	factory *SessionFactory
	port    int
}

func (this *Application) GetSession() *TxSession {
	if this.factory != nil {
		return this.factory.GetSession()
	}
	return nil
}

func (this *Application) getProperty(key string) string {
	if this.config == nil {
		return ""
	}
	return this.config[key]
}

func (this *Application) AddHandler(url string, handler HttpMethodHandler) {
	var rule routerRule
	rule.url = url
	rule.methodHandler = handler
	this.router.addRouter(&rule)
}

func (this *Application) init() {
	if this.config == nil {
		this.config = make(map[string]string)
	}
	if this.config["bingo.system.usedb"] == "true" {
		log.Println("init db")

		var sqlfactory SessionFactory
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
	this.router.setTemplateRoot(this.config["bingo.mvc.template"])
	//设置静态处理
	this.router.staticHandler = &staticController{staticDir: this.config["bingo.mvc.static"]}
}

func (this *Application) Run(file string) {
	if file != "" && isFileExist(file) {
		f, err := os.Open(file)
		if err == nil {
			log.Println("open config file " + file)
			txt, _ := ioutil.ReadAll(f)
			json.Unmarshal(txt, &this.config)
		}

	}

	this.init()
	http.ListenAndServe("localhost:"+strconv.Itoa(this.port), &this.router)
}
