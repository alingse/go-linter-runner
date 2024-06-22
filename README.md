# go-linter-runner

A launcher to run a specific linter for a large number of Go repositories using Github Actions.

# Why

After implementing certain linters, you may want to check if online repositories also have similar issues.

For example, https://github.com/alingse/asasalint and https://github.com/ashanbrown/makezero/pull/15

Both have detected numerous online bugs.

However, manually observing the results of Actions and ignoring some specific errors can be quite exhausting.

Here, you can configure additional `includes` and `excludes` to exclude false-positives that don't need to be addressed.

# Effects

see
1. https://github.com/alingse/go-linter-runner/issues/1
2. https://github.com/alingse/go-linter-runner/issues/2


# Dev

## storage

we use https://getpantry.cloud/ as center storage when run a job

# TODO

- [] Ability to configure more linter yaml
- [] Provide data aggregation (multiple failed Actions consolidated into a Sheet or a specific issue)
- [] Automatic issue creation for projects(may use AI)
- [] Automatic execution
