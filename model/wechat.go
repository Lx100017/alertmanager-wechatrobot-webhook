package model

type WeChatkMessage struct {
}

type WeChatMarkdown struct {
	ToTag    string    `json:"totag"`
	ToUser   string    `json:"touser"`
	ToParty  string    `json:"toparty"`
	AgentId  string    `json:"agentid"`
	MsgType  string    `json:"msgtype"`
	Markdown *Markdown `json:"markdown"`
}

type WeChatText struct {
	ToTag   string `json:"totag"`
	ToUser  string `json:"touser"`
	ToParty string `json:"toparty"`
	AgentId string `json:"agentid"`
	MsgType string `json:"msgtype"`
	Text    *Text  `json:"text"`
}

type Markdown struct {
	Content string `json:"content"`
}

type Text struct {
	Content string `json:"content"`
}

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	ErrCode     int64  `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
}

type WechatMessageCommonResponse struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}
