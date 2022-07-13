package cmd

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/actions/gh-actions-cache/internal"
	"github.com/actions/gh-actions-cache/types"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestListWithIncorrectArguments(t *testing.T) {
	t.Cleanup(gock.Off)

	cmd := NewCmdList()
	cmd.SetArgs([]string{"keyValue"})
	err := cmd.Execute()

	assert.NotNil(t, err)
	assert.Equal(t, err, fmt.Errorf("Invalid argument(s). Expected 0 received 1"))
	assert.True(t, gock.IsDone(), internal.PrintPendingMocks(gock.Pending()))
}

func TestListWithIncorrectRepo(t *testing.T) {
	t.Cleanup(gock.Off)

	cmd := NewCmdList()
	cmd.SetArgs([]string{"--repo", "testOrg/testRepo/123/123"})
	err := cmd.Execute()

	assert.NotNil(t, err)
	assert.Equal(t, err, fmt.Errorf("expected the \"[HOST/]OWNER/REPO\" format, got \"testOrg/testRepo/123/123\""))
	assert.True(t, gock.IsDone(), internal.PrintPendingMocks(gock.Pending()))
}

func TestListWithNegativeLimit(t *testing.T) {
	t.Cleanup(gock.Off)

	cmd := NewCmdList()
	cmd.SetArgs([]string{"--limit", "-1", "--repo", "testOrg/testRepo"})
	err := cmd.Execute()

	assert.NotNil(t, err)
	assert.Equal(t, err, fmt.Errorf("-1 is not a valid integer value for limit flag. Allowed values: 1-100"))
	assert.True(t, gock.IsDone(), internal.PrintPendingMocks(gock.Pending()))
}

func TestListWithIncorrectLimit(t *testing.T) {
	t.Cleanup(gock.Off)

	cmd := NewCmdList()
	cmd.SetArgs([]string{"--limit", "101", "--repo", "testOrg/testRepo"})
	err := cmd.Execute()

	assert.NotNil(t, err)
	assert.Equal(t, err, fmt.Errorf("101 is not a valid integer value for limit flag. Allowed values: 1-100"))
	assert.True(t, gock.IsDone(), internal.PrintPendingMocks(gock.Pending()))
}

func TestListWithIncorrectOrder(t *testing.T) {
	t.Cleanup(gock.Off)

	cmd := NewCmdList()
	cmd.SetArgs([]string{"--order", "incorrectOrderValue", "--repo", "testOrg/testRepo"})
	err := cmd.Execute()

	assert.NotNil(t, err)
	assert.Equal(t, err, fmt.Errorf("incorrectOrderValue is not a valid value for order flag. Allowed values: asc/desc"))
	assert.True(t, gock.IsDone(), internal.PrintPendingMocks(gock.Pending()))
}

func TestListWithIncorrectSort(t *testing.T) {
	t.Cleanup(gock.Off)

	cmd := NewCmdList()
	cmd.SetArgs([]string{"--sort", "incorrectSortValue", "--repo", "testOrg/testRepo"})
	err := cmd.Execute()

	assert.NotNil(t, err)
	assert.Equal(t, err, fmt.Errorf("incorrectSortValue is not a valid value for sort flag. Allowed values: last-used/size/created-at"))
	assert.True(t, gock.IsDone(), internal.PrintPendingMocks(gock.Pending()))
}

func TestListWithIncorrectRepoForListCaches(t *testing.T) {
	t.Cleanup(gock.Off)
	gock.New("https://api.github.com").
		Get("/repos/testOrg/testRepo/actions/cache/usage").
		Reply(200).
		JSON(`{
			"full_name": "t-dedah/vipul-bugbash",
			"active_caches_size_in_bytes": 291205,
			"active_caches_count": 12
		}`)

	gock.New("https://api.github.com").
		Get("/repos/testOrg/testRepo/actions/caches").
		Reply(404).
		JSON(`{
			"message": "Not Found",
			"documentation_url": "https://docs.github.com/rest/reference/actions#get-github-actions-cache-list-for-a-repository"
		}`)

	cmd := NewCmdList()
	cmd.SetArgs([]string{"--repo", "testOrg/testRepo"})
	err := cmd.Execute()

	assert.NotNil(t, err)
	assert.Equal(t, reflect.TypeOf(err), reflect.TypeOf(types.HandledError{}))

	var customError types.HandledError
	errors.As(err, &customError)
	assert.Equal(t, customError.Message, "The given repo does not exist.")

	assert.True(t, gock.IsDone(), internal.PrintPendingMocks(gock.Pending()))
}

func TestListWithUnauthorizedRequestForListCaches(t *testing.T) {
	t.Cleanup(gock.Off)
	gock.New("https://api.github.com").
		Get("/repos/testOrg/testRepo/actions/cache/usage").
		Reply(200).
		JSON(`{
			"full_name": "t-dedah/vipul-bugbash",
			"active_caches_size_in_bytes": 291205,
			"active_caches_count": 12
		}`)

	gock.New("https://api.github.com").
		Get("/repos/testOrg/testRepo/actions/caches").
		Reply(401).
		JSON(`{
			"message": "Must have admin rights to Repository.",
			"documentation_url": "https://docs.github.com/rest/actions/cache#delete-a-github-actions-cache-for-a-repository-using-a-cache-id"
		}`)

	cmd := NewCmdList()
	cmd.SetArgs([]string{"--repo", "testOrg/testRepo"})
	err := cmd.Execute()

	assert.NotNil(t, err)
	assert.Equal(t, reflect.TypeOf(err), reflect.TypeOf(types.HandledError{}))

	var customError types.HandledError
	errors.As(err, &customError)
	assert.Equal(t, customError.Message, "Must have admin rights to Repository.")

	assert.True(t, gock.IsDone(), internal.PrintPendingMocks(gock.Pending()))
}

func TestListWithInternalServerErrorForListCaches(t *testing.T) {
	t.Cleanup(gock.Off)
	gock.New("https://api.github.com").
		Get("/repos/testOrg/testRepo/actions/cache/usage").
		Reply(200).
		JSON(`{
			"full_name": "t-dedah/vipul-bugbash",
			"active_caches_size_in_bytes": 291205,
			"active_caches_count": 12
		}`)

	gock.New("https://api.github.com").
		Get("/repos/testOrg/testRepo/actions/caches").
		Reply(500).
		JSON(`{
			"message": "Internal Server Error",
			"documentation_url": "https://docs.github.com/rest/reference/actions#get-github-actions-cache-list-for-a-repository"
		}`)

	cmd := NewCmdList()
	cmd.SetArgs([]string{"--repo", "testOrg/testRepo"})
	err := cmd.Execute()

	assert.NotNil(t, err)
	assert.Equal(t, reflect.TypeOf(err), reflect.TypeOf(types.HandledError{}))

	var customError types.HandledError
	errors.As(err, &customError)
	assert.Equal(t, customError.Message, "We could not process your request due to internal error.")

	assert.True(t, gock.IsDone(), internal.PrintPendingMocks(gock.Pending()))
}

func TestListSuccess(t *testing.T) {
	t.Cleanup(gock.Off)
	gock.New("https://api.github.com").
		Get("/repos/testOrg/testRepo/actions/cache/usage").
		Reply(200).
		JSON(`{
			"full_name": "t-dedah/vipul-bugbash",
			"active_caches_size_in_bytes": 2432967,
			"active_caches_count": 1
		}`)

	gock.New("https://api.github.com").
		Get("/repos/testOrg/testRepo/actions/caches").
		Reply(200).
		JSON(`{
			"total_count": 1,
			"actions_caches": [
				{
					"id": 29,
					"ref": "refs/heads/master",
					"key": "Linux-build-cache-node-modules-3fd22dd3a926d576e2562e8b76a5ff157cd3b986f3d44195acfe7efa6bc05919-8",
					"version": "7fcda33c1e1d849a13bcc06f49b9ab64efc01ca9dabe4d7a8d0d387feef4fc88",
					"last_accessed_at": "2022-06-22T20:32:45.550000000Z",
					"created_at": "2022-06-22T20:32:45.550000000Z",
					"size_in_bytes": 2432967
				}]
			}`)

	cmd := NewCmdList()
	cmd.SetArgs([]string{"--repo", "testOrg/testRepo"})
	err := cmd.Execute()

	assert.Nil(t, err)
	assert.True(t, gock.IsDone(), internal.PrintPendingMocks(gock.Pending()))
}
