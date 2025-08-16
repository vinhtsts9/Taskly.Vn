package impl

import (
	"context"
	"fmt"
	"time"

	"Taskly.com/m/internal/database"
	model "Taskly.com/m/internal/models"
)

type sDisputeService struct {
	store database.Store
}

func NewDisputeService(store database.Store) *sDisputeService {
	return &sDisputeService{store: store}
}

func (s *sDisputeService) CreateDispute(ctx context.Context, input model.CreateDisputeParams) (model.Dispute, error) {
	order, err := s.store.GetOrderByID(ctx, input.OrderID)
	if err != nil {
		return model.Dispute{}, err
	}

	if order.Status != "delivered" {
		return model.Dispute{}, fmt.Errorf("only delivered orders can be disputed")
	}

	if time.Since(order.OrderDate) > 72*time.Hour {
		return model.Dispute{}, fmt.Errorf("dispute window has passed")
	}

	d, err := s.store.CreateDispute(ctx, database.CreateDisputeParams{
		OrderID: input.OrderID,
		UserID:  input.UserID,
		Reason:  input.Reason,
	})
	if err != nil {
		return model.Dispute{}, err
	}

	return model.Dispute{
		ID:        d.ID,
		OrderID:   d.OrderID,
		UserID:    d.UserID,
		Reason:    d.Reason,
		Status:    d.Status.(string),
		CreatedAt: d.CreatedAt,
	}, nil
}

func (s *sDisputeService) ListDisputes(ctx context.Context) ([]model.Dispute, error) {
	disputes, err := s.store.ListDisputes(ctx)
	if err != nil {
		return nil, err
	}

	var result []model.Dispute
	for _, d := range disputes {
		result = append(result, model.Dispute{
			ID:        d.ID,
			OrderID:   d.OrderID,
			UserID:    d.UserID,
			Reason:    d.Reason,
			Status:    d.Status.(string),
			CreatedAt: d.CreatedAt,
		})
	}
	return result, nil
}

func (s *sDisputeService) UpdateDisputeStatus(ctx context.Context, input model.UpdateDisputeStatusParams) error {
	// 1. Cập nhật trạng thái tranh chấp
	err := s.store.UpdateDisputeStatus(ctx, database.UpdateDisputeStatusParams{
		ID:     input.ID,
		Status: input.Status,
	})
	if err != nil {
		return err
	}

	// 2. Nếu là hoàn tiền → huỷ đơn
	if input.Status == "refunded" {
		dispute, err := s.store.GetDisputeByID(ctx, input.ID)
		if err != nil {
			return err
		}
		err = s.store.UpdateOrderStatus(ctx, database.UpdateOrderStatusParams{
			ID:     dispute.OrderID,
			Status: "cancelled",
		})
		if err != nil {
			return err
		}
	}

	return nil
}
