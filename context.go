package bingo

import (
	"strings"
	"github.com/aosfather/bingo/utils"
	"github.com/aosfather/bingo/sql"
	"os"
	"io/ioutil"
	"encoding/json"
)

type ApplicationContext struct {
	config  map[string]string
	logfactory *utils.LogFactory
	factory *sql.SessionFactory
	services map[string]interface{}
}

func (this *ApplicationContext)GetLog(module string)utils.Log {
	return this.logfactory.GetLog(module)
}

func (this *ApplicationContext) GetSession() *sql.TxSession {
	if this.factory != nil {
		return this.factory.GetSession()
	}
	return nil
}
//不能获取bingo自身的属性，只能获取应用自身的扩展属性
func (this *ApplicationContext) GetPropertyFromConfig(key string) string {
	if strings.HasPrefix(key, "bingo.") {
		return ""
	}
	return this.getProperty(key)
}

func (this *ApplicationContext) RegisterService(name string,service interface{}) {
	if name!="" && service!=nil {
		instance:=this.services[name]
		if instance==nil {
			this.services[name]=service
		}
	}


}

func (this *ApplicationContext)GetService(name string) interface{} {
   if name!="" {
   	   return this.services[name]
   }
   return nil
}

func (this *ApplicationContext) getProperty(key string) string {
	if this.config == nil {
		return ""
	}
	return this.config[key]
}

func(this *ApplicationContext) init(file string){
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
    this.services=make(map[string]interface{})
	this.initLogFactory()
	this.initSessionFactory()
}

func (this *ApplicationContext) initSessionFactory(){
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
}

func (this *ApplicationContext) initLogFactory(){
	this.logfactory=&utils.LogFactory{}
	this.logfactory.SetConfig(utils.LogConfig{true,this.config["bingo.log.file"]})
}

//load函数，如果加载成功返回true，否则返回FALSE
type OnLoad func(context *ApplicationContext)bool
