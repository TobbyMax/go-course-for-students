package user

type User struct {
	ID       int64  `json:"id"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}
