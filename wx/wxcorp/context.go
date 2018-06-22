package wxcorp

/**
企业应用的上下文
*/

type wxAppcontext struct {
	auth        *wxAuthCorpcontext
	AuthAgentId int //应用在授权企业里的id
}

type wxAuthCorpcontext struct {
	AuthCorpId  string //授权企业id
	AccessToken string //授权企业的访问token
}
