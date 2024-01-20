package entitie

type Navigation struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Tooltip string `json:"tooltip"`
	Group   int    `json:"group"`
	Icon    string `json:"icon"`
	Link    string `json:"link"`
}
