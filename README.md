# go-linter-runner

使用 GitHub Actions 为公开的 Go 仓库运行 Go linter 并以 issue comment 的形式报告

# 背景

在实现某些 linter 后,您可能希望检查在线仓库是否也存在类似的问题。例如,https://github.com/alingse/asasalint 和 https://github.com/ashanbrown/makezero#15 都发现了大量的在线 bug
但是,手动观察 Actions 的结果并忽略一些特定的错误可能相当繁琐，因此可以使用此项目自动化

# 使用

推荐集成到 Github Workflow 中

## 配置检查单个 Repo 的 Workflow

参考 [`.github/workflows/go-linter-runner.yml`](https://github.com/alingse/go-linter-runner/blob/main/.github/workflows/go-linter-runner.yml) 配置安装和运行 linter 的参数

```yaml
      - name: Example -> go-linter-runner run with yaml job config
        uses: alingse/go-linter-runner@main
        with:
          action: run
          yaml_config: .github/jobs/alingse-makezero.yaml
          repo_url: ${{ inputs.repo_url }}
        env:
          GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}

      - name: Example -> go-linter-runner use direct job config
        uses: alingse/go-linter-runner@main
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

## 配置提交大量检查任务的 Workflow

参考 [`.github/workflows/go-linter-runner-submit.yml`](https://github.com/alingse/go-linter-runner/blob/main/.github/workflows/go-linter-runner-submit.yml) 配置需要提交的信息

```yaml
      - name: Example -> go-linter-runner submit 10 repos
        uses: alingse/go-linter-runner@main
        with:
          action: submit
          submit_source_file: top.txt
          submit_repo_count: ${{ inputs.count }}
          submit_workflow: ${{ inputs.workflow }}
        env:
          GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

## 以本地二进制文件运行

```bash
go install github.com/alingse/go-linter-runner@latest
```

运行参考

```bash
go-linter-runner --help
```

# 其他

## Github 效果

参考 https://github.com/alingse/go-linter-runner/issues/1 的最新评论效果

# 贡献

欢迎试用, 提交 Issue 和 Pull Request!
