package cmd

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/AlecAivazis/survey/v2"
	"github.com/TwiN/go-color"
	ghRepo "github.com/cli/go-gh/pkg/repository"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringP("repo", "R", "", "Select another repository for finding actions cache.")
	deleteCmd.Flags().StringP("branch", "B", "", "Filter by branch")
	deleteCmd.Flags().Bool("confirm", false, "Delete the cache without asking user for confirmation.")
	deleteCmd.SetHelpTemplate(getDeleteHelp())
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete cache by key",
	Long:  `Delete cache by key`,
	Run: func(cmd *cobra.Command, args []string) {
		COMMAND = "delete"
		r, _ := cmd.Flags().GetString("repo")
		branch, _ := cmd.Flags().GetString("branch")
		confirm, _ := cmd.Flags().GetBool("confirm")
		if len(args) != 1 {
			fmt.Println("accepts 1 arg(s), received " + strconv.Itoa(len(args)))
			return
		}
		key := args[0]

		queryParams := generateQueryParams(branch, 100, key, "", "")

		repo, err := getRepo(r)
		if err != nil {
			fmt.Println(err)
			return
		}

		if !confirm {
			var matchedCaches = getCacheListWithExactMatch(repo, queryParams, key)
			if len(matchedCaches) == 0 {
				fmt.Println("Cache with input key '" + key + "' does not exist")
				return
			}
			fmt.Print("\nYou're going to delete " + strconv.Itoa(len(matchedCaches)) + " cache ")
			if len(matchedCaches) == 1 {
				fmt.Println("entry")
			} else {
				fmt.Println("entries")
			}
			fmt.Println()
			prettyPrintCacheList(matchedCaches)
			choice := ""
			prompt := &survey.Select{
				Message: "Are you sure you want to delete the cache entries?",
				Options: []string{"Delete", "Cancel"},
			}
			err := survey.AskOne(prompt, &choice)
			if err != nil {
				fmt.Println("Error occured while taking input from user while trying to delete cache")
				return
			}
			confirm = choice == "Delete"
			fmt.Println()
		}
		if confirm {
			cachesDeleted := deleteCaches(repo, queryParams)
			if cachesDeleted > 0 {
				src := "\u2713"
				r, _ := utf8.DecodeRuneInString(src)
				fmt.Print(color.Colorize(color.Red, string(r)) + " Deleted " + strconv.FormatFloat(cachesDeleted, 'f', 0, 64) + " cache ")
				if cachesDeleted == 1 {
					fmt.Print("entry")
				} else {
					fmt.Print("entries")
				}
				fmt.Print(" with key " + key + "\n")
			} else {
				fmt.Println("Cache with input key '" + key + "' does not exist")
			}
		}

	},
}

func getDeleteHelp() string {
	return `
gh-actions-cache: Works with GitHub Actions Cache. 

USAGE:
	gh actions-cache delete <key> [flags]

ARGUMENTS:
	key		cache key which needs to be deleted
	
FLAGS:
	-R, --repo <[HOST/]owner/repo>		Select another repository using the [HOST/]OWNER/REPO format
	-B, --branch <string>			Filter by branch
	--confirm				Confirm deletion without prompting

INHERITED FLAGS
	--help		Show help for command

EXAMPLES:
	$ gh actions-cache delete Linux-node-f5dbf39c9d11eba80242ac13
`
}

func getCacheListWithExactMatch(repo ghRepo.Repository, queryParams url.Values, key string) []cacheInfo {
	listApiResponse := listCaches(repo, queryParams)
	var exactMatchedKeys []cacheInfo
	for _, cache := range listApiResponse {
		if strings.EqualFold(key, cache.Key) {
			exactMatchedKeys = append(exactMatchedKeys, cache)
		}
	}
	return exactMatchedKeys
}
