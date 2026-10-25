package handlers

import (
	"strconv"

	"github.com/1inkun/billApp/internal/models"
	"github.com/1inkun/billApp/internal/services"
	"github.com/1inkun/billApp/internal/utils"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService *services.OrderService
}

func NewOrderHandler(orderService *services.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// @Summary					新订单(签到)
// @Tags						orders
// @Accept						json
// @Produce					json
// @Param						requestData	body	models.NewOrderRequest	true	"用户签到时需要提供签到的用户ID"
// @securityDefinitions.apikey	ApiKeyAuth
// @Success					200	{object}	utils.Response{data=models.OrderResponse}
// @Failure					400	{object}	utils.ErrResponse{detail=string}
// @Failure					401	{object}	utils.ErrResponse{detail=string}	"登录状态有误"
// @Failure					500	{object}	utils.ErrResponse{detail=nil}
// @Router						/order [post]
func (h *OrderHandler) NewOrder(c *gin.Context) {
	ctx := c.Request.Context()
	requestData := models.NewOrderRequest{}
	if err := c.ShouldBindBodyWithJSON(&requestData); err != nil {
		c.Error(models.ErrBindingData)
		return
	}
	requestUser, ok := c.Get("userId")
	if !ok {
		c.Error(models.ErrLoginStatus)
		return
	}
	// 验证请求发起者是否与所要操作对象为同一人
	if requestData.UserId != requestUser.(int64) {
		c.Error(models.ErrIllegalOperation)
		return
	}
	data, err := h.orderService.NewOrder(ctx, requestData.UserId)
	if err != nil {
		c.Error(err)
		return
	}
	utils.NewResponse(c, 200, "签到成功", data)
}

// @Summary					查找订单
// @Tags						orders
// @Accept						json
// @Produce					json
// @Param						SearchOrderParams	query	models.SearchOrderParams	true	"查询订单时需要提供相应的信息,根据提供的信息不同,有多种不同的查找模式"
// @securityDefinitions.apikey	ApiKeyAuth
// @Success					200	{object}	utils.Response{data=models.OrderResponse}								"按Id查询时的返回"
// @Success					207	{object}	utils.Response{data=utils.PaginationData{data=[]models.OrderResponse}}	"按userId查询时的返回"
// @Failure					400	{object}	utils.ErrResponse{detail=string}
// @Failure					401	{object}	utils.ErrResponse{detail=string}	"登录状态有误"
// @Failure					500	{object}	utils.ErrResponse{detail=nil}
// @Router						/order [get]
func (h *OrderHandler) GetOrderData(c *gin.Context) {
	ctx := c.Request.Context()
	requestUser, ok := c.Get("userId")
	if !ok {
		c.Error(models.ErrLoginStatus)
		return
	}
	params := models.SearchOrderParams{}
	if err := c.ShouldBindQuery(&params); err != nil {
		c.Error(models.ErrBindingData)
		return
	}
	var data any
	var err error
	if params.Id != 0 {
		// 如果id存在,则为使用id查找
		data, err = h.orderService.GetOrderDataById(ctx, params.Id, requestUser.(int64))
		if err != nil {
			c.Error(err)
			return
		}
	} else {
		// 如果id不存在,则为按照userId查找
		// 比对查询的用户是否为当前用户
		if requestUser != params.UserId {
			c.Error(models.ErrIllegalOperation)
			return
		}
		data, err = h.orderService.GetOrderDataByUserId(ctx, params)
		if err != nil {
			c.Error(err)
			return
		}
	}
	utils.NewResponse(c, 200, "获取成功", data)
}

// @Summary					结算订单
// @Tags						orders
// @Accept						json
// @Produce					json
// @Param						order_id	path	integer	true	"结算时需要提供订单的orderId"
// @securityDefinitions.apikey	ApiKeyAuth
// @Success					200	{object}	utils.Response{data=models.FinishOrderResponse}
// @Failure					400	{object}	utils.ErrResponse{detail=string}
// @Failure					401	{object}	utils.ErrResponse{detail=string}	"登录状态有误"
// @Failure					500	{object}	utils.ErrResponse{detail=nil}
// @Router						/order/:order_id [put]
func (h *OrderHandler) FinishOrder(c *gin.Context) {
	ctx := c.Request.Context()
	orderId, err := strconv.Atoi(c.Param("order_id"))
	// log.Println(orderId)
	if err != nil || orderId <= 0 {
		c.Error(models.ErrBindingData)
		return
	}
	requestUser, ok := c.Get("userId")
	if !ok {
		c.Error(models.ErrLoginStatus)
		return
	}
	data, err := h.orderService.FinishOrder(ctx, int64(orderId), requestUser.(int64))
	if err != nil {
		c.Error(err)
		return
	}
	utils.NewResponse(c, 200, "结算成功", data)
}
