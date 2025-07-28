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
	UploadImage(ctx context.Context, file *multipart.FileHeader, folder string, filename string) (string, error)
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
func (c *cloudinaryServiceImpl) UploadImage(ctx context.Context, file *multipart.FileHeader, folder string, filename string) (string, error) {
	fileContent, err := file.Open()
	if err != nil {
		return "", err
	}
	defer fileContent.Close()

	publicID := fmt.Sprintf("%s/%s", folder, filename)

	uploadParam := uploader.UploadParams{
		PublicID:     publicID,
		Folder:       folder, // folder sudah di PublicID
		Overwrite:    boolPtr(true),
		ResourceType: "image",
	}

	resp, err := c.cld.Upload.Upload(ctx, fileContent, uploadParam)
	if err != nil {
		return "", err
	}
	return resp.SecureURL, nil

}

// DestroyImage implements CloudinaryService.
func (c *cloudinaryServiceImpl) DestroyImage(ctx context.Context, publicID string) error {
	resp, err := c.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID:   publicID,
		Invalidate: boolPtr(true),
	})
	if err != nil {
		fmt.Printf("Error deleting image from cloudinary: %v, PublicID: %s\n", err, publicID) // Added logging
		return fmt.Errorf("failed to delete image from cloudinary: %w", err)
	}
	if resp.Result != "ok" { // Cloudinary indicates success with "ok"
		fmt.Printf("Cloudinary destroy result not 'ok': %s, PublicID: %s\n", resp.Result, publicID) // Added logging for non-ok result
		return fmt.Errorf("cloudinary image deletion failed with result: %s", resp.Result)
	}
	fmt.Println("Deleting public ID:", publicID)
	return nil
}
func boolPtr(b bool) *bool {
	return &b
}
