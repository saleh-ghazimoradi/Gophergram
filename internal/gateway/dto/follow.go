package dto

type FollowUser struct {
	UserId int64 `json:"user_id"`
}

type UnfollowUser struct {
	UserId int64 `json:"user_id"`
}
