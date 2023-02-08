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

func TestDeleteWithIncorrectArguments(t *testing.T) {
	t.Cleanup(gock.Off)

	cmd := NewCmdDelete()
	cmd.SetArgs([]string{})
	err := cmd.Execute()

	assert.NotNil(t, err)
	assert.Equal(t, err, fmt.Errorf("accepts 1 arg(s), received 0"))
	assert.True(t, gock.IsDone(), internal.PrintPendingMocks(gock.Pending()))
}

func TestDeleteWithIncorrectRepo(t *testing.T) {
	t.Cleanup(gock.Off)

	cmd := NewCmdDelete()
	cmd.SetArgs([]string{"--repo", "testOrg/testRepo/123/123", "cacheName"})
	err := cmd.Execute()

	assert.NotNil(t, err)
	assert.Equal(t, err, fmt.Errorf("expected the \"[HOST/]OWNER/REPO\" format, got \"testOrg/testRepo/123/123\""))
	assert.True(t, gock.IsDone(), internal.PrintPendingMocks(gock.Pending()))
}

func TestDeleteWithIncorrectRepoForDeleteCaches(t *testing.T) {
	t.Cleanup(gock.Off)

	gock.New("https://api.github.com").
		Get("/repos/testOrg/testRepo/actions/caches").
		MatchParam("key", "cacheName").
		Reply(404).
		JSON(`{
			"message": "Not Found",
			"documentation_url": "https://docs.github.com/rest/actions/cache#list-github-actions-caches-for-a-repository"
		}`)

	cmd := NewCmdDelete()
	cmd.SetArgs([]string{"--repo", "testOrg/testRepo", "cacheName"})
	err := cmd.Execute()

	assert.NotNil(t, err)
	assert.True(t, gock.IsDone(), internal.PrintPendingMocks(gock.Pending()))
	var customError types.HandledError
	errors.As(err, &customError)
	assert.Equal(t, customError.Message, "The given repo does not exist.")
}

func TestDeleteSuccessWithConfirmFlagProvided(t *testing.T) {
	t.Cleanup(gock.Off)

	gock.New("https://api.github.com").
		Delete("/repos/testOrg/testRepo/actions/caches").
		MatchParam("key", "2022-06-29T13:33:49").
		Reply(200).
		JSON(`{
				"total_count": 1,
				"actions_caches": [
					{
						"id": 1293,
						"ref": "refs/heads/main",
						"key": "2022-06-29T13:33:49",
						"version": "803758043e242677f6b8650742372d82ded436d99b2a8a09bc3b6ed77cd6aec2",
						"last_accessed_at": "2022-06-29T13:33:52.280000000Z",
						"created_at": "2022-06-29T13:33:52.280000000Z",
						"size_in_bytes": 29747
					}
				]
			}`)

	cmd := NewCmdDelete()
	cmd.SetArgs([]string{"--repo", "testOrg/testRepo", "2022-06-29T13:33:49", "--confirm"})
	err := cmd.Execute()

	assert.NoError(t, err)
	assert.True(t, gock.IsDone(), internal.PrintPendingMocks(gock.Pending()))
}

func TestDeleteFailureWhileTakingUserInput(t *testing.T) {
	t.Cleanup(gock.Off)
	choice = "Delete"

	gock.New("https://api.github.com").
		Get("/repos/testOrg/testRepo/actions/caches").
		MatchParam("key", "2022-06-29T13:33:49").
		Reply(200).
		JSON(`{
				"total_count": 1,
				"actions_caches": [
					{
						"id": 1293,
						"ref": "refs/heads/main",
						"key": "2022-06-29T13:33:49",
						"version": "803758043e242677f6b8650742372d82ded436d99b2a8a09bc3b6ed77cd6aec2",
						"last_accessed_at": "2022-06-29T13:33:52.280000000Z",
						"created_at": "2022-06-29T13:33:52.280000000Z",
						"size_in_bytes": 29747
					}
				]
			}`)

	cmd := NewCmdDelete()
	cmd.SetArgs([]string{"--repo", "testOrg/testRepo", "2022-06-29T13:33:49"})
	err := cmd.Execute()

	assert.NotNil(t, err)
	assert.True(t, gock.IsDone(), internal.PrintPendingMocks(gock.Pending()))
}

func TestDeleteWithUnauthorizedRequestForDeleteCaches(t *testing.T) {
	t.Cleanup(gock.Off)

	gock.New("https://api.github.com").
		Delete("/repos/testOrg/testRepo/actions/caches").
		Reply(401).
		JSON(`{
			"message": "Must have admin rights to Repository.",
			"documentation_url": "https://docs.github.com/rest/actions/cache#delete-a-github-actions-cache-for-a-repository-using-a-cache-id"
		}`)

	cmd := NewCmdDelete()
	cmd.SetArgs([]string{"--repo", "testOrg/testRepo", "cacheKey", "--confirm"})
	err := cmd.Execute()

	assert.NotNil(t, err)
	assert.Equal(t, reflect.TypeOf(err), reflect.TypeOf(types.HandledError{}))

	var customError types.HandledError
	errors.As(err, &customError)
	assert.Equal(t, customError.Message, "Must have admin rights to Repository.")

	assert.True(t, gock.IsDone(), internal.PrintPendingMocks(gock.Pending()))
}

func TestDeleteWithInternalServerErrorForDeleteCaches(t *testing.T) {
	t.Cleanup(gock.Off)

	gock.New("https://api.github.com").
		Delete("/repos/testOrg/testRepo/actions/caches").
		Reply(500).
		JSON(`{
			"message": "Internal Server Error",
			"documentation_url": "https://docs.github.com/rest/reference/actions#get-github-actions-cache-Delete-for-a-repository"
		}`)

	cmd := NewCmdDelete()
	cmd.SetArgs([]string{"--repo", "testOrg/testRepo", "cacheKey", "--confirm"})
	err := cmd.Execute()

	assert.NotNil(t, err)
	assert.Equal(t, reflect.TypeOf(err), reflect.TypeOf(types.HandledError{}))

	var customError types.HandledError
	errors.As(err, &customError)
	assert.Equal(t, customError.Message, "We could not process your request due to internal error.")

	assert.True(t, gock.IsDone(), internal.PrintPendingMocks(gock.Pending()))
}
