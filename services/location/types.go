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
	Id         string     `json:"id"         binding:"required"`
	Name       string     `json:"name"       binding:"required"`
	Bio        string     `json:"bio"`
	Distance   string     `json:"distance"   binding:"required"`
	Profession string     `json:"profession"`
	Tags       []Interest `json:"tags"`
}

type Interest struct {
	Id   string `json:"id"   binding:"required"`
	Name string `json:"name" binding:"required"`
}
