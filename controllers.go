package main

import (
	"fmt"
	"github.com/aosfather/bingo_mvc"
)

type FormRequest struct {
	FormName string `Field:"_name"`
	FormType string `Field:"_type"`
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
	Metas   FormMetaManager         `mapper:"name(form);url(/form);method(GET);style(HTML)"`
}

func (this *System) Init() {
	this.engines = make(map[string]RenderEngine)
	this.engines["FORM"] = &FormEngine{}
	this.engines["QUERY"] = &QueryFormEngine{}
}
func (this *System) GetHandles() bingo_mvc.HandleMap {
	result := bingo_mvc.NewHandleMap()
	result.Add("form", this.Form, &FormRequest{})
	//result.Add("query", this.query, bingo_mvc.TypeOfMap())
	return result
}

//界面显示
func (this *System) Form(a interface{}) interface{} {
	request := a.(*FormRequest)
	debug(request)
	if engine, ok := this.engines[request.FormType]; ok {
		if engine != nil {
			//获取meta信息
			meta := this.Metas.GetFormMeta(request.FormName)
			if meta == nil {
				return fmt.Sprintf("request Form type '%s',and Form '%s' not exits! please check", request.FormType, request.FormName)
			}
			//生成模板
			buffers := engine.Render(meta)
			p := make(map[string]string)
			p["FORM_NAME"] = meta.Code
			p["FORM_TITLE"] = meta.Title
			p["FORM_ACTION"] = meta.Action
			for index, key := range engine.GetKeys() {
				p[key] = buffers[index]
			}
			return bingo_mvc.ModelView{engine.GetTemplate(), p}
		}

	}

	return fmt.Sprintf("request Form type '%s',and Form '%s' not exits! please check", request.FormType, request.FormName)

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
