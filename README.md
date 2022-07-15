# gh-actions-cache

✨ A GitHub (`gh`) [CLI](https://cli.github.com/) extension to manage the GitHub Actions [caches](https://docs.github.com/en/actions/using-workflows/caching-dependencies-to-speed-up-workflows) being used in a GitHub repository. 

It enables listing of active caches in a repo along with capability to filter by cache key or branch. This brings transparency, for example by showing how much storage quota a cache is consuming or which branch a cache was created for etc or how recently was the cache used.

It also allows deleting a corrupt, incomplete or dangling cache. A cache can be deleted by cache key. The key can be easily found either using the list capability or by looking at the cache action log in workflow run logs.

This extension builds on top of [cache management](https://docs.github.com/en/actions/using-workflows/caching-dependencies-to-speed-up-workflows#managing-caches) capabilities exposed by GitHub [APIs](https://docs.github.com/en/rest/actions/cache).

## Installation

1. Install the `gh` CLI - see the [installation](https://github.com/cli/cli#installation)
   
   _Installation requires a minimum version (2.0.0) of the the GitHub CLI that supports extensions._

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


## Local Development

1. Build the extension using 

        go build

2. Install the extension


        gh extension install <filepath-to-build>


    If you are already in the same directory use this `gh extension install .`

3. Run the command

        gh actions-cache <command> [Flags]

4. Update the golang code and generate new binary using step 1. There is no need to install the binary again ie step 2.
