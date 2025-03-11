package service_models

import (
	"github.com/lib/pq"
	"time"
)

type Post struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserID    int64     `json:"user_id"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Version   int       `json:"version"`
	Comments  []Comment `json:"comments"`
	User      Users     `json:"user"`
}

type PostWithMetadata struct {
	ID           int64     `json:"id" boil:"id"`
	Content      string    `json:"content" boil:"content"`
	Title        string    `json:"title" boil:"title"`
	UserID       int64     `json:"user_id" boil:"user_id"`
	Tags         []string  `json:"tags" boil:"tags"`
	CreatedAt    time.Time `json:"created_at" boil:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" boil:"updated_at"`
	Version      int       `json:"version" boil:"version"`
	Comments     []Comment `json:"comments" boil:"-"`
	User         Users     `json:"user" boil:"-"`
	CommentCount int64     `json:"comment_count" boil:"comment_count"`
}

type RawPostWithMetadata struct {
	ID           int64
	UserID       int64
	Title        string
	Content      string
	CreatedAt    time.Time
	Version      int
	Tags         pq.StringArray
	Username     string
	CommentCount int64
}
