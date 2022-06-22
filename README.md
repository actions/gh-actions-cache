# gh-actions-cache

## Local Development

1. Build the extension using 

        go build

2. Install the extension


        gh extension install <filepath-to-build>


    If you are already in the same directory use this `gh extension install .`

3. Run the command

        gh actions-cache <command> [Flags]


## Troubleshooting

1. `symlink /Users/.../gh-actions-cache /Users/.../share/gh/extensions/gh-actions-cache: file exists`

Uninstall the current version of extension using

    gh extension remove gh-actions-cache