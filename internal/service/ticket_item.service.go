package service

import (
	"context"

	"go-ecommerce-backend-api.com/internal/model"
)

type (
	ITicketHome interface{}
	ITicketItem interface {
		GetTicketItemById(ctx context.Context, in *model.TicketItemRequest) (out *model.GetTicketItemByIdOutput, err error)
	}
)

var (
	localTicketItem ITicketItem
	localTicketHome ITicketHome
)

func TicketHome() ITicketHome {
	if localTicketHome == nil {
		panic("implement localTicketHome not found interface ITicketHome")
	}
	return localTicketHome
}
func InitTicketHome(i ITicketHome) {
	localTicketHome = i
}

func TicketItem() ITicketItem {
	if localTicketItem == nil {
		panic("implement localTicketItem not found interface ITicketItem")
	}
	return localTicketItem
}
func InitTicketItem(i ITicketItem) {
	localTicketItem = i
}
