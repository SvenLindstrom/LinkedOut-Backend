package location

type Location struct {
	Lat float64 `json:"lat" binding:"required"`
	Lon float64 `json:"lon" binding:"required"`
}

type UserProx struct {
	Name string
}
