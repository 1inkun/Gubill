package handlers

import (
	"github.com/1inkun/billApp/internal/models"
	"github.com/1inkun/billApp/internal/services"
	"github.com/1inkun/billApp/internal/utils"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// @Summary	注册新账户
// @Tags		users
// @Accept		json
// @Produce	json
// @Param		registerInfo	body		models.RegisterRequest	true	"用户的注册信息"
// @Success	200				{object}	utils.Response{data=nil}
// @Failure	400				{object}	utils.ErrResponse{detail=string}
// @Failure	500				{object}	utils.ErrResponse{detail=nil}
// @Router		/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	ctx := c.Request.Context()
	// 验证，绑定数据
	registerInfo := models.RegisterRequest{}
	if err := c.ShouldBindBodyWithJSON(&registerInfo); err != nil {
		c.Error(models.ErrBindingData)
		return
	}
	// 服务层处理
	err := h.userService.Register(ctx, registerInfo)
	if err != nil {
		c.Error(err)
		return
	}
	utils.NewResponse(c, 200, "注册成功", nil)
}

// @Summary	登录
// @Tags		users
// @Accept		json
// @Produce	json
// @Param		logInfo	body		models.LoginRequest	true	"用户的登录信息"
// @Success	200		{object}	utils.Response{data=models.LoginSuccess}
// @Failure	400		{object}	utils.ErrResponse{detail=string}
// @Failure	500		{object}	utils.ErrResponse{detail=nil}
// @Router		/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()
	logInfo := models.LoginRequest{}
	if err := c.ShouldBindBodyWithJSON(&logInfo); err != nil {
		c.Error(models.ErrBindingData)
		return
	}
	tokenStr, err := h.userService.Login(ctx, logInfo)
	if err != nil {
		c.Error(err)
		return
	}
	data := models.LoginSuccess{
		Token: tokenStr,
	}
	utils.NewResponse(c, 200, "登录成功", data)
}

func (h *UserHandler) Reflash(c *gin.Context) {

}
