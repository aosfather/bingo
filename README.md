# bingo
the simple framework for light web application

#冰果
一个简单的web框架，用于快速的编写web应用。
该框架是在使用beego过程中慢慢构建的，旨在直接简单，构建简单应用，
以少写啰嗦代码为目标。


##Ver 0.1 特性
*  提供参数校验Tag
*  参数自动绑定
*  根据返回对象直接转json
*  一个简单的mvc实现
*  使用标准的sql接口，没有复杂的orm映射，提供简单的结果集到struct对象的映射
*  contoller支持 Get、Post、Put、Delete方法，如果配置上数据库，则Get不提供事务控制，其它都自动提供了事务的控制

##Example
hello world
###hello.go

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

