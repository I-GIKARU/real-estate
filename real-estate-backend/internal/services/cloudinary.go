package services

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"real-estate-backend/internal/config"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/google/uuid"
)

// CloudinaryService handles image uploads to Cloudinary
type CloudinaryService struct {
	client *cloudinary.Cloudinary
	config *config.CloudinaryConfig
}

// UploadResponse represents the response from Cloudinary upload
type UploadResponse struct {
	PublicID   string `json:"public_id"`
	URL        string `json:"url"`
	SecureURL  string `json:"secure_url"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	Format     string `json:"format"`
	ResourceType string `json:"resource_type"`
	Bytes      int    `json:"bytes"`
}

// NewCloudinaryService creates a new Cloudinary service
func NewCloudinaryService(cfg *config.CloudinaryConfig) (*CloudinaryService, error) {
	if cfg.CloudName == "" || cfg.APIKey == "" || cfg.APISecret == "" {
		return nil, fmt.Errorf("cloudinary configuration is incomplete")
	}

	cld, err := cloudinary.NewFromParams(cfg.CloudName, cfg.APIKey, cfg.APISecret)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize cloudinary: %w", err)
	}

	return &CloudinaryService{
		client: cld,
		config: cfg,
	}, nil
}

// UploadImage uploads an image file to Cloudinary
func (s *CloudinaryService) UploadImage(ctx context.Context, file multipart.File, header *multipart.FileHeader, propertyID uuid.UUID) (*UploadResponse, error) {
	// Generate unique public ID
	publicID := s.generatePublicID(propertyID, header.Filename)

	// Upload parameters
	invalidate := true
	uploadParams := uploader.UploadParams{
		PublicID:        publicID,
		Folder:          s.config.Folder,
		ResourceType:    "image",
		Transformation:  "q_auto,f_auto", // Auto quality and format optimization
		AllowedFormats:  []string{"jpg", "jpeg", "png", "webp"},
		Invalidate:      &invalidate, // Invalidate CDN cache
		Tags:            []string{"property", "real-estate", "kenya"},
	}

	// Upload the file
	result, err := s.client.Upload.Upload(ctx, file, uploadParams)
	if err != nil {
		return nil, fmt.Errorf("failed to upload image to cloudinary: %w", err)
	}

	return &UploadResponse{
		PublicID:     result.PublicID,
		URL:          result.URL,
		SecureURL:    result.SecureURL,
		Width:        result.Width,
		Height:       result.Height,
		Format:       result.Format,
		ResourceType: result.ResourceType,
		Bytes:        result.Bytes,
	}, nil
}

// DeleteImage deletes an image from Cloudinary
func (s *CloudinaryService) DeleteImage(ctx context.Context, publicID string) error {
	invalidate := true
	_, err := s.client.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID:     publicID,
		ResourceType: "image",
		Invalidate:   &invalidate,
	})
	if err != nil {
		return fmt.Errorf("failed to delete image from cloudinary: %w", err)
	}
	return nil
}

// GenerateTransformationURL generates a transformed image URL
func (s *CloudinaryService) GenerateTransformationURL(publicID string, transformation string) string {
	img, _ := s.client.Image(publicID)
	transformedURL, _ := img.String()
	return transformedURL
}

// GetImageThumbnail generates a thumbnail URL for an image
func (s *CloudinaryService) GetImageThumbnail(publicID string, width, height int) string {
	transformation := fmt.Sprintf("w_%d,h_%d,c_fill,q_auto,f_auto", width, height)
	return s.GenerateTransformationURL(publicID, transformation)
}

// generatePublicID creates a unique public ID for the image
func (s *CloudinaryService) generatePublicID(propertyID uuid.UUID, filename string) string {
	// Remove extension from filename
	name := strings.TrimSuffix(filename, filepath.Ext(filename))
	
	// Clean the filename
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ToLower(name)
	
	// Generate timestamp
	timestamp := time.Now().Unix()
	
	// Create unique public ID
	return fmt.Sprintf("property_%s_%s_%d", propertyID.String(), name, timestamp)
}

// ValidateImageFile validates if the uploaded file is a valid image
func (s *CloudinaryService) ValidateImageFile(header *multipart.FileHeader, allowedTypes []string) error {
	// Check file size
	if header.Size > 10485760 { // 10MB
		return fmt.Errorf("file size exceeds maximum limit of 10MB")
	}

	// Check content type
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		return fmt.Errorf("content type not specified")
	}

	// Validate against allowed types
	allowed := false
	for _, allowedType := range allowedTypes {
		if contentType == allowedType {
			allowed = true
			break
		}
	}

	if !allowed {
		return fmt.Errorf("file type %s not allowed. Allowed types: %v", contentType, allowedTypes)
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowedExts := []string{".jpg", ".jpeg", ".png", ".webp"}
	extAllowed := false
	for _, allowedExt := range allowedExts {
		if ext == allowedExt {
			extAllowed = true
			break
		}
	}

	if !extAllowed {
		return fmt.Errorf("file extension %s not allowed. Allowed extensions: %v", ext, allowedExts)
	}

	return nil
}
