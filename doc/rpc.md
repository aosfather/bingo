# 远程调用
## bingo rpc-json
与rpc-json协议唯一的不同就是服务端使用的是标准的url，而没有一个专门的服务来处理rpc的方法调用。
基本的映射关系 method -> url,而对应的url的服务返回json格式数据
method =>url的转换规则如下
* 英文点号 => /
* 驼峰形式转换为下划线，XxYy=>xx_yy

只有一个参数，该参数会转换成 json格式post到服务，因此所有的服务都将是支持post方法。
### example


## rpc-json
实现了标准的rpc-json协议的客户端


## rpc-gob
基于gob的rpc客户端调用