# bingo
the simple framework for light web application

#冰果
一个简单的web框架，用于快速的编写web应用。

##Ver 0.1 特性
1、提供参数校验Tag
2、参数自动绑定
3、根据返回对象直接转json
4、一个简单的mvc实现
5、使用标准的sql接口，没有复杂的orm映射，提供简单的结果集到struct对象的映射

##Example
hello world
hello.go

package main
import "github.com/aosfather/bingo"
type hello struct {
   bingo.Controller
}
func (this *hello) Get()(interface{},bingo.BingoError){
            return "hello world",nil
}

func main(){
    application:=bingo.Application{}
    application.addHandler("/hello",&hello{})
    application.Run("app.conf")
}


