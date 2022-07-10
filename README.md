# gh-actions-cache

✨ A GitHub (`gh`) CLI extension to list and delete actions cache using filters.

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

``` 
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
```

### Delete 

```
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
```


## Local Development

1. Build the extension using 

        go build

2. Install the extension


        gh extension install <filepath-to-build>


    If you are already in the same directory use this `gh extension install .`

3. Run the command

        gh actions-cache <command> [Flags]

4. Update the golang code and generate new binary using step 1. There is no need to install the binary again ie step 2.
