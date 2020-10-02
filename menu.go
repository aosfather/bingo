package main

import (
	"github.com/aosfather/bingo_mvc/sqltemplate"
	"github.com/aosfather/bingo_utils/contain"
	"time"
)

//菜单项，
type MenuItem struct {
	Code        string      `Option:"pk"`
	Label       string      `Table:"bingo_menus"`
	Level       int         //层级
	ParentCode  string      `Field:"pcode"` //上级菜单
	Sortindex   int         //排序
	Description string      //描述
	Url         string      //对应的访问url
	Enabled     int         //是否启用
	Children    []*MenuItem `Option:"not"`
}

func (this *MenuItem) AddChild(m *MenuItem) {
	if m != nil {
		if m.Level != this.Level+1 {
			m.Level = this.Level + 1
		}
		m.ParentCode = this.Code
		this.Children = append(this.Children, m)
	}
}
func (this *MenuItem) GetChild(code string) *MenuItem {
	if code != "" {
		for _, c := range this.Children {
			if c.Code == code {
				return c
			}
		}

	}
	return nil
}

type MenuTree struct {
	Root    []*MenuItem
	menuMap map[string]*MenuItem
	Dao     *sqltemplate.MapperDao `Inject:"menu"`
	cache   *contain.Cache
}

func (this *MenuTree) Init() {
	this.menuMap = make(map[string]*MenuItem)
	this.cache = contain.New(10*time.Minute, 0)
	this.Load()
}

func (this *MenuTree) Load() {
	i := &MenuItem{Enabled: 1}
	menus := this.Dao.QueryAll(i, "getallmenus")
	for _, item := range menus {
		menu := item.(*MenuItem)
		debug(menu)
		if menu.Level == 0 {
			this.Root = append(this.Root, menu)
		}

		this.menuMap[menu.Code] = menu
		if menu.ParentCode != "" && menu.ParentCode != "-" {
			pmenu := this.menuMap[menu.ParentCode]
			if pmenu != nil {
				pmenu.AddChild(menu)
			}

		}

	}
}

func (this *MenuTree) GetUserMenu(user string) []*MenuItem {
	return this.Root
}
