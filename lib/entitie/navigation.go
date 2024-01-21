package entitie

type Navigation struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Tooltip string `json:"tooltip"`
	Group   int    `json:"group"`
	Icon    string `json:"icon"`
	Link    string `json:"link"`
}

type NavigationInfoPage struct {
	Id      int    `json:"id"`
	Title   string `josn:"title"`
	Enabled bool   `json:"enabled"`
}

type NavigationAvailable struct {
	UserId   int
	EntityId int
	GroupId  int
}
