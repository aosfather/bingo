package bingo

import (
	"encoding/json"
	"github.com/aosfather/bingo_dao"
	"github.com/aosfather/bingo_mvc"
	utils "github.com/aosfather/bingo_utils"
	"io/ioutil"
	"os"
	"strings"
)

type ApplicationContext struct {
	config     map[string]string
	logfactory *utils.LogFactory
	services   InjectMan
	holder     ValuesHolder
	ds         *bingo_dao.DataSource
	mvc        *bingo_mvc.HttpDispatcher
}

func (this *ApplicationContext) shutdown() {

	//关闭所有service

	this.logfactory.Close()

}

func (this *ApplicationContext) completedLoaded() {
	this.mvc.Run()
}
func (this *ApplicationContext) GetLog(module string) utils.Log {
	return this.logfactory.GetLog(module)
}

func (this *ApplicationContext) GetConnection() *bingo_dao.Connection {
	if this.ds != nil {
		return this.ds.GetConnection()
	}
	return nil
}

func (this *ApplicationContext) CreateDao() *bingo_dao.BaseDao {
	dao := bingo_dao.BaseDao{}
	dao.Init(this.ds)
	return &dao
}

//不能获取bingo自身的属性，只能获取应用自身的扩展属性
func (this *ApplicationContext) GetPropertyFromConfig(key string) string {
	if strings.HasPrefix(key, "bingo.") {
		return ""
	}
	return this.getProperty(key)
}

func (this *ApplicationContext) RegisterService(name string, service interface{}) {
	if name != "" && service != nil {
		instance := this.services.GetObjectByName(name)
		if instance == nil {
			this.holder.ProcessValueTag(service)
			this.services.AddObjectByName(name, service)
		}
	}

}

func (this *ApplicationContext) GetService(name string) interface{} {
	if name != "" {
		//return this.services[name]
		return this.services.GetObjectByName(name)
	}
	return nil
}

func (this *ApplicationContext) getProperty(key string) string {
	if this.config == nil {
		return ""
	}
	return this.config[key]
}

func (this *ApplicationContext) init(file string) {
	if file != "" && utils.IsFileExist(file) {
		f, err := os.Open(file)
		if err == nil {
			txt, _ := ioutil.ReadAll(f)
			json.Unmarshal(txt, &this.config)
		}

	}
	if this.config == nil {
		this.config = make(map[string]string)
	}
	//this.services=make(map[string]interface{})
	this.services.Init(nil)
	this.services.AddObject(this)
	this.holder.InitByFunction(this.GetPropertyFromConfig)
	this.initLogFactory()
	this.initSessionFactory()
	this.mvc = new(bingo_mvc.HttpDispatcher)
	this.mvc.SetLog(this.logfactory.GetLog("mvc"))
	//this.mvc.SetPort(this.context.getProperty(""))
	//this.mvc.SetRoot()
	this.mvc.Init()
}

func (this *ApplicationContext) initSessionFactory() {
	if this.config["bingo.system.usedb"] == "true" {
		this.logfactory.Write(utils.LEVEL_INFO, "bingo", "init db")

		var ds bingo_dao.DataSource
		bingo_dao.SetLogger(this.logfactory.GetLog("dao"))
		ds.DBtype = this.config["bingo.db.type"]
		ds.DBname = this.config["bingo.db.name"]
		ds.DBurl = this.config["bingo.db.url"]
		ds.DBuser = this.config["bingo.db.user"]
		ds.DBpassword = this.config["bingo.db.password"]
		ds.Init()

		this.ds = &ds
	}
}

func (this *ApplicationContext) initLogFactory() {
	this.logfactory = &utils.LogFactory{}
	this.logfactory.SetConfig(utils.LogConfig{true, this.config["bingo.log.file"]})
}

func (this *ApplicationContext) initServices() {

}

func (this *ApplicationContext) AddControl(c bingo_mvc.HttpController) {
	this.services.InjectObject(c)
	this.mvc.AddController(c)
}
