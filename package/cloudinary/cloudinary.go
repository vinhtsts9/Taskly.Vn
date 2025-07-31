package cloudinary

import (
	"context"
	"errors"
	"fmt"
	"io"

	"bytes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// CloudinaryService là dịch vụ để tương tác với Cloudinary
type CloudinaryService struct {
	cld *cloudinary.Cloudinary
}

// Hàm khởi tạo Cloudinary
func InitCloudinary(cloudName, apiKey, apiSecret string) (*CloudinaryService, error) {
	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return nil, fmt.Errorf("Failed to init cloudinary")
	}
	return &CloudinaryService{
		cld: cld,
	}, nil
}

// Hàm upload ảnh từ URL
func (c *CloudinaryService) UploadImageFromURLToCloudinary(imageUrl string) (string, error) {
	// Kiểm tra Cloudinary đã được khởi tạo hay chưa
	if c.cld == nil {
		return "", errors.New("cloudinary not initialized")
	}

	// Upload từ URL
	resp, err := c.cld.Upload.Upload(context.Background(), imageUrl, uploader.UploadParams{})
	if err != nil {
		return "", err
	}
	return resp.SecureURL, nil
}

// Hàm upload ảnh từ file
func (c *CloudinaryService) UploadImageToCloudinaryFromReader(file io.Reader, folder string) (string, error) {
	// Kiểm tra xem Cloudinary đã được khởi tạo hay chưa
	if c.cld == nil {
		return "", fmt.Errorf("cloudinary not initialized")
	}

	// Upload file từ io.Reader lên Cloudinary
	resp, err := c.cld.Upload.Upload(context.Background(), file, uploader.UploadParams{
		Folder: folder, // Chọn thư mục lưu trữ trên Cloudinary (tùy chỉnh theo yêu cầu)
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload image to Cloudinary: %w", err)
	}

	// Trả về URL của ảnh đã upload
	return resp.SecureURL, nil
}

// Hàm upload nhiều ảnh từ file local lên Cloudinary
func (c *CloudinaryService) UploadImagesFromLocal(files []string, folderName string) ([]string, error) {
	if c.cld == nil {
		return nil, errors.New("cloudinary not initialized")
	}

	var urls []string
	for _, filePath := range files {
		resp, err := c.cld.Upload.Upload(context.Background(), filePath, uploader.UploadParams{
			Folder: folderName,
		})
		if err != nil {
			return nil, err
		}
		urls = append(urls, resp.SecureURL)
	}

	return urls, nil
}

// Hàm upload ảnh từ S3
func (c *CloudinaryService) UploadImageFromS3ToCloudinary(bucketName, imageName string) (string, error) {
	// Kiểm tra Cloudinary đã được khởi tạo hay chưa
	if c.cld == nil {
		return "", errors.New("cloudinary not initialized")
	}

	// Tạo session S3
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"), // Thay đổi theo vùng của bạn
	})
	if err != nil {
		return "", err
	}

	svc := s3.New(sess)

	// Lấy đối tượng từ S3
	result, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(imageName),
	})
	if err != nil {
		return "", err
	}
	defer result.Body.Close()

	// Đọc nội dung đối tượng
	buf := new(bytes.Buffer)
	buf.ReadFrom(result.Body)

	// Upload lên Cloudinary
	resp, err := c.cld.Upload.Upload(context.Background(), buf, uploader.UploadParams{})
	if err != nil {
		return "", err
	}
	return resp.SecureURL, nil
}
