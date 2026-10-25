package handlers

import (
	"strconv"

	"github.com/1inkun/billApp/internal/models"
	"github.com/1inkun/billApp/internal/services"
	"github.com/1inkun/billApp/internal/utils"
	"github.com/gin-gonic/gin"
)

type GoodHandler struct {
	goodService *services.GoodService
}

func NewGoodHandler(goodService *services.GoodService) *GoodHandler {
	return &GoodHandler{
		goodService: goodService,
	}
}

// @Summary	获取当前在售的会员商品信息
// @Tags		good
// @Produce	json
// @Success	200	{object}	utils.Response{data=[]models.Good}
// @Failure	500	{object}	utils.ErrResponse
// @Router		/good [get]
func (h *GoodHandler) GetGoods(c *gin.Context) {
	ctx := c.Request.Context()
	data, err := h.goodService.GetGoods(ctx)
	if err != nil {
		c.Error(err)
		return
	}
	utils.NewResponse(c, 200, "获取成功", data)
}

// @Summary					购买会员商品
// @Tags						good
// @Accept						json
// @Produce					json
// @param						good_id	path	integer	true	"所要购买的商品id"
// @securityDefinitions.apikey	ApiKeyAuth
// @Success					200	{object}	utils.Response{data=[]models.BuyGoodResponse}
// @Failure					400	{object}	utils.ErrResponse
// @Failure					401	{object}	utils.ErrResponse{detail=string}	"登录状态有误"
// @Failure					500	{object}	utils.ErrResponse
// @Router						/good/:good_id [post]
func (h *GoodHandler) BuyGoods(c *gin.Context) {
	ctx := c.Request.Context()
	goodId, err := strconv.Atoi(c.Param("good_id"))
	if err != nil {
		c.Error(err)
		return
	}
	requestUser, ok := c.Get("userId")
	if !ok {
		c.Error(models.ErrLoginStatus)
		return
	}
	data, err := h.goodService.BuyGood(ctx, int64(goodId), requestUser.(int64))
	if err != nil {
		c.Error(err)
		return
	}
	utils.NewResponse(c, 200, "订单创建成功", data)
}
