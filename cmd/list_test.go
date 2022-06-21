package cmd

import (
	"fmt"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/cli/go-gh/pkg/api"
	"github.com/actions/gh-actions-cache/mock_client"
)

func TestList_WithNoFlag(t *testing.T) {
	opts := api.ClientOptions{
		Headers: map[string]string{"User-Agent": fmt.Sprintf("gh-actions-cache/%s/%s", "0.0.1", "list")},
	}
	artifactCache := mock_client.NewMockArtifactCache()
	cmd := NewCmdList(opts, artifactCache)
	fmt.Println(cmd)
	assert.Equal(t, "v1.0.0", "v1.0.0")
}
