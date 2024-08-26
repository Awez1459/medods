package models


type TokenReq struct {
	UserId  string `json:"user_id"`
}

type TokenResp struct {
	Access  string `json:"access_token"`
	Refresh string `json:"refresh_token"`
}
