# bingo
the simple admin web application engine  
一个简单的管理web应用开发引擎实现
# 冰果
一个简单的基于元模型的web应用开发引擎。  
适合熟悉后台开发或脚本开发人员快速开发一个全功能的管理应用。
## 特性
* 支持输入参数的检查
* 支持LUA脚本作为主逻辑处理
* 具有统一的桌面管理界面
* 支持权限控制
* 表单元模型描述，自动生成交互界面


## 引入的技术
* bingo_dao 数据元模型库
* bingo_mvc 双引擎的mvc框架
* goperlua，基于go实现的lua脚本vm

 
  
## Road Map
### version1.X
目标：实现简单的后台管理类web应用的开发 
具有的特性
* 支持表单类型元信息  [V]
* 支持数据字典       [V]
* 支持一个数据库链接  [V]
* 支持redis作为缓存  [V]
* 支持LUA脚本       [V]
  * HTTP访问        [O]
  * 数据库访问       [O]
  * JSON解析        [V]
  * 缓存访问         [O]
* 支持用户鉴权及权限控制 [O]
* 支持菜单的配置    [O]    
