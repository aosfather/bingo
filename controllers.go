package main

import (
	"fmt"
	"github.com/aosfather/bingo_mvc"
)

type FormRequest struct {
	FormName string `Field:"name"`
	FormType string `Field:"type"`
}

/*
  系统接口
   1、表单显示
   2、数据提交
      新增、更新、删除
   3、数据查询
*/
type System struct {
	engines map[string]RenderEngine //引擎
}

func (this *System) Init() {
	this.engines = make(map[string]RenderEngine)
	this.engines[""] = nil
}

//界面显示
func (this *System) Form(a interface{}) interface{} {
	request := a.(*FormRequest)
	if engine, ok := this.engines[request.FormType]; ok {
		engine.Render(nil, nil)
		return bingo_mvc.ModelView{engine.GetTemplate(), nil}
	}

	return fmt.Sprintf("request Form type %s,and %s not exits! please check", request.FormType, request.FormName)

}

//新增
func (this *System) Add() {

}

//更新
func (this *System) Update() {

}

//删除
func (this *System) Delete() {

}

//查询
func (this *System) Query() {

}

//上传

//下载
//导入
//导出
