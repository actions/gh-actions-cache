package internal

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRepo_IncorrectRepoString(t *testing.T) {
	r := "testOrg/testRepo/123/123"
	repo, err := GetRepo(r)

	assert.Error(t, err)
	assert.Nil(t, repo)
	assert.Equal(t, err.Error(), fmt.Sprintf("expected the \"[HOST/]OWNER/REPO\" format, got \"%s\"", r))
}

func TestGetRepo_CorrectRepoString(t *testing.T) {
	r := "testOrg/testRepo"
	repo, err := GetRepo(r)

	assert.NotNil(t, repo)
	assert.Nil(t, err)
	assert.Equal(t, repo.Host(), "github.com")
	assert.Equal(t, repo.Owner(), "testOrg")
	assert.Equal(t, repo.Name(), "testRepo")
}

func TestGetRepo_CorrectRepoStringWithCustomHost(t *testing.T) {
	r := "api.testEnterprise.com/testOrg/testRepo"
	repo, err := GetRepo(r)

	assert.NotNil(t, repo)
	assert.Nil(t, err)
	assert.Equal(t, repo.Host(), "api.testEnterprise.com")
	assert.Equal(t, repo.Owner(), "testOrg")
	assert.Equal(t, repo.Name(), "testRepo")
}

func TestFormatCacheSize_MB(t *testing.T) {
	cacheSizeInBytes := 1024 * 1024 * 1.5
	cacheSizeDetailString := FormatCacheSize(cacheSizeInBytes)

	assert.NotNil(t, cacheSizeDetailString)
	assert.Equal(t, cacheSizeDetailString, "1.50 MB")
}

func TestFormatCacheSize_GB(t *testing.T) {
	cacheSizeInBytes := 1024 * 1024 * 1024 * 1.5
	cacheSizeDetailString := FormatCacheSize(cacheSizeInBytes)

	assert.NotNil(t, cacheSizeDetailString)
	assert.Equal(t, cacheSizeDetailString, "1.50 GB")
}

func TestGenerateQueryParams_DefaultFlags(t *testing.T) {
	branch := ""
	limit := 30
	key := ""
	order := ""
	sort := ""

	queryParams := GenerateQueryParams(branch, limit, key, order, sort)

	assert.Equal(t, queryParams.Encode(), "")
}

func TestGenerateQueryParams_NonDefaultFlag(t *testing.T) {
	branch := "main"
	limit := 30
	key := ""
	order := ""
	sort := ""

	queryParams := GenerateQueryParams(branch, limit, key, order, sort)

	assert.Equal(t, queryParams.Encode(), "ref=refs%2Fheads%2Fmain")
}
