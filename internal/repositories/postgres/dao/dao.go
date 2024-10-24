package dao

type UserTable struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"password"`
	FullName     string `json:"full_name"`
}
