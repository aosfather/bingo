# bingo
the simple framework for light web application

# 冰果
> 一个简单的web框架，用于快速的编写web应用。支持国内一些常用平台的应用开发，如微信公众号开发、钉钉、云之家等。
该框架是在使用beego过程中慢慢构建的，旨在直接简单，构建简单应用，该框架不是走大而全，走精干的路线，重心在以微服务应用搭建方面为主。
在一些特性支持方面往往会偏向于约定而不是配置和自定义，如果有些特性的确很重要，而约定的又不能完全解决问题，更加倾向于框架自身实现的组件接口直接公开，由应用自己来实现一个替换框架默认的实现。
>  由于作者是多年的javaer，在一些特性设计上也会借鉴spring mvc等java框架中个人觉的比较好的方式和包装。


# bingo V2.0
> 整体结构重构，分为了bingo、bingo_mvc、bingo_dao、bingo_utils、bingo_wx、bingo_dingding
> 其中：
* bingo 具有ioc的特性，boot的特性，提供了整合cache、mvc、dao、mq等，可以快速构建一个微服务应用
* bingo_mvc 是一个目标实现go lang版本的spring mvc框架，但目标是微服务，而不是旧的web应用，会对一些特性会消减
* bingo_dao 是一个很轻量级的orm实现，实现了简单
* bingo_utils 一些好用的工具类，或者是对于java转过来的程序员更加友好的封装。
* bingo_wx 开发企业微信和微信公众号服务的封装，不用关注微信平台的通讯与接口等，只需要实现相应的响应方法即可。
* bingo_dingding 开发钉钉第三方服务的go的封装。
* search搜索的特性移到candy项目中，做为candy的一部分。

# Road Map
## V2.1计划
* 支持YML格式的配置文件
* 支持Rest服务XML格式的序列化输出

## V2.2计划
* 提供文件上传的默认实现
* 提供透明缓存的实现
* 提供接口安全认证

## V2.3计划
* 优化参数校验，新增默认校验规则
* 自定义错误界面优化
* 优化序列化输出，提供新的注解
* 优化拦截器实现，丰富拦截策略

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
    application:=bingo.Application{}
    application.Run("")
}

默认端口 8990

```
