package bingo

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type RpcObject interface {
	MethodName() string
	Call(result interface{}) error
}

//rpc客户端
type RpcClient interface {
	BuildRpcObject(method string, param interface{}) RpcObject
}

//bingo json rpc
type BingoJsonRpcClient struct {
	server string
}

func (this *BingoJsonRpcClient) BuildRpcObject(method string, param interface{}) RpcObject {
	return &BingoJsonRpcObject{method, param, this}
}

func (this *BingoJsonRpcClient) invoke(obj *BingoJsonRpcObject, result interface{}) error {
	if obj != nil {

		url := this.server + "/" + methodtoUrl(obj.method)
		buf, _ := json.Marshal(obj.param)
		body := bytes.NewBuffer(buf)
		r, _ := http.NewRequest("post", url, body)
		r.Header.Set("Content-Type", "application/json")
		res, _ := http.DefaultClient.Do(r)
		resbody, err := ioutil.ReadAll(res.Body)
		defer res.Body.Close()
		if err == nil {
			json.Unmarshal(resbody, result)
		} else {
			return err
		}

		return nil

	}
	return CreateError(500, "rpc object is nil")
}

func methodtoUrl(method string) string {
	return method
}

type BingoJsonRpcObject struct {
	method string
	param  interface{}

	client *BingoJsonRpcClient
}

func (this *BingoJsonRpcObject) MethodName() string {
	return this.method
}

func (this *BingoJsonRpcObject) Call(result interface{}) error {
	return this.client.invoke(this, result)
}

//stands json-rpc
type JsonRpcClient struct {
	ServerUrl string
	Instance  string
}

func (this *JsonRpcClient) BuildRpcObject(method string, param interface{}) RpcObject {
	return nil
}

func (this *JsonRpcClient) invoke(obj *JsonRpcObject, result interface{}) error {
	return nil
}

type JsonRpcObject struct {
	Method string         `json:"method"` //方法名
	Param  [1]interface{} `json:"params"` //参数
	Id     uint64         `json:"id"`     //请求唯一标识
	client *JsonRpcClient
}

func (this *JsonRpcObject) MethodName() string {
	return this.Method
}

func (this *JsonRpcObject) Call(result interface{}) error {
	return this.client.invoke(this, result)
}
