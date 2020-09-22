package main

import (
	"fmt"
	"github.com/aosfather/bingo_utils/contain"
	"github.com/aosfather/bingo_utils/files"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
)

//表单元信息存放目录
const (
	_FormDir        = "forms"
	_FilePathFormat = "%s/%s/%s.yaml"
)

//应用的定义
//应用由一系列forms是构成
type Application struct {
	Root  string //应用根目录
	Cache *contain.Cache
}

func (this *Application) Init() {
	this.Cache = contain.New(10*time.Minute, 0)
}

func (this *Application) GetFilePath(p string) string {
	return fmt.Sprintf("%s/%s", this.Root, p)
}

func (this *Application) GetFormMeta(name string) *FormMeta {
	if meta, exist := this.Cache.Get(name); exist {
		return meta.(*FormMeta)
	}
	//查找文件目录,从文件中加载
	filename := fmt.Sprintf(_FilePathFormat, this.Root, _FormDir, name)
	if files.IsFileExist(filename) {
		fm := &FormMeta{}
		data, err := ioutil.ReadFile(filename)
		if err == nil {
			err = yaml.Unmarshal(data, &fm)
		}
		if err != nil {
			errs("load form meta error:", err.Error())
			return nil
		}
		//放入缓存中
		this.Cache.SetDefault(name, fm)
		return fm

	}
	return nil
}

//刷新所有的表单
func (this *Application) RefreshFormAll() {
	this.Cache.Flush()
}

//重新加载form，用于从缓存中unload掉
func (this *Application) RefreshForm(formname string) {
	this.Cache.Delete(formname)
}
