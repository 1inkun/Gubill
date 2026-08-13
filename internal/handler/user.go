package handler

import (
	"errors"
	"net/http"

	"github.com/1inkun/Gubill/internal/models"
	"github.com/1inkun/Gubill/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

var (
	ErrBindDataError = errors.New("数据绑定出错")
	ErrWrongReqData  = errors.New("请求所用的数据不符合要求")
)

type LoginData struct {
	Usernmae string `json:"username" validate:"required,alphanum,min=5,max=20"`
	Password string `json:"password" validate:"required,min=6,max=30,containsany=abcdefghijklmnopqrstuvwxyz,containsany=0123456789"`
}

type RegisterData struct {
	UserName string `json:"username" validate:"required,alphanum,min=4,max=20"`
	NickName string `json:"nickname" validate:"min=0,max=16"`
	Password string `json:"password" validate:"required,min=6,max=30,containsany=abcdefghijklmnopqrstuvwxyz,containsany=0123456789"`
	Email    string `json:"email" validate:"required,email"`
}

func (h *UserHandler) Login(c *gin.Context) {
	response := models.NewResponse()
	response["msg"] = "登录成功"
	bodyData := LoginData{}
	ctx := c.Request.Context()
	err := c.ShouldBindJSON(&bodyData)
	if err != nil {
		models.DataBindErr(c, 400, ErrBindDataError.Error())
		return
	}
	if bodyData.Usernmae == "" || bodyData.Password == "" {
		models.DataBindErr(c, 400, ErrWrongReqData.Error())
		return
	}
	validate := validator.New()
	if err = validate.StructCtx(ctx, bodyData); err != nil {
		models.DataBindErr(c, 400, ErrWrongReqData.Error())
		return
	}
	// 数据处理完毕,具体逻辑交由服务层函数处理
	res, err := h.userService.Login(ctx, bodyData.Usernmae, bodyData.Password)
	if err != nil {
		c.Error(err)
		return
	}
	response["data"] = gin.H{
		"token": res,
	}
	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) Register(c *gin.Context) {
	response := models.NewResponse()
	response["msg"] = "注册成功"
	bodyData := RegisterData{}
	ctx := c.Request.Context()
	if err := c.ShouldBindJSON(&bodyData); err != nil {
		models.DataBindErr(c, 400, ErrBindDataError.Error())
		return
	}
	if bodyData.UserName == "" || bodyData.Password == "" || bodyData.Email == "" {
		models.DataBindErr(c, 400, ErrWrongReqData.Error())
		return
	}
	validate := validator.New()
	if err := validate.StructCtx(ctx, bodyData); err != nil {
		models.DataBindErr(c, 400, ErrWrongReqData.Error())
		return
	}
	// 数据处理完毕,调用服务层函数继续
	res, err := h.userService.Register(ctx, bodyData.UserName, bodyData.NickName, bodyData.Password, bodyData.Email)
	if err != nil {
		c.Error(err)
		return
	}
	response["data"] = gin.H{
		"userId": res,
	}
	c.JSON(http.StatusOK, response)
}
