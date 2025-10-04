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

type ImageUploadPayload struct {
	ShopID    uuid.UUID `json:"shop_id"`
	Type      string    `json:"type"` // "cover" atau "gallery"
	FileBytes []byte    `json:"file_bytes"`
	Folder    string    `json:"folder"`
	Filename  string    `json:"filename"`
}
type ImageEventUploadPayload struct {
	EventID   uuid.UUID `json:"event_id"`
	Type      string    `json:"type"`
	FileBytes []byte    `json:"file_bytes"`
	Folder    string    `json:"folder"`
	Filename  string    `json:"filename"`
}
