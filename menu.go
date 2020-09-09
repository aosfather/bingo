package main

import (
	"github.com/aosfather/bingo_mvc/sqltemplate"
)

/**
  菜单
*/
type Menu struct {
	Code     string
	Label    string
	Children []*Menu
}

func (this *Menu) AddChild(m *Menu) {
	if m != nil {
		this.Children = append(this.Children, m)
	}
}
func (this *Menu) GetChild(code string) *Menu {
	if code != "" {
		for _, c := range this.Children {
			if c.Code == code {
				return c
			}
		}

	}
	return nil
}

//菜单项，
type MenuItem struct {
	Code       string `Option:"pk"`
	Name       string `Table:"bingo_menus"`
	Level      int    //层级
	ParentCode string //上级菜单
	Sortindex  int
	Descript   string
	Url        string
	Enabled    int
}

type MenuTree struct {
	Root    []*Menu
	menuMap map[string]*Menu
	Dao     *sqltemplate.MapperDao `Inject:"menu"`
}
