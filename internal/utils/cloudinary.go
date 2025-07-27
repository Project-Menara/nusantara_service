package utils

import "strings"

func ExtractPublicIDFromCloudinaryURL(url string) string {
	parts := strings.Split(url, "/upload/")
	if len(parts) < 2 {
		return ""
	}

	pathAfterUpload := parts[1]
	if len(pathAfterUpload) > 1 && pathAfterUpload[0] == 'v' {
		if idx := strings.Index(pathAfterUpload, "/"); idx != -1 {
			pathAfterUpload = pathAfterUpload[idx+1:]
		}
	}

	if idx := strings.LastIndex(pathAfterUpload, "?"); idx != -1 {
		pathAfterUpload = pathAfterUpload[:idx]
	}

	return pathAfterUpload
}
