package service_models

import "time"

type Post struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserID    int64     `json:"user_id"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreatePostPayload struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}
