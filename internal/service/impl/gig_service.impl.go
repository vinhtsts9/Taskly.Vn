package impl

import (
	"context"

	"Taskly.com/m/global"
	"go.uber.org/zap"

	"Taskly.com/m/internal/database"
	model "Taskly.com/m/internal/models"
	utils "Taskly.com/m/package/utils"
	"github.com/google/uuid"
)

type sGigService struct {
	store database.Store
}

func NewGigService(store database.Store) *sGigService {
	return &sGigService{
		store: store,
	}
}

// 1. Lấy danh sách dịch vụ (có lọc, tìm kiếm)
// func (s *sGigService) ListServices(ctx context.Context, params model.ListServicesParams) ([]model.Gig, error) {
// 	dbGigs, err := s.store.ListServices(ctx, database.ListServicesParams{
// 		Search:     utils.ToNullString(params.Search),
// 		CategoryID: utils.ToNullInt32(params.CategoryID).Int32,
// 		Status:     params.Status,
// 		Limit:      params.Limit,
// 		Offset:     params.Offset,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	var gigs []model.Gig
// 	for _, g := range dbGigs {
// 		gigs = append(gigs, model.Gig{
// 			ID:           g.ID,
// 			UserID:       g.UserID,
// 			Title:        g.Title,
// 			Description:  g.Description,
// 			CategoryID:   g.CategoryID,
// 			Price:        g.Price,
// 			DeliveryTime: g.DeliveryTime,
// 			ImageURL:     utils.PtrIfValid(g.ImageUrl),
// 			Status:       g.Status,
// 			CreatedAt:    g.CreatedAt,
// 		})
// 	}
//
// 	return gigs, nil
// }

// 2. Tạo dịch vụ mới
func (s *sGigService) CreateService(ctx context.Context, input model.CreateServiceParams) (model.Gig, error) {
	gig, err := s.store.CreateService(ctx, database.CreateServiceParams{
		UserID:       input.UserID,
		Title:        input.Title,
		Description:  input.Description,
		CategoryID:   input.CategoryID,
		Price:        input.Price,
		DeliveryTime: input.DeliveryTime,
		ImageUrl:     utils.ToNullString(input.ImageURL),
		Status:       input.Status,
	})
	if err != nil {
		return model.Gig{}, err
	}

	return model.Gig{
		ID:           gig.ID,
		UserID:       gig.UserID,
		Title:        gig.Title,
		Description:  gig.Description,
		CategoryID:   gig.CategoryID,
		Price:        gig.Price,
		DeliveryTime: gig.DeliveryTime,
		ImageURL:     utils.PtrIfValid(gig.ImageUrl),
		Status:       gig.Status,
		CreatedAt:    gig.CreatedAt,
	}, nil
}

// 3. Lấy chi tiết dịch vụ
func (s *sGigService) GetServiceByID(ctx context.Context, id uuid.UUID) (model.GigDetailDTO, error) {
	global.Logger.Info("Service: Calling database to get service by ID", zap.String("id", id.String()))
	gig, err := s.store.GetService(ctx, id)
	if err != nil {
		global.Logger.Error("Service: Database returned an error", zap.String("id", id.String()), zap.Error(err))
		return model.GigDetailDTO{}, err
	}

	global.Logger.Info("Service: Successfully retrieved service from database", zap.String("id", id.String()))
	return model.GigDetailDTO{
		Gig: model.Gig{
			ID:           gig.ID,
			UserID:       gig.UserID,
			Title:        gig.Title,
			Description:  gig.Description,
			Price:        gig.Price,
			DeliveryTime: gig.DeliveryTime,
			ImageURL:     utils.PtrIfValid(gig.ImageUrl),
			Status:       gig.Status,
			CreatedAt:    gig.CreatedAt,
		},
		SellerInfo: model.SellerInfo{
			UserID:         gig.UserID,
			UserName:       gig.UserName,
			UserProfilePic: utils.PtrIfValid(gig.UserProfilePic),
		},
		CategoryInfo: model.CategoryInfo{
			CategoryName: gig.CategoryName,
		},
	}, nil
}

// 4. Cập nhật dịch vụ
func (s *sGigService) UpdateService(ctx context.Context, input model.UpdateServiceParams) (model.Gig, error) {
	gig, err := s.store.UpdateService(ctx, database.UpdateServiceParams{
		ID:           input.ID,
		Title:        input.Title,
		Description:  input.Description,
		CategoryID:   input.CategoryID,
		Price:        input.Price,
		DeliveryTime: input.DeliveryTime,
		ImageUrl:     utils.ToNullString(input.ImageURL),
		Status:       input.Status,
	})
	if err != nil {
		return model.Gig{}, err
	}

	return model.Gig{
		ID:           gig.ID,
		UserID:       gig.UserID,
		Title:        gig.Title,
		Description:  gig.Description,
		CategoryID:   gig.CategoryID,
		Price:        gig.Price,
		DeliveryTime: gig.DeliveryTime,
		ImageURL:     utils.PtrIfValid(gig.ImageUrl),
		Status:       gig.Status,
		CreatedAt:    gig.CreatedAt,
	}, nil
}

// 5. Xoá dịch vụ
func (s *sGigService) DeleteService(ctx context.Context, id uuid.UUID) error {
	return s.store.DeleteService(ctx, id)
}
