package internal

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRepo_IncorrectRepoString(t *testing.T) {
	r := "testOrg/testRepo/123/123"
	repo, err := GetRepo(r)

	assert.ErrorContains(t, err, fmt.Sprintf("expected the \"[HOST/]OWNER/REPO\" format, got \"%s\"", r))
	assert.Nil(t, repo)
}

func TestGetRepo_CorrectRepoString(t *testing.T) {
	r := "testOrg/testRepo"
	repo, err := GetRepo(r)

	assert.NoError(t, err)
	if assert.NotNil(t, repo) {
		assert.Equal(t, "github.com", repo.Host())
		assert.Equal(t, "testOrg", repo.Owner())
		assert.Equal(t, "testRepo", repo.Name())
	}
}

func TestGetRepo_CorrectRepoStringWithCustomHost(t *testing.T) {
	r := "api.testEnterprise.com/testOrg/testRepo"
	repo, err := GetRepo(r)

	assert.NoError(t, err)
	if assert.NotNil(t, repo) {
		assert.Equal(t, "api.testEnterprise.com", repo.Host())
		assert.Equal(t, "testOrg", repo.Owner())
		assert.Equal(t, "testRepo", repo.Name())
	}
}

func TestFormatCacheSize_MB(t *testing.T) {
	cacheSizeInBytes := 1024 * 1024 * 1.5
	cacheSizeDetailString := FormatCacheSize(cacheSizeInBytes)

	assert.Equal(t, "1.50 MB", cacheSizeDetailString)
}

func TestFormatCacheSize_GB(t *testing.T) {
	cacheSizeInBytes := 1024 * 1024 * 1024 * 1.5
	cacheSizeDetailString := FormatCacheSize(cacheSizeInBytes)

	assert.Equal(t, "1.50 GB", cacheSizeDetailString)
}
