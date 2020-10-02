package main

import "github.com/aosfather/bingo_mvc"

/**桌面

 */
type Desktop struct {
	Title string    `mapper:"name(desktop);url(/desktop);method(GET);style(HTML)" Value:"app.title"`
	Tree  *MenuTree `mapper:"name(index);url(/index);method(GET);style(HTML)" Inject:""`
}

func (this *Desktop) GetHandles() bingo_mvc.HandleMap {
	result := bingo_mvc.NewHandleMap()
	result.Add("index", this.Index, bingo_mvc.TypeOfMap())
	result.Add("desktop", this.Desktop, bingo_mvc.TypeOfMap())
	return result
}

//整个工作界面框架
func (this *Desktop) Index(a interface{}) interface{} {
	//获取用户信息

	datas := make(map[string]interface{})
	datas["Title"] = this.Title
	datas["Name"] = "测试用户"
	datas["Menus"] = this.Tree.GetUserMenu("")
	return bingo_mvc.ModelView{"index", datas}
}

//桌面
func (this *Desktop) Desktop(a interface{}) interface{} {
	datas := make(map[string]string)
	datas["Title"] = "测试用桌面"
	return bingo_mvc.ModelView{"desktop", datas}
}
