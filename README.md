# go-linter-runner

GitHub Action for running Go linters on multi repositories and posting issue comments. Supports single repository runs and batch task submissions.

## Commend Example

### Go-linter-runner Report

**Linter**:     `nilnesserr`
**Repository**:  [https://github.com/pingcap/tidb-dashboard](https://github.com/pingcap/tidb-dashboard)


**⭐ Stars**:    187
**🍴 Forks**:    139
**⌨ Pushed**:    2025-03-31T03:39:07Z

**🧐 Found Issues**:  1

View Action Log: https://github.com/alingse/nilnesserr/actions/runs/14202053344
Report issue:    https://github.com/pingcap/tidb-dashboard/issues

<details>
<summary>Show details (1 issues)</summary>

- <a href="https://github.com/pingcap/tidb-dashboard/blob/master/pkg/apiserver/configuration/service.go#L239">pkg/apiserver/configuration/service.go#L239</a> call function with a nil value error after check error
</details>

## Quick Start

### Single Repository Run
```yaml
- uses: alingse/go-linter-runner@v1.0.1
  with:
    action: run
    repo_url: https://github.com/owner/repo
    install_command: go install github.com/example/linter@latest
    linter_command: linter --flags
  env:
    GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

### Batch Task Submission
```yaml
- uses: alingse/go-linter-runner@v1.0.1
  with:
    action: submit
    submit_source_file: source/top.txt
    submit_repo_count: 1000
  env:
    GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

## Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `go_version` | Go version | 1.22 |
| `action` | Action type: 'run' or 'submit' | run |
| `yaml_config` | run action: YAML config file path for linter | empty |
| `repo_url` | Repository URL to check | required |
| `workdir` | run action: Working directory | . |
| `install_command` | run action: Command to install linter | empty |
| `linter_command` | run action: Command to run linter | empty |
| `excludes` | run action: Strings to exclude in output (regex supported) | [] |
| `includes` | run action: Strings to include in output (regex supported) | [".go"] |
| `issue_id` | run action: Issue ID to comment when problems found | empty |
| `enable_testfile` | run action: Whether to check _test.go files | false |
| `submit_source_file` | submit action: Repository list file | top.txt |
| `submit_repo_count` | submit action: Number of repositories to submit | 1000 |
| `submit_workflow` | submit action: Workflow file to run | go-linter-runner.yml |
| `submit_workflow_ref` | submit action: Workflow branch reference | empty |
| `submit_rate` | submit action: Rate limit (requests per second) | 0.25 |

## Detailed Usage

### Single Repository Run Examples

```yaml
- name: Run with YAML config
  uses: alingse/go-linter-runner@v1.0.1
  with:
    action: run
    yaml_config: .github/jobs/linter-config.yaml
    repo_url: ${{ inputs.repo_url }}
  env:
    GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}

- name: Run with direct parameters
  uses: alingse/go-linter-runner@v1.0.1
  with:
    action: run
    install_command: go install github.com/example/linter@version
    linter_command: linter --flags
    includes: '["go", "github"]'
    excludes: '["ignore this"]'
    issue_id: 1
    repo_url: ${{ inputs.repo_url }}
  env:
    GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

### Batch Task Submission Example

```yaml
- name: Submit batch tasks
  uses: alingse/go-linter-runner@v1.0.1
  with:
    action: submit
    submit_source_file: ${{ inputs.source }}
    submit_repo_count: ${{ inputs.count }}
    submit_workflow: ${{ inputs.workflow }}
  env:
    GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

## Local Execution

```bash
go install github.com/alingse/go-linter-runner@latest

# Submit tasks
go-linter-runner submit -s source/top.txt -c 10000 -w go-linter-runner.yml

# Or using gh client
tail -1000 source/awesome.txt | xargs -I {} gh workflow run go-linter-runner.yml -F repo_url={}
```

## Example Results

See latest comments in https://github.com/alingse/go-linter-runner/issues/1.

## Contribution

Welcome to try, submit Issues and Pull Requests!
