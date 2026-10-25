package models

type CurrentUser struct {
	Username   string `json:"username" example:"admin"`
	Nickname   string `json:"nickname" example:"admin"`
	AvatarFile string `json:"avatar_file" example:"./avatar.png"`
}
