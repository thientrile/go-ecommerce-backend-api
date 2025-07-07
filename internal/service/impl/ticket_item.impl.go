package impl

import (
	"context"
	"fmt"

	"go-ecommerce-backend-api.com/internal/database"
	"go-ecommerce-backend-api.com/internal/model"
)

type sTicketItem struct {
	r *database.Queries
}

func NewTicketItem(r *database.Queries) *sTicketItem {
	return &sTicketItem{
		r: r,
	}
}

func (s *sTicketItem) GetTicketItemById(ctx context.Context, in *model.TicketItemRequest) (out *model.GetTicketItemByIdOutput, err error) {
	fmt.Println("GetTicketItemById called with TicketId:", in.TicketId)
	ticketItem, err := s.r.GetTicketItemByID(ctx, in.TicketId)
	if err != nil {
		fmt.Println("Error fetching ticket item:", err)
		return nil, fmt.Errorf("ticket item not found: %w", err)
	}
	if ticketItem.ID == 0 {
		fmt.Println("Ticket item not found for ID:", in.TicketId)
		return nil, fmt.Errorf("ticket item not found with ID: %d", in.TicketId)
	}
	// mapper
	out = &model.GetTicketItemByIdOutput{
		TichketId:      ticketItem.ID,
		TicketName:     ticketItem.Name,
		StockInitial:   ticketItem.StockInitial,
		StockAvailable: ticketItem.StockAvailable,
	}
	fmt.Println("Fetched ticket item:", out)
	return out, nil
}
