package main

import (
	"encoding/json"
	"fmt"
	"github.com/aosfather/bingo_mvc"
)

type FormRequest struct {
	FormName string `Field:"_name"`
	FormType string `Field:"_type"`
}

type FormResult struct {
	Code  int           `json:"code"`
	Msg   string        `json:"msg"`
	Count int           `json:"count"`
	Data  []interface{} `json:"data"`
}

type FormRawResult struct {
	Code  int               `json:"code"`
	Msg   string            `json:"msg"`
	Count int               `json:"count"`
	Data  []json.RawMessage `json:"data"`
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
	Metas   FormMetaManager         `mapper:"name(form);url(/form/:_name);method(GET);style(HTML)" Inject:""`
	Action  *FormActions            `mapper:"name(action);url(/do/:_form_);method(POST);method(GET);style(JSON)" Inject:""`
}

func (this *System) Init() {
	this.engines = make(map[string]RenderEngine)
	this.engines["FORM"] = &FormEngine{}
	this.engines["QUERY"] = &QueryFormEngine{}
}
func (this *System) GetHandles() bingo_mvc.HandleMap {
	result := bingo_mvc.NewHandleMap()
	result.Add("form", this.Form, &FormRequest{})
	result.Add("action", this.FormAction, bingo_mvc.TypeOfMap())
	return result
}

//界面显示
func (this *System) Form(a interface{}) interface{} {
	request := a.(*FormRequest)
	debug(request)
	//获取meta信息
	meta := this.Metas.GetFormMeta(request.FormName)
	if meta == nil {
		return fmt.Sprintf("request Form type '%s',and Form '%s' not exits! please check", request.FormType, request.FormName)
	}
	if engine, ok := this.engines[meta.FormType]; ok {
		if engine != nil {
			//生成模板
			buffers, script, exscript := engine.Render(meta)
			p := make(map[string]string)
			p["FORM_NAME"] = meta.Code
			p["FORM_TITLE"] = meta.Title
			if meta.Action == "" {
				p["FORM_ACTION"] = "/do/" + meta.Code
			} else {
				p["FORM_ACTION"] = meta.Action
			}

			p["FORM_VERIFY"] = script
			p["FORM_SCRIPT"] = meta.JSscript
			p["COMPONENT_SCRIPT"] = exscript
			for index, key := range engine.GetKeys() {
				p[key] = buffers[index]
			}
			return bingo_mvc.ModelView{engine.GetTemplate(), &p}
		}

	}

	return fmt.Sprintf("request Form type '%s',and Form '%s' not exits! please check", request.FormType, request.FormName)

}

func (this *System) FormAction(a interface{}) interface{} {
	request := a.(map[string]interface{})
	debug(request)
	if formcode, ok := request["_form_"]; ok {
		meta := this.Metas.GetFormMeta(formcode.(string))
		if meta != nil {
			r, e := this.Action.Execute(meta, request)
			if e != nil {
				return FormResult{Code: 500, Msg: e.Error()}
			}

			switch meta.Response.Type {
			case PT_DIRECT:
				if result, ok := r.(string); ok {
					var body json.RawMessage
					//这个是可以直接返回结果的，用于中转服务是比较合适
					body = []byte(result)
					return body
				}
				return r

			default:
				if result, ok := r.(string); ok {
					return FormRawResult{Code: 0, Msg: "ok", Count: 1, Data: []json.RawMessage{[]byte(result)}}
				} else if result, ok := r.([]string); ok {
					datas := []json.RawMessage{}
					for _, item := range result {
						datas = append(datas, []byte(item))
					}
					return FormRawResult{Code: 0, Msg: "ok", Count: 1, Data: datas}
				}

				return FormResult{Code: 0, Msg: "ok", Count: 1, Data: []interface{}{r}}
			}
		}
	}

	return FormResult{Code: 500, Msg: "the form not exist!"}
}

//上传

//下载
//导入
//导出
//接口调用
