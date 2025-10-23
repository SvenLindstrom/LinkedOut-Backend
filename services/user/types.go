package user

type UserInfo struct {
	Id         string     `json:"id"`
	Name       string     `json:"name"`
	Bio        string     `json:"bio"`
	Profession string     `json:"Profession"`
	Interests  []Interest `json:"interests"`
}

type Interest struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
