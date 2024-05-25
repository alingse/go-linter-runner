# go-linter-runner

利用 Github Actions 来为大量go仓库跑某个特定 linter 的启动器

## Why

在实现某些 linter 之后，想检查一下线上仓库是否也会犯这样的错误,

比如 https://github.com/alingse/asasalint 和 https://github.com/ashanbrown/makezero/pull/15

都曾经检测出来不少在线 bug

但是原来人肉观察 Actions 的结果, 以及忽略一些特别的错误，总是比较耗费精力,

这里可以额外配置一些 includes 和 excludes 以排除不需要关注的 false-positive

## 效果


## TODO

- 期望能配置更多的 linter yaml
- 提供数据汇聚(多个失败的 Actions 集中放到一个 Sheet 或者某个issue)
- 自动给项目创建 Issues 功能
- 自动运行
