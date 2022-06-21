package cmd

import "time"

type ListApiResponse struct {
	Total_count    int            `json:"total_count"`
	Actions_caches []ActionsCache `json:"actions_caches"`
}

type ActionsCache struct {
	Id             int       `json:"id"`
	Ref            string    `json:"ref"`
	Key            string    `json:"key"`
	Version        string    `json:"version"`
	LastAccessedAt time.Time `json:"last_accessed_at"`
	CreatedAt      time.Time `json:"created_at"`
	SizeInBytes    uint64    `json:"size_in_bytes"`
}
