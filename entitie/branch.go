package entitie

type Branch struct {
	Id         int    `json:"id,omitempty"`
	OrgId      int    `json:"orgId,omitempty"`
	Name       string `json:"name,omitempty"`
	Address    string `json:"address,omitempty"`
	Phone      string `json:"phone,omitempty"`
	WorkTime   string `json:"workTime,omitempty"`
	CreateTime int    `json:"createTime"`
	UpdateTime int    `json:"updateTime"`
}
