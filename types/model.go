package types

type RepoLevelUsageApiResponse struct {
	FullName               string  `json:"full_name"`
	ActiveCacheSizeInBytes float64 `json:"active_caches_size_in_bytes"`
	ActiveCacheCount       float64 `json:"active_caches_count"`
}

type ListApiResponse struct {
	TotalCount    int            `json:"total_count"`
	ActionsCaches []ActionsCache `json:"actions_caches"`
}

type DeleteApiResponse struct {
	TotalCount    int            `json:"total_count"`
	ActionsCaches []ActionsCache `json:"actions_caches"`
}

type ActionsCache struct {
	Id             int     `json:"id"`
	Ref            string  `json:"ref"`
	Key            string  `json:"key"`
	Version        string  `json:"version"`
	LastAccessedAt string  `json:"last_accessed_at"`
	CreatedAt      string  `json:"created_at"`
	SizeInBytes    float64 `json:"size_in_bytes"`
}

type InputFlags struct {
	Repo    string
	Branch  string
	Limit   int
	Key     string
	Order   string
	Sort    string
	Confirm bool
}
