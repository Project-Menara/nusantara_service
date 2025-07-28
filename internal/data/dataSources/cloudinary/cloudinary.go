package cloudinary

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryService interface {
	UploadImage(ctx context.Context, file *multipart.FileHeader, folder string) (string, error)
	DestroyImage(ctx context.Context, publicID string) error
}

type cloudinaryServiceImpl struct {
	cld *cloudinary.Cloudinary
}

func NewCloudinaryService() (CloudinaryService, error) {
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	if cloudName == "" || apiKey == "" || apiSecret == "" {
		return nil, errors.New("cloudinary credentials are not set in environment variables")
	}

	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize cloudinary: %w", err)
	}

	return &cloudinaryServiceImpl{cld: cld}, nil
}

// UploadImage implements CloudinaryService.
func (c *cloudinaryServiceImpl) UploadImage(ctx context.Context, file *multipart.FileHeader, folder string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}

	defer src.Close()

	uploadResult, err := c.cld.Upload.Upload(ctx, src, uploader.UploadParams{
		Folder: folder,
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload file to cloudinary: %w", err)
	}

	return uploadResult.SecureURL, nil

}

// DestroyImage implements CloudinaryService.
func (c *cloudinaryServiceImpl) DestroyImage(ctx context.Context, publicID string) error {
	invalidate := true
	_, err := c.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID:   publicID,
		Invalidate: &invalidate,
	})
	if err != nil {
		return fmt.Errorf("failed to delete image from cloudinary: %w", err)
	}
	return nil
}
