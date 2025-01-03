package service_models

import "time"

type Follower struct {
	UserId     int64     `json:"user_id"`
	FollowerId int64     `json:"follower_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type FollowUser struct {
	UserID int64 `json:"user_id"`
}
