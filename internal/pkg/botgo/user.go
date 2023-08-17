package botgo

type User struct {
	Id               string `json:"id"`
	Username         string `json:"username"`
	Avatar           string `json:"avatar"`
	Bot              bool   `json:"bot"`
	UnionOpenid      string `json:"union_openid"`
	UnionUserAccount string `json:"union_user_account"`
}
