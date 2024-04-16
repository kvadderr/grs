package server

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/golang-grpc-proxy/docs"
	"github.com/golang-grpc-proxy/src/entities"
	"github.com/golang-grpc-proxy/src/services"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type HttpServer struct {
	api services.Api
}

func RegisterHttpServer(engine *gin.Engine, api services.Api) {

	server := HttpServer{api: api}

	engine.POST("/register", server.Register)
	engine.POST("/getUser", server.GetUser)
	engine.POST("/getToken", server.GetToken)

	engine.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

}
//	@Summary		регистрирует нового пользователя
//	@Description	регистрирует нового пользователя на основе уникальности email и phone
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			data	body		entities.User	true "body data"
//	@Success		200		{object}	RegisterResponseDto
//	@Failure		400		{object}	RegisterResponseDto
//	@Failure		500		{object}	RegisterResponseDto
//	@Router			/register [post]
func (s *HttpServer) Register(ctx *gin.Context) {
	var body entities.User
	err := ctx.ShouldBindJSON(&body)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"result": false,
			"errorMessage": services.ErrInvalidBody.Error(),
			"id": 0,
		})
		return
	}

	id, err := s.api.Register(body)

	if err == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"result": true,
			"errorMessage": "",
			"id": id,
		})
		return
	}
	
	var code = http.StatusBadRequest

	if errors.Is(err, services.ErrInternal) {
		code = http.StatusInternalServerError
	}

	ctx.JSON(code, gin.H{
		"result": false,
		"errorMessage": err.Error(),
		"id": 0,
	})
}

//	@Summary		проверяет наличие пользователя
//	@Description	проверяет наличие пользователя с таким email и password
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			data	body		GetUserDto	true "body data"
//	@Success		200		{object}	ResponseDto
//	@Failure		400		{object}	ResponseDto
//	@Failure		500		{object}	ResponseDto
//	@Router			/getUser [post]
func (s *HttpServer) GetUser(ctx *gin.Context) {
	var body GetUserDto
	err := ctx.ShouldBindJSON(&body)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"result": false,
			"errorMessage": services.ErrInvalidBody.Error(),
		})
		return
	}

	err = s.api.AuthentificateUser(body.Email, body.Password)

	if err == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"result": true,
			"errorMessage": "",
		})
		return
	}

	var code = http.StatusBadRequest

	if errors.Is(err, services.ErrInternal) {
		code = http.StatusInternalServerError
	}

	ctx.JSON(code, gin.H{
		"result": false,
		"errorMessage": err.Error(),
	})
}

//	@Summary		возвращает jwt токен
//	@Description	проверяет наличие пользователя с таким email и password, в случае наличия генерирует jwt токен с его email на один час
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			data	body		GetUserDto	true "body data"
//	@Success		200		{object}	GetTokenResponseDto
//	@Failure		400		{object}	GetTokenResponseDto
//	@Failure		500		{object}	GetTokenResponseDto
//	@Router			/getToken [post]
func (s *HttpServer) GetToken(ctx *gin.Context) {
	var body GetUserDto

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"result": false,
			"errorMessage": services.ErrInvalidBody.Error(),
			"token": "",
		})
		return 
	}

	token, err := s.api.GetToken(body.Email, body.Password)
	if err == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"result": true,
			"errorMessage": "",
			"token": token,
		})
		return
	}

	var code = http.StatusBadRequest

	if errors.Is(err, services.ErrInternal) {
		code = http.StatusInternalServerError
	}

	ctx.JSON(code, gin.H{
		"result": false,
		"errorMessage": err.Error(),
		"token": "",
	})
}