package entitie

type User struct {
	Id             int
	OrganizationId int
	UserName       string
	Password       string
	ActionCode     int
	Rights         int
	CreateTime     int64
	UpdateTime     int64
	Icon           string
}

type UserShort struct {
	Id         int    `json:"id,omitempty"`
	UserName   string `json:"userName,omitempty"`
	CreateTime int    `json:"createTime"`
	UpdateTime int    `json:"updateTime"`
}

type UserIcon struct {
	ImageName string `json:"image-name"`
}
