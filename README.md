# gh-actions-cache

✨ A GitHub (`gh`) CLI extension to list and delete actions cache using filters.

## Installation

1. Install the `gh` CLI - see the [installation](https://github.com/cli/cli#installation)
   
   _Installation requires a minimum version (2.0.0) of the the GitHub CLI that supports extensions._

2. Install this extension:

        gh extension install actions/gh-actions-cache

## Usage

        gh actions-cache <command> [flags]

CORE COMMANDS:

	list:		list caches with result length cap of 100
	delete:		delete caches with a key


## Local Development

1. Build the extension using 

        go build

2. Install the extension


        gh extension install <filepath-to-build>


    If you are already in the same directory use this `gh extension install .`

3. Run the command

        gh actions-cache <command> [Flags]

4. Update the golang code and generate new binary using step 1. There is no need to install the binary again ie step 2.
