package auth

type UserInfo struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type oAuthPayload struct {
	Code string `json:"code" binding:"required"`
}

type LoginRes struct {
	Token    string   `json:"token"`
	New      bool     `json:"new"`
	UserInfo UserInfo `json:"userInfo"`
}
