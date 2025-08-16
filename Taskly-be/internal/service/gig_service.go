package service

import (
	"context"
	"mime/multipart"

	model "Taskly.com/m/internal/models"
	"github.com/google/uuid"
)

type IGigService interface {
	// 1. Lấy danh sách dịch vụ (có lọc, tìm kiếm)
	// ListServices(ctx context.Context, params model.ListServicesParams) ([]model.Gig, error)

	// 2. Tạo dịch vụ mới
	CreateService(ctx context.Context, input model.CreateServiceParams) (model.Gig, error)

	// 3. Lấy chi tiết dịch vụ
	GetServiceByID(ctx context.Context, id uuid.UUID) (model.GigDetailDTO, error)

	// 4. Cập nhật dịch vụ
	UpdateService(ctx context.Context, input model.UpdateServiceParams) (model.Gig, error)

	// 5. Xoá dịch vụ
	DeleteService(ctx context.Context, id uuid.UUID) error

	// 6.Lấy danh mục
	GetCategories(cxt context.Context) ([]model.GetCategoriesRs, error)

	// 7. Upload nhiều ảnh cho Gig
	UploadGigImages(ctx context.Context, files []*multipart.FileHeader) ([]string, error)

	// ... (các hàm hiện có) ...

	// 8. Tìm kiếm dịch vụ theo tiêu đề hoặc mô tả
	SearchGigs(ctx context.Context, params model.SearchGigParams) ([]model.SearchGigDTO, error)
}

var (
	localGigService IGigService
)

func GetGigService() IGigService {
	if localGigService == nil {
		panic("implement localGigService not found")
	}
	return localGigService
}

func InitGigService(i IGigService) {
	localGigService = i
}
