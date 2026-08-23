package handler

import (
	"log"
	"net/http"

	"github.com/1inkun/Gubill/internal/models"
	"github.com/1inkun/Gubill/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	userService *service.UserService
	validate    *validator.Validate
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	validate := validator.New()
	return &UserHandler{
		userService: userService,
		validate:    validate,
	}
}

var (
	ErrBindDataError     = models.NewBusinessError(400, "请求参数有误")
	ErrPasswdANDUsername = models.NewBusinessError(400, "用户名或密码不符合要求")
)

// type LoginData struct {
// 	Username string `json:"username" validate:"required,alphanum,min=4,max=20" binding:"required"`
// 	Password string `json:"password" validate:"required,min=6,max=30,containsany=abcdefghijklmnopqrstuvwxyz,containsany=0123456789" binding:"required"`
// }

// type RegisterData struct {
// 	UserName string `json:"username" validate:"required,alphanum,min=4,max=20" binding:"required"`
// 	NickName string `json:"nickname" validate:"min=0,max=16"`
// 	Password string `json:"password" validate:"required,min=6,max=30,containsany=abcdefghijklmnopqrstuvwxyz,containsany=0123456789" binding:"required"`
// 	Email    string `json:"email" validate:"required,email" binding:"required,email"`
// }

func (h *UserHandler) Login(c *gin.Context) {
	response := models.NewResponse(200, "success", "登录成功", nil)
	bodyData := models.LoginData{}
	ctx := c.Request.Context()
	err := c.ShouldBindJSON(&bodyData)
	if err != nil {
		log.Println(err.Error())
		c.Error(ErrBindDataError)
		return
	}
	// if err = h.validate.StructCtx(ctx, bodyData); err != nil {
	// 	c.Error(ErrPasswdANDUsername)
	// 	return
	// }
	// 数据处理完毕,具体逻辑交由服务层函数处理
	res, err := h.userService.Login(ctx, bodyData)
	if err != nil {
		c.Error(err)
		return
	}
	response.Data = gin.H{
		"token": res,
	}
	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) Register(c *gin.Context) {
	response := models.NewResponse(200, "success", "登录成功", nil)
	bodyData := models.RegisterData{}
	ctx := c.Request.Context()
	if err := c.ShouldBindJSON(&bodyData); err != nil {
		c.Error(ErrBindDataError)
		return
	}
	// if err := h.validate.StructCtx(ctx, bodyData); err != nil {
	// 	c.Error(ErrPasswdANDUsername)
	// 	return
	// }
	// 数据处理完毕,调用服务层函数继续
	res, err := h.userService.Register(ctx, bodyData)
	if err != nil {
		c.Error(err)
		return
	}
	response.Data = gin.H{
		"userId": res,
	}
	c.JSON(http.StatusOK, response)
}
