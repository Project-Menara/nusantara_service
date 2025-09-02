package payload

import (
	"github.com/google/uuid"
)

type ImageDeleteTask struct {
	PublicIDs []string `json:"public_ids"`
}

type CacheInvalidateTask struct {
	Keys []string `json:"keys"`
}

type ImageSendTask struct {
	UserID    uuid.UUID `json:"user_id"`
	FileBytes []byte    `json:"file_bytes"`
	Folder    string    `json:"folder"`
	Filename  string    `json:"filename"`
}
