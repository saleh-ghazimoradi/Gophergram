package service_modles

import "time"

type Post struct {
	ID        int64      `json:"id"`
	Content   string     `json:"content"`
	Title     string     `json:"title"`
	UserID    int64      `json:"user_id"`
	Tags      []string   `json:"tags"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Version   int        `json:"version"`
	Comment   []Comments `json:"comments"`
	User      Users      `json:"user"`
}

type PostWithMetaData struct {
	Post
	CommentsCount int `json:"comments_count"`
}
