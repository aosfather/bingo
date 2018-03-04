# bingo
the simple framework for light web application

# 冰果
> 一个简单的web框架，用于快速的编写web应用。
该框架是在使用beego过程中慢慢构建的，旨在直接简单，构建简单应用，该框架不是走大而全，走精干的路线，重心在以微服务应用搭建方面为主。
在一些特性支持方面往往会偏向于约定而不是配置和自定义，如果有些特性的确很重要，而约定的又不能完全解决问题，更加倾向于框架自身实现的组件接口直接公开，由应用自己来实现一个替换框架默认的。
>  > 用于作者是多年的javaer，在一些特性设计上也会借鉴spring mvc等java框架中个人觉的比较好的方式和包装。

# Road Map
## Ver 0.4 
* 增加脚本插件实现，支持lua脚本

------------------------------------------------
## Ver 0.3 the current development version
* 支持xml
* 支持微信
* 支持redis cache
* 支持MQ


## Ver 0.2
* 提供基于json的rpc调用支持
* 实现和基于bingo框架的服务互相调用
* template变成可选设置
* 监听端口可以进行配置


## Ver 0.1 
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
package main
import "github.com/aosfather/bingo"
type hello struct {
   bingo.Controller
}
func (this *hello) Get(bingo.Context c,p interface{})(interface{},bingo.BingoError){
            return "hello world",nil
}

func main(){
    application:=bingo.Application{}
    application.addHandler("/hello",&hello{})
    application.Run("app.conf")
}
```

## conf文件样例-----json格式的配置
* bingo.mvc.template 模板目录配置
* bingo.mvc.static  静态资源目录
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

