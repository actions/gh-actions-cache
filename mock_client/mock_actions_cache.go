package mock_client

import (
	"github.com/actions/gh-actions-cache/types"
	"github.com/cli/go-gh/pkg/api"
	ghRepo "github.com/cli/go-gh/pkg/repository"
	// gh "github.com/cli/go-gh"
	"github.com/actions/gh-actions-cache/client"
	"net/url"
)

type MockArtifactCache struct{
	HttpClient	api.RESTClient
}

func NewMockArtifactCache() client.ArtifactCacheService {
	return &MockArtifactCache{nil}
}

func (a *MockArtifactCache) GetCacheUsage(repo ghRepo.Repository) float64 {
	return 4*1024*1024*1024
}

func (a *MockArtifactCache) ListCaches(repo ghRepo.Repository, queryParams url.Values) []types.CacheInfo {
	caches := []types.CacheInfo{types.CacheInfo{
		Key:            "123",
		Ref:            "123",
		LastAccessedAt: "123",
		Size:           123,
	},}
	return caches
}

func (a *MockArtifactCache) DeleteCaches(repo ghRepo.Repository, queryParams url.Values) int {
	return 2
}