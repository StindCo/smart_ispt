package dto

type UserDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Fullname string `json:"fullname"`
	RoleTag  string `json:"role_tag"`
}
