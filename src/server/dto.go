package server

type GetUserDto struct {
	Email    string `json:"email" binding:"gt=0,required,email"`
	Password string `json:"password" binding:"gt=0,required"`
}

type ResponseDto struct {
	Result       bool   `json:"result" example:"true"`
	ErrorMessage string `json:"errorMessage" example:"bad request"`
}

type RegisterResponseDto struct {
	ResponseDto
	ID int `json:"id" example:"1"`
}

type GetTokenResponseDto struct {
	ResponseDto
	Token string `json:"token"`
}