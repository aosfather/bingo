package wxcorp

import "time"

/**
企业应用的上下文
*/

type wxAppcontext struct {
	auth        *wxAuthCorpcontext
	AuthAgentId int //应用在授权企业里的id
}

type wxAuthCorpcontext struct {
	suit        *WxCorpSuite
	AuthCorpId  string //授权企业id
	AccessToken string //授权企业的访问token
	expire      int64
}

func CreateAuthCorpContext(suit *WxCorpSuite, corpId string) *wxAuthCorpcontext {
	context := wxAuthCorpcontext{}
	context.suit = suit
	context.AuthCorpId = corpId
	context.Refresh()
	return &context
}

func (this *wxAuthCorpcontext) setToken(token CorpAccessToken) {
	now := time.Now().Unix()
	this.expire = now + int64(token.Expire)
	this.AccessToken = token.Token
}

func (this *wxAuthCorpcontext) Refresh() {
	if this.AccessToken == "" || time.Now().Unix() > this.expire {
		pcode := this.suit.dataStage.GetCode(CATALOG_PERMANENT, this.AuthCorpId)
		token := this.suit.GetCorpAccessToken(this.AuthCorpId, pcode)
		this.setToken(token)
	}
}
