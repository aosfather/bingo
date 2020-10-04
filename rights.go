package main

import (
	"github.com/aosfather/bingo_mvc/hippo"
	"github.com/aosfather/bingo_mvc/sqltemplate"
	"github.com/aosfather/bingo_utils/contain"
	"time"
)

type Permission struct {
	Role      string `Field:"role_code"`
	Table     string `Field:"table_name"`
	Rowid     int
	Field     string `Field:"field_code"`
	IsVeto    bool
	Negation  bool
	ValueType string `Field:"value_type"`
	Value     string `Field:"value_data"`
}

//权限管理
type Permissions struct {
	Dao   *sqltemplate.MapperDao `Inject:"permission"`
	cache *contain.Cache
}

func (this *Permissions) Init() {
	this.cache = contain.New(10*time.Minute, 0)
}
func (this *Permissions) GetRole(r string) *hippo.Role {
	if role, ok := this.cache.Get(r); ok {
		return role.(*hippo.Role)
	}
	p := &Permission{Role: r}
	//从数据库中加载角色定义
	if this.Dao.Exist(p, "existRole") {
		datas := this.Dao.QueryAll(p, "getrole")
		role := this.createRoleByData(r, datas)
		this.cache.SetDefault(r, role)
		return role

	}
	return nil
}

func (this *Permissions) createRoleByData(r string, datas []interface{}) *hippo.Role {
	role := hippo.Role{Code: r}
	role.Init()
	var table string
	var rowid int = -1
	var authrow *hippo.AuthRow
	for _, v := range datas {
		p := v.(*Permission)
		table = p.Table
		if p.Rowid != rowid {
			row := make(hippo.AuthRow)
			authrow = &row
			role.AddRow(table, authrow, p.IsVeto)
		}

		authrow.Add(&hippo.AuthField{Code: p.Field}, hippo.ValueType(p.ValueType), p.Value, p.Negation)
	}

	return &role
}
