package impl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"

	"Taskly.com/m/global"
	"go.uber.org/zap"

	"Taskly.com/m/internal/database"
	model "Taskly.com/m/internal/models"
	utils "Taskly.com/m/package/utils"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)

type sGigService struct {
	store database.Store
}

func NewGigService(store database.Store) *sGigService {
	return &sGigService{
		store: store,
	}
}

// 2. Tạo dịch vụ mới
func (s *sGigService) CreateService(ctx context.Context, input model.CreateServiceParams) (model.Gig, error) {
	var createdGig model.Gig

	err := s.store.ExecTx(ctx, func(q *database.Queries) error {
		// 1. Tạo Gig chính
		gig, err := q.CreateGig(ctx, database.CreateGigParams{
			UserID:      input.UserID,
			Title:       input.Title,
			Description: input.Description,
			CategoryID:  input.CategoryID,
			ImageUrl:    input.ImageURL,
			PricingMode: input.PricingMode,
			Status:      input.Status,
		})
		if err != nil {
			return fmt.Errorf("failed to create gig: %w", err)
		}

		createdGig = model.Gig{
			ID:          gig.ID,
			UserID:      gig.UserID,
			Title:       gig.Title,
			Description: gig.Description,
			CategoryID:  gig.CategoryID,
			ImageURL:    utils.PtrSliceIfValid(gig.ImageUrl),
			PricingMode: gig.PricingMode,
			Status:      gig.Status,
			CreatedAt:   gig.CreatedAt,
			UpdatedAt:   gig.UpdatedAt,
		}

		// 2. Tạo Gig Packages
		for _, pkg := range input.Packages {
			optionsJSON, err := json.Marshal(pkg.Options)
			if err != nil {
				return fmt.Errorf("failed to marshal package options: %w", err)
			}

			_, err = q.CreateGigPackage(ctx, database.CreateGigPackageParams{
				GigID:        gig.ID,
				Tier:         pkg.Tier,
				Price:        pkg.Price,
				DeliveryTime: pkg.DeliveryDays,
				Options:      pqtype.NullRawMessage{RawMessage: optionsJSON, Valid: true},
			})
			if err != nil {
				return fmt.Errorf("failed to create gig package: %w", err)
			}
		}

		// 3. Tạo Gig Requirements
		for _, req := range input.Requirements.Questions {
			_, err = q.CreateGigRequirement(ctx, database.CreateGigRequirementParams{
				GigID:    gig.ID,
				Question: req.Question,
				Required: req.Required,
			})
			if err != nil {
				return fmt.Errorf("failed to create gig requirement: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return model.Gig{}, err
	}

	return createdGig, nil
}

// GetServiceByID lấy chi tiết dịch vụ theo ID
func (s *sGigService) GetServiceByID(ctx context.Context, id uuid.UUID) (model.GigDetailDTO, error) {
	// Lấy thông tin Gig
	gig, err := s.store.GetService(ctx, id)
	if err != nil {
		global.Logger.Error(
			"Service: Database returned an error",
			zap.String("id", id.String()),
			zap.Error(err),
		)
		return model.GigDetailDTO{}, err
	}

	// Lấy Gig Packages
	dbPackages, err := s.store.GetGigPackagesByGigID(ctx, id)
	if err != nil {
		global.Logger.Error(
			"Service: Failed to get gig packages",
			zap.String("gig_id", id.String()),
			zap.Error(err),
		)
		return model.GigDetailDTO{}, fmt.Errorf("failed to get gig packages: %w", err)
	}

	var gigPackages []model.GigPackage
	for _, dbPkg := range dbPackages {
		var options model.GigPackageOptions
		if err := json.Unmarshal(dbPkg.Options.RawMessage, &options); err != nil {
			global.Logger.Error(
				"Service: Failed to unmarshal package options",
				zap.Error(err),
			)
			// Có thể bỏ qua lỗi này nếu muốn gig vẫn trả về mà không có options
		}

		gigPackages = append(gigPackages, model.GigPackage{
			Tier:         dbPkg.Tier,
			Price:        dbPkg.Price,
			DeliveryDays: dbPkg.DeliveryTime,
			Options:      options,
		})
	}

	// Lấy Gig Requirements
	dbRequirements, err := s.store.GetGigRequirementsByGigID(ctx, id)
	if err != nil {
		global.Logger.Error(
			"Service: Failed to get gig requirements",
			zap.String("gig_id", id.String()),
			zap.Error(err),
		)
		return model.GigDetailDTO{}, fmt.Errorf("failed to get gig requirements: %w", err)
	}

	var gigRequirements []model.Question
	for _, dbReq := range dbRequirements {
		gigRequirements = append(gigRequirements, model.Question{
			ID: 	  dbReq.ID,
			Question: dbReq.Question,
			Required: dbReq.Required,
		})
	}

	// Trả về kết quả
	global.Logger.Info(
		"Service: Successfully retrieved service from database",
		zap.String("id", id.String()),
	)

	return model.GigDetailDTO{
		Gig: model.Gig{
			ID:          gig.ID,
			UserID:      gig.UserID,
			Title:       gig.Title,
			Description: (gig.Description),
			CategoryID:  gig.CategoryID,
			ImageURL:    utils.PtrSliceIfValid(gig.ImageUrl),
			PricingMode: (gig.PricingMode),
			Status:      gig.Status,
			CreatedAt:   gig.CreatedAt,
			UpdatedAt:   gig.UpdatedAt,
		},
		SellerInfo: model.SellerInfo{
			UserName:       gig.UserName,
			UserProfilePic: utils.PtrStringIfValid(gig.UserProfilePic),
		},
		CategoryInfo: model.CategoryInfo{
			CategoryName: gig.CategoryName,
		},
		GigPackage: gigPackages,
		Question:   gigRequirements,
	}, nil
}

// 4. Cập nhật dịch vụ
func (s *sGigService) UpdateService(ctx context.Context, input model.UpdateServiceParams) (model.Gig, error) {
	gig, err := s.store.UpdateService(ctx, database.UpdateServiceParams{
		ID:         input.ID,
		Title:      input.Title,
		CategoryID: input.CategoryID,
		ImageUrl:   input.ImageURL,
		Status:     input.Status,
	})
	if err != nil {
		return model.Gig{}, err
	}

	return model.Gig{
		ID:         gig.ID,
		UserID:     gig.UserID,
		Title:      gig.Title,
		CategoryID: gig.CategoryID,
		ImageURL:   utils.PtrSliceIfValid(gig.ImageUrl),
		Status:     gig.Status,
		CreatedAt:  gig.CreatedAt,
	}, nil
}

// 5. Xoá dịch vụ
func (s *sGigService) DeleteService(ctx context.Context, id uuid.UUID) error {
	return s.store.DeleteService(ctx, id)
}

// 6. Lấy categories
func (s *sGigService) GetCategories(ctx context.Context) ([]model.GetCategoriesRs, error) {
	rows, err := s.store.GetCategories(ctx)
	var rs []model.GetCategoriesRs
	for _, row := range rows {
		rs = append(rs, model.GetCategoriesRs{
			ParentID:     row.ParentID,
			ParentName:   row.ParentName,
			ChildrenID:   utils.PtrIntIfValid(row.ChildrenID),
			ChildrenName: utils.PtrStringIfValid(row.ChildrenName),
		})
	}
	if err != nil {
		return []model.GetCategoriesRs{}, err
	}
	return rs, nil
}

func (s *sGigService) UploadGigImages(ctx context.Context, files []*multipart.FileHeader) ([]string, error) {
	if global.Cloudinary == nil {
		return nil, errors.New("Cloudinary service not initialized")
	}

	// Đặt tên thư mục trên Cloudinary cho ảnh gig. Bạn có thể thay đổi tùy ý.
	folder := "Taskly_gigs"

	urls, err := global.Cloudinary.UploadMultipleFiles(files, folder)
	if err != nil {
		return nil, fmt.Errorf("failed to upload gig images: %w", err)
	}
	return urls, nil
}

// SearchGigs tìm kiếm các gig dựa trên tiêu chí
func (s *sGigService) SearchGigs(ctx context.Context, params model.SearchGigParams) ([]model.SearchGigDTO, error) {
	var minPrice float64
	if params.MinPrice != nil {
		minPrice = *params.MinPrice
	}

	var maxPrice float64
	if params.MaxPrice != nil {
		maxPrice = *params.MaxPrice
	}

	dbParams := database.SearchGigsParams{
		SearchTerm:  params.SearchTerm,
		MinPrice:    minPrice,
		MaxPrice:    maxPrice,
		CategoryIds: params.CategoryIDs,
		LastGigID:   params.LastGigID,
	}
	fmt.Println("Service dbParams:", dbParams)

	rows, err := s.store.SearchGigs(ctx, dbParams)
	if err != nil {
		global.Logger.Error(
			"Service: Failed to search gigs",
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to search gigs: %w", err)
	}
	fmt.Println("Service rows from DB:", rows)

	var gigs []model.SearchGigDTO
	for _, row := range rows {
		gigs = append(gigs, model.SearchGigDTO{
			ID:          row.ID,
			Title:       row.Title,
			Description: row.Description,
			ImageURL:    utils.PtrSliceIfValid(row.ImageUrl),
			PricingMode: row.PricingMode,
			CreatedAt:   row.CreatedAt,
			UpdatedAt:   row.UpdatedAt,
			BasicPrice:  row.BasicPrice,
		})
	}

	return gigs, nil
}
