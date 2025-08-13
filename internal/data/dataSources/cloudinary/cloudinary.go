// infra/cloudinary/service.go
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

type UploadResult struct {
	URL      string
	PublicID string
}

type CloudinaryService interface {
	UploadImage(ctx context.Context, file *multipart.FileHeader, folder, filename string) (*UploadResult, error)
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
		return nil, errors.New("cloudinary credentials are not set")
	}
	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize cloudinary: %w", err)
	}
	return &cloudinaryServiceImpl{cld: cld}, nil
}

func (c *cloudinaryServiceImpl) UploadImage(ctx context.Context, file *multipart.FileHeader, folder, filename string) (*UploadResult, error) {
	f, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	publicID := fmt.Sprintf("%s/%s", folder, filename)
	resp, err := c.cld.Upload.Upload(ctx, f, uploader.UploadParams{
		PublicID:     publicID,
		Folder:       folder,
		Overwrite:    boolPtr(true),
		ResourceType: "image",
	})
	if err != nil {
		return nil, err
	}
	return &UploadResult{URL: resp.SecureURL, PublicID: resp.PublicID}, nil
}

func (c *cloudinaryServiceImpl) DestroyImage(ctx context.Context, publicID string) error {
	resp, err := c.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID:   publicID,
		Invalidate: boolPtr(true),
	})
	if err != nil {
		return fmt.Errorf("failed to delete image: %w", err)
	}
	if resp.Result != "ok" {
		return fmt.Errorf("deletion result: %s", resp.Result)
	}
	return nil
}

func boolPtr(b bool) *bool { return &b }
