package payload

type ImageDeleteTask struct {
	PublicIDs []string `json:"public_ids"`
}

type CacheInvalidateTask struct {
	Keys []string `json:"keys"`
}
