package dto

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginUserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type LoginResponse struct {
	Token string            `json:"token"`
	User  LoginUserResponse `json:"user"`
}

type ProfileResponse struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}