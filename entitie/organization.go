package entitie

type Organization struct {
	Id         int    `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	Host       string `json:"host,omitempty"`
	CreateTime int    `json:"create-time"`
	UpdateTime int    `json:"update-time"`
	Creator    int    `json:"creator"`
}
