package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	MaxFileSize      = 10 << 20 // 10 MB
	MaxImagesPerVehicle = 8
)

var allowedImageExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
}

// UploadVehicleImages handles uploading multiple vehicle images
func UploadVehicleImages(files []*multipart.FileHeader, uploadDir string) ([]string, error) {
	if len(files) > MaxImagesPerVehicle {
		return nil, fmt.Errorf("maximum %d images allowed", MaxImagesPerVehicle)
	}

	// Create upload directory if it doesn't exist
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %v", err)
	}

	var uploadedPaths []string

	for _, fileHeader := range files {
		// Validate file size
		if fileHeader.Size > MaxFileSize {
			return uploadedPaths, fmt.Errorf("file %s exceeds maximum size of 10MB", fileHeader.Filename)
		}

		// Validate file extension
		ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
		if !allowedImageExtensions[ext] {
			return uploadedPaths, fmt.Errorf("file %s has invalid extension. Allowed: jpg, jpeg, png, webp", fileHeader.Filename)
		}

		// Open uploaded file
		file, err := fileHeader.Open()
		if err != nil {
			return uploadedPaths, fmt.Errorf("failed to open file %s: %v", fileHeader.Filename, err)
		}
		defer file.Close()

		// Generate unique filename
		filename := generateUniqueFilename(ext)
		filepath := filepath.Join(uploadDir, filename)

		// Create destination file
		dst, err := os.Create(filepath)
		if err != nil {
			return uploadedPaths, fmt.Errorf("failed to create file %s: %v", filename, err)
		}
		defer dst.Close()

		// Copy file content
		if _, err := io.Copy(dst, file); err != nil {
			return uploadedPaths, fmt.Errorf("failed to save file %s: %v", filename, err)
		}

		// Store relative path for database
		relativePath := fmt.Sprintf("/uploads/vehicles/%s", filename)
		uploadedPaths = append(uploadedPaths, relativePath)
	}

	return uploadedPaths, nil
}

// generateUniqueFilename creates a unique filename using UUID and timestamp
func generateUniqueFilename(extension string) string {
	timestamp := time.Now().Unix()
	uniqueID := uuid.New().String()
	return fmt.Sprintf("%d_%s%s", timestamp, uniqueID, extension)
}

// DeleteFile removes a file from the filesystem
func DeleteFile(filePath string) error {
	// Convert relative path to absolute
	absPath := filepath.Join(".", filePath)
	return os.Remove(absPath)
}
