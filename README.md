# go-linter-runner

Use GitHub Actions to run Go linters on public Go repositories and report the results as issue comments.

## Background
After implementing certain linters or linter idea, you may want to check if online repositories also have similar problem. For example, https://github.com/alingse/asasalint and https://github.com/ashanbrown/makezero#15 have discovered a large number of online bugs.

However, manually observing the Actions results and ignoring certain specific errors can be quite tedious, so this project can be used to automate the process.

# Usage

It is recommended to integrate this into your GitHub Workflow.

## Configure a workflow to run for a single Repository

Refer to the [.github/workflows/go-linter-runner.yml](https://github.com/alingse/go-linter-runner/blob/main/.github/workflows/go-linter-runner.yml) configuration to set up the parameters for installing and running the linter.

```yaml
- name: Example -> go-linter-runner run with yaml job config
  uses: alingse/go-linter-runner@v1.0.1
  with:
    action: run
    yaml_config: .github/jobs/alingse-makezero.yaml
    repo_url: ${{ inputs.repo_url }}
  env:
    GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}

- name: Example -> go-linter-runner use direct job config
  uses: alingse/go-linter-runner@v1.0.1
  with:
    action: run
    install_command: go install github.com/alingse/makezero@f6a823578e89de5cdfdfef50d4a5d9a09ade16dd
    linter_command: makezero
    includes: '["go", "github"]'
    excludes: '["assigned by index", "funcall", "called by funcs"]'
    issue_id: 1
    repo_url: ${{ inputs.repo_url }}
  env:
    GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

## Configure a workflow to submit multiple run tasks

Refer to the [.github/workflows/go-linter-runner-submit.yml](https://github.com/alingse/go-linter-runner/blob/main/.github/workflows/go-linter-runner-submit.yml) configuration to set up the information needed to submit the tasks.

```yaml
- name: Submit go-linter-runner actions for repos
  uses: alingse/go-linter-runner@v1.0.1
  with:
    action: submit
    submit_source_file: ${{ inputs.source }}
    submit_repo_count: ${{ inputs.count }}
    submit_workflow: ${{ inputs.workflow }}
  env:
    GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

## Run locally with a binary

```bash
go install github.com/alingse/go-linter-runner@latest
go-linter-runner --help
```

# Other

## GitHub Effects

Refer to the latest comments in https://github.com/alingse/go-linter-runner/issues/1 for the effect.

# Contribution

Welcome to try, submit Issues and Pull Requests!