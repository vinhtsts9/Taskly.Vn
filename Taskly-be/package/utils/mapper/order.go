package mapper

import (
	"Taskly.com/m/internal/database"
	model "Taskly.com/m/internal/models"
	"Taskly.com/m/package/utils"
)

func ConvertDBOrderToModel(o database.Order) model.OrderResult {
	return model.OrderResult{
		ID:           o.ID,
		GigID:        o.GigID,
		BuyerID:      o.BuyerID,
		SellerID:     o.SellerID,
		Status:       o.Status,
		TotalPrice:   o.TotalPrice,
		OrderDate:    o.OrderDate,
		DeliveryDate: utils.PtrTimeIfValid(o.DeliveryDate),
	}
}
func ConvertDBOrderListToModelList(dbOrders []database.Order) []model.OrderResult {
	var result []model.OrderResult
	for _, o := range dbOrders {
		result = append(result, ConvertDBOrderToModel(o))
	}
	return result
}
