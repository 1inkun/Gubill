package handlers

import (
	"github.com/1inkun/billApp/internal/services"
	"github.com/1inkun/billApp/internal/utils"
	"github.com/gin-gonic/gin"
)

type StatusHandler struct {
	statusService *services.StatusService
}

func NewStatusService(statusService *services.StatusService) *StatusHandler {
	return &StatusHandler{
		statusService: statusService,
	}
}

// @Summary	查看当前总体签到情况
// @Tags		orders
// @Produce	json
// @Success	200	{object}	utils.Response{data=[]models.CurrentUser}
// @Failure	500	{object}	utils.ErrResponse
// @Router		/status [get]
func (h *StatusHandler) GetSginStatus(c *gin.Context) {
	ctx := c.Request.Context()
	data, err := h.statusService.GetSginStatus(ctx)
	if err != nil {
		c.Error(err)
		return
	}
	utils.NewResponse(c, 200, "获取成功", data)
}
