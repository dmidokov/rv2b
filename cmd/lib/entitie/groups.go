package entitie

type Group struct {
	GroupId               int    `json:"group_id,omitempty"`
	CreatorOrganizationId int    `json:"creator_organization_id,omitempty"`
	GroupName             string `json:"group_name,omitempty"`
	GroupRights1          int64  `json:"group_rights_1,omitempty"`
	CreatorId             int    `json:"creator_id"`
	CreatorName           string `json:"creator_name"`
}

type GroupNameAndIds struct {
	GroupId   int    `json:"group_id,omitempty"`
	UserId    int    `json:"user_id,omitempty"`
	GroupName string `json:"group_name,omitempty"`
}
