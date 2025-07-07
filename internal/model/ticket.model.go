package model

// get ticket item by id output
type GetTicketItemByIdOutput struct {
	TichketId      int64  `json:"id"`
	TicketName     string `json:"ticket_name"`
	StockInitial   int32  `json:"stock_initial"`
	StockAvailable int32  `json:"stock_available"`
}

// DTO
type TicketItemRequest struct {
	TicketId int64 `uri:"id" binding:"required"`
}
