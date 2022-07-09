# Releasing

Our build system automatically compiles and attaches cross-platform binaries to any git tag named `vX.Y.Z`. The changelog is [generated from git commit log](https://docs.github.com/en/repositories/releasing-projects-on-github/automatically-generated-release-notes).

Users who run official builds of `gh` on their machines will get notified about the new version within a 24 hour period.

To test out the build system, publish a prerelease tag with a name such as `vX.Y.Z-pre.0` or `vX.Y.Z-rc.1`. Note that such a release will still be public, but it will be marked as a "prerelease", meaning that it will never show up as a "latest" release nor trigger upgrade notifications for users.

## General guidelines

* Features to be released should be reviewed and approved at least one day prior to the release.
* Feature releases should bump up the minor version number.

## Tagging a new release

1. `git tag v1.2.3 && git push origin v1.2.3`
2. Wait several minutes for builds to run: <https://github.com/actions/gh-actions-cache/actions>
3. Verify release is displayed and has correct assets: <https://github.com/actions/gh-actions-cache/releases>
4. Scan generated release notes and optionally add a human touch by grouping items under topic sections
6. (Optional) Delete any pre-releases related to this release

If the build fails, there is not a clean way to re-run it. The easiest way would be to start over by deleting the partial release on GitHub and re-publishing the tag. Note that this might be disruptive to users or tooling that were already notified about an upgrade. If a functional release and its binaries are already out there, it might be better to try to manually fix up only the specific workflow tasks that failed. Use your best judgement depending on the failure type.

## Release locally for debugging

A local release can be created for testing without creating anything official on the release page.

1. Make sure GoReleaser is installed: `brew install goreleaser`
2. `goreleaser --skip-validate --skip-publish --rm-dist`
3. Find the built products under `dist/`.
