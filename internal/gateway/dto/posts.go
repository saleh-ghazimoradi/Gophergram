package dto

type Post struct {
	Content string   `json:"content" validate:"required,max=1000"`
	Title   string   `json:"title" validate:"required,max=100"`
	Tags    []string `json:"tags"`
}

type UpdatePost struct {
	Content *string `json:"content" validate:"omitempty,max=1000"`
	Title   *string `json:"title" validate:"omitempty,max=100"`
}
