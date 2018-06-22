package wxcorp

import (
	"encoding/json"
	"fmt"
)

/**
通讯录组织机构管理

*/
const (
	CONTACT_DEPARTMENT_API = "https://qyapi.weixin.qq.com/cgi-bin/department/%s?access_token=%s&id=%s"

	CONTACT_DEPARTMENT_USERLIST_API = "https://qyapi.weixin.qq.com/cgi-bin/user/list?access_token=%s&department_id=%d&fetch_child=%d"
)

//通讯录部门
type contactDepart struct {
	Id     int64  `json:"id"`       //部门id
	Name   string `json:"name"`     //部门名称
	Parent int64  `json:"parentid"` //父亲部门id。根部门为1
	Order  int64  `json:"order"`    //在父部门中的次序值。order值大的排序靠前。值范围是[0, 2^32)
}

func (this *contactDepart) ToWxDepart() WxDepart {
	return WxDepart{this.Id, this.Name, this.Parent, this.Order}

}

type contactDepartList struct {
	baseMessage
	List []contactDepart `json:"department"`
}

type contactUserList struct {
	baseMessage
	List []contactUser `json:"userlist"`
}

//用户信息
type contactUser struct {
	WxUserDetail
	Order       []int           `json:"order"`
	IsLeader    int             `json:"isleader"`
	Telephone   string          `json:"telephone"`
	EnglishName string          `json:"english_name"`
	Status      int             `json:"status"`
	Extattr     contactExtAttrs `json:"extattr"`
}

func (this *contactUser) ToWxUser() WxUser {

	user := WxUser{this.Id, this.Id, this.Name, this.Avatar, this.Departments, this.Mobile, this.Position, this.Gender, this.Email, this.Status, this.EnglishName, this.IsLeader, this.Telephone, nil}

	return user

}

type contactExtAttrs struct {
	Data []contactExtAttr `json:"attrs"`
}

type contactExtAttr struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type WxContactApi struct {
	context *wxAuthCorpcontext
}

//获取权限下的所有部门列表
func (this *WxContactApi) GetDepartmentList() contactDepartList {
	if this.context == nil {
		return contactDepartList{}
	}
	urlstr := fmt.Sprintf(CONTACT_DEPARTMENT_API, "list", this.context.AccessToken, "")
	fmt.Println(urlstr)
	content, err := HTTPGet(urlstr)
	if err == nil {
		list := contactDepartList{}
		fmt.Println(string(content))
		json.Unmarshal(content, &list)
		return list
	}

	return contactDepartList{}

}

//取某部门下的员工详情列表，不递归取部门下的子部门成员
func (this *WxContactApi) GetDepartmentUserList(depid int64) contactUserList {
	if this.context == nil {
		return contactUserList{}
	}
	urlstr := fmt.Sprintf(CONTACT_DEPARTMENT_USERLIST_API, this.context.AccessToken, depid, 0)
	fmt.Println(urlstr)
	content, err := HTTPGet(urlstr)
	if err == nil {
		list := contactUserList{}
		fmt.Println(string(content))
		json.Unmarshal(content, &list)
		return list
	}

	return contactUserList{}

}
