package utils

import (
	"net/url"
	"path"
	"strings"
)

func ExtractPublicIDFromCloudinaryURL(imageURL string) string {
	parts := strings.Split(imageURL, "/upload/")
	if len(parts) != 2 {
		return ""
	}

	// Ambil bagian setelah /upload/
	afterUpload := parts[1]

	// Hilangkan versi (v12345678/...)
	slashIndex := strings.Index(afterUpload, "/")
	if slashIndex == -1 {
		return ""
	}

	// Ambil path tanpa versi
	pathWithFile := afterUpload[slashIndex+1:]

	// Decode URL-encoded string (e.g., %20 → space)
	decodedPath, err := url.PathUnescape(pathWithFile)
	if err != nil {
		return ""
	}

	// Hilangkan ekstensi file (.webp, .png, .jpg, dll)
	publicID := strings.TrimSuffix(decodedPath, path.Ext(decodedPath))
	return publicID
}
