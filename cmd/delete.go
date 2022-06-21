package cmd

import (
	"fmt"
	"strconv"
	"unicode/utf8"

	"github.com/AlecAivazis/survey/v2"
	"github.com/TwiN/go-color"
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
		r, _ := cmd.Flags().GetString("repo")
		branch, _ := cmd.Flags().GetString("branch")
		confirm, _ := cmd.Flags().GetBool("confirm")
		if len(args) == 0 {
			fmt.Println(getDeleteHelp())
			return
		}
		key := args[0]

		queryParams := parseInputFlags(branch, 30, key, "", "")

		repo, err := getRepo(r)
		if err != nil {
			fmt.Println(err)
		}
		var userConfirmation bool = true
		var matchedCaches = getCacheListWithExactMatch(repo, queryParams)
		if len(matchedCaches) == 0 {
			fmt.Println("Cache with input key '" + key + "' does not exist")
			return
		}
		if !confirm {
			prettyPrintCacheList(matchedCaches)
			choice := ""
			prompt := &survey.Select{
				Message: "Are you sure you want to delete the cache entries?",
				Options: []string{"Delete", "Cancel"},
			}
			survey.AskOne(prompt, &choice)
			if choice == "Delete" {
				userConfirmation = true
			} else {
				userConfirmation = false
			}
			fmt.Println()
		}
		if userConfirmation {
			if branch != "" {
				queryParams.Add("ref", branch)
			}
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
