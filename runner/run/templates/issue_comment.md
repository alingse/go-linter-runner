{{ $repoInfo := .RepoInfo }}
{{ $warning := buildWarning $repoInfo}}
### Go-linter-runner Report

**Linter**:     `{{ .Linter }}`
**Repository**:  [{{ .RepositoryURL }}]({{ .RepositoryURL }})

{{if $repoInfo}}
**⭐ Stars**:    {{if $repoInfo}}{{ formatCount $repoInfo.StargazerCount }}{{end}}
**🍴 Forks**:    {{if $repoInfo}}{{ formatCount $repoInfo.ForkCount }}{{end}}
**⌨ Pushed**:    {{if $repoInfo}}{{ $repoInfo.PushedAt }}{{end}}{{end}}{{if $warning}}
**🚨 Warning**:  {{$warning}}{{end}}

**🧐 Found Issues**:  {{len .Lines}}

View Action Log: {{ .GithubActionLink }}
Report issue:    {{ .RepositoryURL }}/issues

<details>
<summary>Show details ({{len .Lines}} issues)</summary>
{{range $index, $line := .Lines}}
- {{$line}}
{{- end}}
</details>
