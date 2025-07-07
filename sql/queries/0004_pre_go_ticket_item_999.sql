-- file: pre_go_ticket_item_999.sql



-- name: GetTicketItemByID :one
SELECT id, name, stock_initial ,stock_available
FROM pre_go_ticket_item_999 
WHERE id = ?;