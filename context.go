package bingo

import (
	"strings"
	"github.com/aosfather/bingo/utils"
	"github.com/aosfather/bingo/sql"
	"os"
	"io/ioutil"
	"encoding/json"
	"fmt"
)

type ApplicationContext struct {
	config  map[string]string
	logfactory *utils.LogFactory
	factory *sql.SessionFactory
	//services map[string]interface{}
	services utils.InjectMan
	holder utils.ValuesHolder
}

func (this *ApplicationContext) shutdown(){

	//关闭所有service

	this.logfactory.Close()

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

func (this *ApplicationContext)CreateDao() *BaseDao {
	return &BaseDao{this}
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
		instance:=this.services.GetObjectByName(name)
		if instance==nil {
			this.holder.ProcessValueTag(service)
			this.services.AddObjectByName(name,service)
		}
	}


}

func (this *ApplicationContext)GetService(name string) interface{} {
   if name!="" {
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
    //this.services=make(map[string]interface{})
	this.services.Init(nil)
	this.services.AddObject(this)
	this.holder.InitByFunction(this.GetPropertyFromConfig)
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

func (this *ApplicationContext)initServices(){

}

//load函数，如果加载成功返回true，否则返回FALSE
type OnLoad func(context *ApplicationContext)bool



const MAX_ROW_COUNT = 100000//最大获取条数
//基础数据操作对象
type BaseDao struct {
	context *ApplicationContext
}

func (this *BaseDao)Init(c *ApplicationContext){
	this.context=c
}

//插入，返回auto id和错误信息
func(this *BaseDao) Insert(obj utils.Object)(int64,error){
	if obj==nil {
		return 0,fmt.Errorf("nil object!")
	}

	session:=this.context.GetSession()
	if session!=nil {
		defer session.Close()
		session.Begin()
		id,count,err:=session.Insert(obj)

		if err==nil&&count==1{
			session.Commit()
			return id,nil
		}else {
			session.Rollback()
			return 0,err
		}

	}

	return 0,fmt.Errorf("session is nil")

}

func (this *BaseDao) Find(obj utils.Object,cols ...string) bool{
	if obj==nil {
		return false
	}
	session:=this.context.GetSession()
	if session!=nil {
		defer session.Close()
		return session.Find(obj,cols...)
	}
   return false
}

func (this *BaseDao) FindBySql(obj utils.Object,sqlTemplate string,args ...interface{}) bool{
	if obj==nil {
		return false
	}
	session:=this.context.GetSession()
	if session!=nil {
		defer session.Close()
		return session.Query(obj,sqlTemplate,args...)
	}
	return false
}

//更新，返回更新的条数和错误信息
func (this *BaseDao) Update(obj utils.Object,cols ... string)(int64,error) {
	if obj==nil {
		return 0,fmt.Errorf("nil object!")
	}

	session:=this.context.GetSession()
	if session!=nil {
		defer session.Close()
		session.Begin()
		_,count,err:=session.Update(obj,cols...)
		if err!=nil {
			session.Rollback()
			return 0,err
		}else {
			session.Commit()
			return count,nil
		}
	}

	return 0,fmt.Errorf("session is nil")
}


func (this *BaseDao) Delete(obj utils.Object,cols ... string)(int64,error) {
	if obj==nil {
		return 0,fmt.Errorf("nil object!")
	}

	session:=this.context.GetSession()
	if session!=nil {
		defer session.Close()
		session.Begin()
		_,count,err:=session.Delete(obj,cols...)
		if err!=nil {
			session.Rollback()
			return 0,err
		}else {
			session.Commit()
			return count,nil
		}
	}

	return 0,fmt.Errorf("session is nil")
}

func (this *BaseDao) QueryAll(obj utils.Object,cols ...string)([]interface{}){

	return this.Query(obj,sql.Page{MAX_ROW_COUNT,1,0},cols...)
}

func (this *BaseDao) Query(obj utils.Object,page sql.Page,cols ...string)([]interface{}){

	if obj==nil {
		return nil
	}
	session:=this.context.GetSession()
	if session!=nil {
		defer session.Close()
		theSql,args,err:=sql.CreateQuerySql(obj,cols...)
		if err==nil {
			return session.QueryByPage(obj,page,theSql,args...)
		}

	}

	return nil
}

//执行单条sql
func(this *BaseDao)Exec(sqltemplate string,objs ...interface{})(int64,error) {
	session:=this.context.GetSession()
	if session!=nil {
		defer session.Close()
		session.Begin()
		_,count,err:=session.ExeSql(sqltemplate,objs...)
		if err==nil {
			session.Commit()
			return count,err
		}

		session.Rollback()
		return 0,err
	}
	return 0,fmt.Errorf("session is nil")
}

//批量执行简单的sql语句
func (this *BaseDao)ExecSqlBatch(sqls ...string) error {
	session:=this.context.GetSession()
	if session!=nil {
		defer session.Close()
		session.Begin()
		for _,sql:=range sqls {
		_,_,err:=session.ExeSql(sql)
		if err!=nil {
			session.Rollback()
			return err
		}

		}
		session.Commit()
		return nil
	}
	return fmt.Errorf("session is nil")
}

func (this *BaseDao)GetSession() *sql.TxSession{
	return this.context.GetSession()
}


//插入对象，并根据对象的id，更新后续的sql语句，一般为update，其中这个关联id必须是第一个参数
func (this *BaseDao)InsertAndUpdate(iobj interface{},sqltemplate string,args ...interface{}) error{
	if iobj==nil{
		return fmt.Errorf("nil object!")
	}

	session:=this.context.GetSession()
	if session!=nil {
		defer session.Close()
		session.Begin()
		id,count,err:=session.Insert(iobj)
		if err==nil&&count==1 {
			p:=[]interface{}{id}
			if args!=nil&&len(args)>0 {
				p=append(p,args...)
			}

			_,_,err=session.ExeSql(sqltemplate,p)

			if err==nil {
				session.Commit()
			}

		}

		session.Rollback()
		return err
	}

	return fmt.Errorf("session is nil")

}

