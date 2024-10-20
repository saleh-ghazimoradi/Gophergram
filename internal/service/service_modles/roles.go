package service_modles

type Roles struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Level       int    `json:"level"`
}
