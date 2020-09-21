package main

import (
	"github.com/aosfather/bingo_utils/contain"
	"time"
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

func (this *Application) GetFormMeta(name string) *FormMeta {

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
