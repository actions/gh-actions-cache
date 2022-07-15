## Contributing

Hi! Thanks for your interest in contributing to the GitHub Actions Cache CLI Extension! 

We accept pull requests for bug fixes and features where we've discussed the approach in an issue and given the go-ahead for a community member to work on it. We'd also love to hear about ideas for new features as issues.

Please do:

* Check existing issues to verify that the [bug][bug issues] or [feature request][feature request issues] has not already been submitted.
* Open an issue if things aren't working as expected.
* Open an issue to propose a significant change.
* Open a pull request to fix a bug.
* Open a pull request to fix documentation about a command.

Please avoid:

* Opening pull requests for issues marked `needs-design`, `needs-investigation`, or `blocked`.
* Adding installation instructions specifically for your OS/package manager.

## Building the project

Prerequisites:
- Go 1.16+
- [GitHub CLI][cli]

Build with:
* Unix-like systems/ MacOS/ Windows: `go build`

Install the binary extension as:
* Unix-like systems/ MacOS/ Windows: `gh extension install .` from the root folder the project.

Run tests with: `go test -v ./...`

## Submitting a pull request

1. Fork/Create a new branch: `git checkout -b my-branch-name`
1. Make your change, add tests, and ensure tests pass
1. Submit a pull request: `gh pr create --web`

Contributions to this project are [released][legal] to the public under the [project's open source license][license].

Please note that this project adheres to a [Contributor Code of Conduct][code-of-conduct]. By participating in this project you agree to abide by its terms.

We generate manual pages from source on every release. You do not need to submit pull requests for documentation specifically; manual pages for commands will automatically get updated after your pull requests gets accepted.

## Resources

- [How to Contribute to Open Source][]
- [Using Pull Requests][]
- [GitHub Help][]


[bug issues]: https://github.com/actions/gh-actions-cache/issues?q=is%3Aopen+is%3Aissue+label%3Abug
[feature request issues]: https://github.com/actions/gh-actions-cache/issues?q=is%3Aopen+is%3Aissue+label%3Aenhancement
[legal]: https://docs.github.com/en/free-pro-team@latest/github/site-policy/github-terms-of-service#6-contributions-under-repository-license
[license]: ./LICENSE
[code-of-conduct]: ./CODE_OF_CONDUCT.md
[How to Contribute to Open Source]: https://opensource.guide/how-to-contribute/
[Using Pull Requests]: https://docs.github.com/en/free-pro-team@latest/github/collaborating-with-issues-and-pull-requests/about-pull-requests
[GitHub Help]: https://docs.github.com/
[cli]: https://cli.github.com/