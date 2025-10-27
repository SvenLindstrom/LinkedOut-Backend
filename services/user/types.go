package user

type UserInfo struct {
	Id         string     `json:"id"`
	Name       string     `json:"name"`
	Profession string     `json:"profession"`
	Bio        string     `json:"bio"`
	Interests  []Interest `json:"interests"`
}

type Interest struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
