package cmd

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/actions/gh-actions-cache/internal"
	"github.com/actions/gh-actions-cache/service"
	"github.com/actions/gh-actions-cache/types"
	"github.com/cli/go-gh/pkg/api"
	"github.com/spf13/cobra"
)

func NewCmdList() *cobra.Command {
	COMMAND = "list"

	f := types.ListOptions{}

	var listCmd = &cobra.Command {
		Use:   "list",
		Short: "Lists the actions cache",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 0 {
				return fmt.Errorf(fmt.Sprintf("Invalid argument(s). Expected 0 received %d", len(args)))
			}

			repo, err := internal.GetRepo(f.Repo)
			if err != nil {
				return err
			}

			err = f.Validate()
			if err != nil {
				return err
			}

			artifactCache, err := service.NewArtifactCache(repo, COMMAND, VERSION)
			if err != nil {
				fmt.Printf("error connecting to %s\n", repo.Host())
				fmt.Println("check your internet connection or https://githubstatus.com")
				return nil
			}

			if f.Branch == "" && f.Key == "" {
				totalCacheSize, err := artifactCache.GetCacheUsage()
				if err == nil {
					fmt.Printf("Total caches size %s\n\n", internal.FormatCacheSize(totalCacheSize))
				}
			}

			queryParams := url.Values{}
			f.GenerateQueryParams(queryParams)
			listCacheResponse, err := artifactCache.ListCaches(queryParams)
			if err != nil {
				var httpError api.HTTPError
				if errors.As(err, &httpError) && httpError.StatusCode == 404 {
					fmt.Println("The given repo does not exist.")
					return nil
				} else {
					fmt.Println("We could not process your request due to internal error.")
					return nil
				}
			}

			totalCaches := listCacheResponse.TotalCount
			caches := listCacheResponse.ActionsCaches

			fmt.Printf("Showing %d of %d cache entries in %s/%s\n\n", displayedEntriesCount(len(caches), f.Limit), totalCaches, repo.Owner(), repo.Name())
			internal.PrettyPrintCacheList(caches)
      		return nil
		},
	}

	listCmd.Flags().StringVarP(&f.Repo, "repo", "R", "", "Select another repository for finding actions cache.")
	listCmd.Flags().StringVarP(&f.Branch, "branch", "B", "", "Filter by branch")
	listCmd.Flags().IntVarP(&f.Limit, "limit", "", 30, "Maximum number of items to fetch (default is 30, max limit is 100)")
	listCmd.Flags().StringVarP(&f.Key, "key", "", "", "Filter by key")
	listCmd.Flags().StringVarP(&f.Order, "order", "", "", "Order of caches returned (asc/desc)")
	listCmd.Flags().StringVarP(&f.Sort, "sort", "", "", "Sort fetched caches (last-used/size/created-at)")
	listCmd.SetHelpTemplate(getListHelp())

	return listCmd
}

func displayedEntriesCount(totalCaches int, limit int) int {
	if totalCaches < limit {
		return totalCaches
	}
	return limit
}

func getListHelp() string {
	return `
gh-actions-cache: Works with GitHub Actions Cache. 

USAGE:
	gh actions-cache list [flags]

ARGUMENTS:
	No Arguments

FLAGS:
	-R, --repo <[HOST/]owner/repo>		Select another repository using the [HOST/]OWNER/REPO format
	-B, --branch <string>			Filter by branch
	-L, --limit <int>			Maximum number of items to fetch (default is 30, max limit is 100)
	--key <string>				Filter by key
	--order <string>			Order of caches returned (asc/desc)
	--sort <string>				Sort fetched caches (last-used/size/created-at)

INHERITED FLAGS
	--help		Show help for command

EXAMPLES:
	$ gh actions-cache list
	$ gh actions-cache list --limit 100
	$ gh actions-cache list --order desc
`
}
