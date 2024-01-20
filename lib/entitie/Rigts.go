package entitie

type RightNameValue struct {
	Name  string
	Value int
}

type Right struct {
	Id    int    `json:"-"`
	Name  string `json:"name,omitempty"`
	Value int    `json:"value,omitempty"`
}
