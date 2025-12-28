package dto

type UserRegisterRequest struct {
	Login    string `json:"login" binding:"required,min=3,max=25"`
	Password string `json:"password" binding:"required,min=6"`
}

type UserLoginRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserUpdateRequest struct {
	Login    string `json:"login" binding:"min=3,max=25"`
	Password string `json:"password" binding:"min=6,max=100"`
}

type UserResponse struct {
	ID          uint   `json:"id"`
	Login       string `json:"login"`
	IsModerator bool   `json:"is_moderator"`
}

type LoginResponse struct {
	User UserResponse `json:"user"`
}

type MessageResponse struct {
	Message string `json:"message"`
}
