package models

type LoginRequest struct {
	Username string `json:"username" binding:"required,min=4,max=20" example:"username"`
	Password string `json:"password" binding:"required,min=6,max=32" example:"password"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=4,max=20" example:"username"`
	Nickname string `json:"nickname" example:"宇宙无敌暴龙战神"`
	Password string `json:"password" binding:"required,min=6,max=32" example:"zxc123"`
	Email    string `json:"email" binding:"required,email" example:"2233@123.com"`
}

// 向数据库添加的用户信息
type RegisterData struct {
	Username   string
	Nickname   string
	PasswdHash string
	Email      string
}

// JWT载荷
type JwtClaims struct {
	UserId   int64  `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

// 登录成功

type LoginSuccess struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IlRlc3QgVXNlciIsImlhdCI6MTc4MDAwMDAwMH0.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"`
}
