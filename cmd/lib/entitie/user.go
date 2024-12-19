package entitie

type User struct {
	Id             int    `json:"id,omitempty"`
	OrganizationId int    `json:"organizationId,omitempty"`
	UserName       string `json:"userName,omitempty"`
	Password       string `json:"-"`
	ActionCode     int    `json:"actionCode,omitempty"`
	Rights         int64  `json:"rights,omitempty"`
	CreateTime     int64  `json:"createTime,omitempty"`
	UpdateTime     int64  `json:"updateTime,omitempty"`
	Icon           string `json:"icon,omitempty"`
	Type           int    `json:"type,omitempty"`
	StartPage      string `json:"startPage,omitempty"`
}

type UserShort struct {
	Id         int    `json:"id,omitempty"`
	UserName   string `json:"userName,omitempty"`
	CreateTime int    `json:"createTime"`
	UpdateTime int    `json:"updateTime"`
	Type       int    `json:"type"`
	StartPage  string `json:"startPage"`
}

type UserSwitcher struct {
	Id       int    `json:"id,omitempty"`
	UserName string `json:"userName,omitempty"`
	Icon     string `json:"icon,omitempty"`
}

type UserIcon struct {
	ImageName string `json:"image-name"`
}

type UserInfoFull struct {
	Id                         int                  `json:"id,omitempty"`
	UserName                   string               `json:"userName,omitempty"`
	CreateTime                 int64                `json:"createTime"`
	UpdateTime                 int64                `json:"updateTime"`
	Type                       int                  `json:"type"`
	StartPage                  string               `json:"startPage"`
	OrganizationId             int                  `json:"organizationId"`
	OrganizationName           string               `json:"organizationName"`
	MapRightNameToRightId      map[string]int       `json:"mapRightNameToRightId"`
	UserRights                 []int64              `json:"userRights"`
	Icon                       string               `json:"icon"`
	UserRightsWithDescriptions []Right              `json:"userRightsWithDesription"`
	Password                   string               `json:"-"`
	ActionCode                 int                  `json:"-"`
	Navigation                 []NavigationInfoPage `json:"navigation"`
	Childs                     []UserIdAndLogin     `json:"childs"`
	HotSwitch                  []UserIdAndLogin     `json:"hotSwitch"`
	Groups                     []Group              `json:"groups"`
	AssignedGroups             []GroupNameAndIds    `json:"assignedGroups"`
}

type UserIdAndLogin struct {
	Id    int    `json:"id"`
	Login string `json:"login"`
}

func ConvertUserToUserLogin(user UserShort) UserIdAndLogin {
	var u UserIdAndLogin

	u.Id = user.Id
	u.Login = user.UserName

	return u
}
