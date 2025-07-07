package ticket

import (
	"github.com/gin-gonic/gin"
	"go-ecommerce-backend-api.com/internal/model"
	"go-ecommerce-backend-api.com/internal/service"
	"go-ecommerce-backend-api.com/pkg/response"
	"go-ecommerce-backend-api.com/pkg/utils"
)

var TicketItem = new(cTicketItem)

type cTicketItem struct{}

// @Summary      Get Ticket Item By ID
// @Description  Lấy thông tin chi tiết của ticket item theo ID
// @Tags         TicketItem
// @Accept       json
// @Produce      json
// @Param        id   path           int  true  "Ticket Item ID"
// @Success      200  {object}  model.GetTicketItemByIdOutput
// @Failure      400  {object}  response.ErrorResponseData
// @Failure      404  {object}  response.ErrorResponseData
// @Router       /ticket/item/{id} [get]
func (c *cTicketItem) GetTicketItemById(ctx *gin.Context) {
	var params model.TicketItemRequest
	if !utils.CheckValidParams(ctx, &params) {
		return
	}
	out, err := service.TicketItem().GetTicketItemById(ctx, &params)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrCodeTicketItemNotFound, err.Error())
		return
	}
	response.SuccessResponse(ctx, response.ErrCodeSuccess, out)
}
