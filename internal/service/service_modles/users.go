package service_modles

type Users struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"`
	CreateAt int64  `json:"create_at"`
}
