# gh-actions-cache

✨ A GitHub (`gh`) [CLI](https://cli.github.com/) extension to manage the GitHub Actions [caches](https://docs.github.com/en/actions/using-workflows/caching-dependencies-to-speed-up-workflows) being used in a GitHub repository. 

It enables listing of active caches in a repo along with capability to filter by cache key or branch. This brings transparency, for example by showing how much storage quota a cache is consuming or which branch a cache was created for etc or how recently was the cache used.

It also allows deleting a corrupt, incomplete or dangling cache. A cache can be deleted by cache key. The key can be easily found either using the list capability or by looking at the cache action log in workflow run logs.

This extension builds on top of [cache management](https://docs.github.com/en/actions/using-workflows/caching-dependencies-to-speed-up-workflows#managing-caches) capabilities exposed by GitHub [APIs](https://docs.github.com/en/rest/actions/cache).

**Note:** This extension only supports github.com and GitHub Enterprise Server 3.7 & above.

## Installation

1. Install the `gh` CLI - see the [installation](https://github.com/cli/cli#installation)
   
   _Installation requires a minimum version (2.0.0) of the GitHub CLI that supports extensions._

2. Install this extension:

        gh extension install actions/gh-actions-cache

## Usage

    gh actions-cache <command> [flags]

#### Commands:

S.No  | Commands | Description
------------- | ------------- | -------------
1  | list | list caches with result length cap of 100
2  | delete | delete caches with a key

### List

List active Actions caches in a repository with ability to filter and sort.

``` 
USAGE:
	gh actions-cache list [flags]


ARGUMENTS:
	No Arguments


FLAGS:
	-R, --repo <[HOST/]owner/repo>		Select another repository using the [HOST/]OWNER/REPO format
	-B, --branch <string>			Filter by branch
	-L, --limit <int>			Maximum number of items to fetch (default is 30, max limit is 100)
	--key <string>				Filter by a key or key prefix
	--order <string>			Order of caches returned (asc/desc)
	--sort <string>				Sort fetched caches (last-used/size/created-at)


INHERITED FLAGS
	--help		Show help for command


EXAMPLES:
	$ gh actions-cache list
	$ gh actions-cache list --key 564-node-a68c45df0f45f888039d32cd3a579992574e837406488e8904431197f20521d6
	$ gh actions-cache list --key 564-node-           // key prefix match
	$ gh actions-cache list -B main
	$ gh actions-cache list -B refs/pull/2/merge      // Use the full ref format for PR branches
	$ gh actions-cache list --limit 100
	$ gh actions-cache list --sort size --order desc  // biggest caches first
```

### Delete 

Deletes actions caches with specific cache key. It asks for confirmation before deletion.

```
USAGE:
	gh actions-cache delete <key> [flags]


ARGUMENTS:
	key		cache key which needs to be deleted

	
FLAGS:
	-R, --repo <[HOST/]owner/repo>		Select another repository using the [HOST/]OWNER/REPO format
	-B, --branch <string>			Delete caches specific to branch. Use the full ref format e.g. refs/heads/main
	--confirm				Confirm deletion without prompting


INHERITED FLAGS
	--help		Show help for command
        

EXAMPLES:
	$ gh actions-cache delete Linux-node-f5dbf39c9d11eba80242ac13
```


> ℹ️ There could be multiple caches in a repo with same key. This can happen when different caches with same key have been created for different branches. it may also happen if the `version` property of the cache is different which usually means that cache with same key was created for different OS or with different [paths](https://github.com/actions/cache#inputs).

## FAQs

### How the current repository is selected?

This extension currently uses the `go-gh` module's [`CurrentRepository` function](https://github.com/actions/gh-actions-cache/blob/d3293b69e1c5bc17686d815ab2c64618618c95df/internal/utils.go#L26) to determine the current repo. This function returns the first element of the list returned by the `git.Remotes()` internal function, which [sorts remotes such that `upstream` precedes `github`, which precedes `origin`](https://github.com/cli/go-gh/blob/c2fc965daac88a8a38dd8af02f236095b5dd48f1/internal/git/remote.go#L30). As such, if an `upstream` remote is present, this extension's default behavior is to return its caches. 

User's input `--repo <owner>/<name>` will override any current git repository and extension will fetch caches for the same.

### How to remove trimming in results

We support a table printer that allows users to pipe output for further processing. If we want to list down certain columns without trimming then just selecting the column number in the below command will work.

`gh actions-cache list -R <owner>/<repo_name> | cut -f 1,2,3`

This will print columns 1,2 and 3 without any trimming.

### Delete all caches for a branch

Please refers to this doc - [Force deleting cache entries](https://docs.github.com/en/actions/using-workflows/caching-dependencies-to-speed-up-workflows#force-deleting-cache-entries)

## Contributing
If anything feels off, or if you feel that some functionality is missing, please check out the [contributing page](CONTRIBUTING.md). There you will find instructions for sharing your feedback, building the tool locally, and submitting pull requests to the project.
