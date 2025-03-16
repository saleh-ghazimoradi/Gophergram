package service_models

import "time"

type Comment struct {
	Id        int64     `json:"id"`
	PostId    int64     `json:"post_id"`
	UserId    int64     `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	User      Users     `json:"user"`
}
