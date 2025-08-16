package service

import (
	"context"

	model "Taskly.com/m/internal/models"
)

type IDisputeService interface {
	CreateDispute(ctx context.Context, input model.CreateDisputeParams) (model.Dispute, error)
	ListDisputes(ctx context.Context) ([]model.Dispute, error)
	UpdateDisputeStatus(ctx context.Context, input model.UpdateDisputeStatusParams) error
}

var localDisputeService IDisputeService

func GetDisputeService() IDisputeService {
	if localDisputeService == nil {
		panic("DisputeService not implemented")
	}
	return localDisputeService
}

func InitDisputeService(s IDisputeService) {
	localDisputeService = s
}
