package location

type Location struct {
	Lat float64 `json:"lat" binding:"required"`
	Lon float64 `json:"lon" binding:"required"`
}

type Proximity struct {
	Location Location `json:"location" binding:"required"`
	Distance int32    `json:"distance" binding:"required"`
}

type UserProx struct {
	Info     UserInfo `json:"info"     binding:"required"`
	Distance string   `json:"distance" binding:"required"`
}

type UserInfo struct {
	Id         string     `json:"id"`
	Name       string     `json:"name"`
	Profession string     `json:"profession"`
	Bio        string     `json:"bio"`
	Interests  []Interest `json:"interests"`
}

type Interest struct {
	Id   string `json:"id"   binding:"required"`
	Name string `json:"name" binding:"required"`
}
