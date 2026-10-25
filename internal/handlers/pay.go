package handlers

import (
	"log"

	"github.com/1inkun/billApp/internal/models"
	"github.com/1inkun/billApp/internal/services"
	"github.com/1inkun/billApp/internal/utils"
	"github.com/gin-gonic/gin"
)

type PayHandler struct {
	payService *services.PayService
}

func NewPayHandler(payService *services.PayService) *PayHandler {
	return &PayHandler{
		payService: payService,
	}
}

// @Summary					查找支付订单
// @Tags						pay
// @Accept						json
// @Produce					json
// @Param						SearchPayOrderDataParams	query	models.SearchPayOrderDataParams	true	"查询订单时需要提供相应的信息,根据提供的信息不同,有多种不同的查找模式"
// @securityDefinitions.apikey	ApiKeyAuth
// @Success					200	{object}	utils.Response{data=models.PayOrderData}								"按Id查询时的返回"
// @Success					207	{object}	utils.Response{data=utils.PaginationData{data=[]models.PayOrderData}}	"按userId查询时的返回"
// @Failure					400	{object}	utils.ErrResponse{detail=string}
// @Failure					401	{object}	utils.ErrResponse{detail=string}	"登录状态有误"
// @Failure					500	{object}	utils.ErrResponse{detail=nil}
// @Router						/pay [get]
func (h *PayHandler) GetPayOrderData(c *gin.Context) {
	var data any
	ctx := c.Request.Context()
	param := models.SearchPayOrderDataParams{}
	err := c.ShouldBindQuery(&param)
	if err != nil {
		c.Error(models.ErrBindingData)
		return
	}
	if param.Id > 0 {
		// 按照payOrderId查找
		data, err = h.payService.GetPayOrderDataById(ctx, param.Id)
		if err != nil {
			c.Error(err)
			return
		}
	} else {
		// 按照userId查找
		requestId, ok := c.Get("userId")
		if !ok {
			c.Error(models.ErrLoginStatus)
			return
		}
		log.Println(requestId)
		if param.UserId != requestId.(int64) {
			c.Error(models.ErrIllegalOperation)
			return
		}
		data, err = h.payService.GetPayOrderDataByUserId(ctx, param)
		if err != nil {
			c.Error(err)
			return
		}
	}
	utils.NewResponse(c, 200, "获取成功", data)
}
