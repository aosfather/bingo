# bingo
the simple framework for light web application

# 冰果
> 一个简单的web框架，用于快速的编写web应用。支持国内一些常用平台的应用开发，如微信公众号开发、钉钉、云之家等。
该框架是在使用beego过程中慢慢构建的，旨在直接简单，构建简单应用，该框架不是走大而全，走精干的路线，重心在以微服务应用搭建方面为主。
在一些特性支持方面往往会偏向于约定而不是配置和自定义，如果有些特性的确很重要，而约定的又不能完全解决问题，更加倾向于框架自身实现的组件接口直接公开，由应用自己来实现一个替换框架默认的实现。
>  > 用于作者是多年的javaer，在一些特性设计上也会借鉴spring mvc等java框架中个人觉的比较好的方式和包装。


# bingo 结构已经做了一定的重构，并且增加了一些默认约定来减少代码

# Road Map
* 支持通过tag来设置controller的url绑定
* gRpc服务的支持
* 钉钉应用开发支持
* 云之家应用开发支持

------------------------------------------------
## 已经实现的特性
* url的绑定支持了rest风格，允许在url上有参数，例如 /xx/:name/:id/info,其中name和id为参数
* 提供了企业微信应用的开发的支持
* 提供微信公众号开发的支持
* 提供session的实现
* 提供一些常用open api的封装
* 新增BaseDao实现，减少应用自行写操作数据库存储的代码
* 提供了Value，Inject tag，用于自动初始化对象，减少构建和初始化struct的代码
* 提供一个支持滚动的log实现
* 提供分布式lock实现
* 支持xml
* 支持redis cache
* 支持MQ-rabbit mq
* 提供了一个基于redis的倒排序搜索的实现，可以简单的实现一个简易的搜索引擎
* template变成可选设置
* 监听端口可以进行配置
*  提供参数校验Tag
*  参数自动绑定
*  根据返回对象直接转json
*  一个简单的mvc实现
*  使用标准的sql接口，没有复杂的orm映射，提供简单的结果集到struct对象的映射
*  contoller支持 Get、Post、Put、Delete方法，如果配置上数据库，则Get不提供事务控制，其它都自动提供了事务的控制

## Example
hello world
### hello.go

```c
简单例子
package main
import "github.com/aosfather/bingo"
func main(){
    application:=bingo.TApplication{}
    application.Run("")
}

默认端口 8090

```
##一个复杂点的服务
```c
package main
import (
  "fmt"
  "github.com/aosfather/bingo"
  "github.com/aosfather/bingo/mvc"
  "github.com/aosfather/bingo/utils"
)
type myStruct struct {
     name string `Value:"mytest"`        //自动赋值属性的tag，配置文件存在 mytest的属性值。如果不是公开的，则会调用Setxxx方法进行赋值
     Content *secondStruct `Inject:""`   //自动装配的tag，不指名名称，会自动装配对应类型的。如果属性不是公开的，则会调用Setxxx方法
}
func (this *myStruct)SetName(t string){
  this.name=t
}
type secondStruct struct{
    text map[string]string `Inject:""`
}

func (this *secondStruct)SetText(m map[string]string){
  fmt.Print("set text")
  this.text=m
}

func (this *secondStruct)Init(){
   fmt.Println("call the init method!")
}


func main(){
  app:=bingo.TApplication{}
  app.SetHandler(loadservice,loadcontroller) //设置服务和自动装配control的方法
  app.SetOnDestoryHandler(destory)   //设置应用被杀掉时候的响应代码
  app.Run("config.conf")
}

func destory()bool{
  fmt.Println("destory")
  return true

}

//加载服务
func loadservice(context *bingo.ApplicationContext) bool{
  context.RegisterService("test",&secondStruct{})
  context.RegisterService("test1",&myStruct{})
return true
}

//加载mvc中的control，也就是请求处理
func loadcontroller(mvc *bingo.MvcEngine,context *bingo.ApplicationContext) bool {
  mvc.AddController(&mybook{})
  fmt.Println(context.GetService("test"))
  p:=context.GetService("test1")  //获取对应的服务对象引用
  if v,ok:=p.(*myStruct);ok {
    v.Content.text["123"]="123"
    fmt.Printf("%s",v)
  }

  return true
}

//一个control，用于响应网络请求
//如果不指名对应的url，默认使用类型的名称，例如响应 /mybook
type mybook struct {
  mvc.SimpleController
  logger utils.Log
}
//初始化
func (this *mybook)Init() {
  this.logger=this.GetBeanFactory().GetLog("mybook")
}

//响应get请求
func (this *mybook) Get(c mvc.Context, p interface{}) (interface{}, mvc.BingoError) {
  this.logger.Info("haha call the mybook!")
  return "test",nil
}


```



## conf文件样例-----json格式的配置
* bingo.mvc.template 模板目录配置  --可选
* bingo.mvc.static  静态资源目录   --可选
* bingo.system.usedb 是否启用数据库，默认不启用
* bingo.system.port 指定服务监听的端口，如果不指定默认端口8990
* bingo.db 数据库的配置
* * type 类型
* * name 数据库名
* * url 格式tcp（ip：port） 
* * user 数据库用户名
* * password 数据库密码

样例如下
### app.conf
```json
{
	"bingo.mvc.template":"/",
	"bingo.mvc.static":"static",
	"bingo.system.usedb":"true",
	"bingo.db.type":"mysql",
	"bingo.db.name":"configs",
	"bingo.db.url":"tcp(127.0.0.1:3306)",
	"bingo.db.user":"root",
	"bingo.db.password":"root"
	
}
```

