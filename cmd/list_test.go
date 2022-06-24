package cmd
import (
	"fmt"
	// "net/http"
	// "os"
	"strings"
	"testing"

	// "github.com/cli/go-gh/pkg/api"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	// "github.com/actions/gh-actions-cache/internal"
	// "github.com/actions/gh-actions-cache/service"
)

func TestRESTClient(t *testing.T) {
	// stubConfig(t, testConfig())
	// t.Cleanup(gock.Off)

	gock.New("https://api.github.com").
		Get("/repos/actions/gh-actions-cache/actions/cache/usage").
		Reply(200).
		JSON(`{
			"full_name": "t-dedah/vipul-bugbash",
			"active_caches_size_in_bytes": 29152205,
			"active_caches_count": 12
		}`)
	
	gock.New("https://api.github.com").
		Get("/repos/actions/gh-actions-cache/actions/caches").
		Reply(200).
		JSON(`{
				"total_count": 2,
				"actions_caches": [
					{
						"id": 29,
						"ref": "refs/heads/master",
						"key": "Linux-build-cache-node-modules-3fd22dd3a926d576e2562e8b76a5ff157cd3b986f3d44195acfe7efa6bc05919-8",
						"version": "7fcda33c1e1d849a13bcc06f49b9ab64efc01ca9dabe4d7a8d0d387feef4fc88",
						"last_accessed_at": "2022-06-22T20:32:45.550000000Z",
						"created_at": "2022-06-22T20:32:45.550000000Z",
						"size_in_bytes": 2432967
					},
					{
						"id": 27,
						"ref": "refs/heads/master",
						"key": "Linux-build-cache-node-modules-3fd22dd3a926d576e2562e8b76a5ff157cd3b986f3d44195acfe7efa6bc05919-10",
						"version": "7fcda33c1e1d849a13bcc06f49b9ab64efc01ca9dabe4d7a8d0d387feef4fc88",
						"last_accessed_at": "2022-06-22T20:32:41.143333300Z",
						"created_at": "2022-06-22T20:32:41.143333300Z",
						"size_in_bytes": 2429442
					}
				]
			}`)

	cmd := NewCmdList()
	// cmd.SetArgs([]string{"--limit",  "10"})
	stdout := cmd.Execute()
	fmt.Println(stdout)

	// client, err := RESTClient(nil)
	// assert.NoError(t, err)

	// res := struct{ Message string }{}
	// err = client.Do("GET", "some/test/path", nil, &res)
	assert.NoError(t, err)
	// assert.True(t, gock.IsDone(), printPendingMocks(gock.Pending()))
	// assert.Equal(t, "success", stdout)
	// assert.Nil(err)
	assert.True(t, gock.IsDone(), printPendingMocks(gock.Pending()))
}

func printPendingMocks(mocks []gock.Mock) string {
	paths := []string{}
	for _, mock := range mocks {
		paths = append(paths, mock.Request().URLStruct.String())
	}
	return fmt.Sprintf("%d unmatched mocks: %s", len(paths), strings.Join(paths, ", "))
}
