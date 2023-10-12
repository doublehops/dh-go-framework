package model

type Author struct {
	BaseModel
	Name string `json:"author"`
}
