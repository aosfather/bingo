package main

//应用的定义
//应用由一系列forms是构成
type Application struct {
	Root string //应用根目录

}

func (this *Application) GetFormMeta(name string) *FormMeta {

	return nil
}
