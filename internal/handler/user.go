package handler

import (
	"log"
	"net/http"

	"github.com/1inkun/Gubill/internal/models"
	"github.com/1inkun/Gubill/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Login(c *gin.Context) {
	response := models.Response
	response["msg"] = "登录成功"
	bodyData := struct {
		Usernmae string `json:"username" validate:"required,alphanum,min=6,max=20"`
		Password string `json:"password" validate:"required,min=8,max=30,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ,containsany=abcdefghijklmnopqrstuvwxyz,containsany=0123456789"`
	}{}
	err := c.ShouldBindJSON(&bodyData)
	if err != nil {
		models.DataBindErr(c, 400, "数据绑定出错")
		return
	}
	log.Println(bodyData)

	c.JSON(http.StatusOK, response)
}
