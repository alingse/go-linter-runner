{{ $repoName := repoName .RepositoryURL }}
{{ $repoInfo := .RepoInfo }}
Run `{{.Linter}}` on Repo: {{.RepositoryURL}}

### Repo

{{if $repoInfo}}
- ⭐[Stars]({{.RepositoryURL}}/stargazers): {{ formatCount $repoInfo.StargazerCount }}
- 🍴[Forks]({{.RepositoryURL}}/network/members): {{ formatCount $repoInfo.ForkCount }}
- PushedAt: {{ formatDate $repoInfo.PushedAt }}
- Status: {{if $repoInfo.IsArchived}}⚠ Archived{{end}}{{if isOldDate $repoInfo.PushedAt}}, ⚠ Last Commit {{ yearsSince $repoInfo.PushedAt }} years ago{{end}}
{{else}}
- Status: ⚠ Failed to get repository details
{{end}}

### Result

Got total {{len .Lines}} lines output in action: {{ .GithubActionLink }}

<details open>
<summary>Click to expand details</summary>
<ol>{{range $index, $line := .Lines}}
<li>{{$line}}</li>
{{- end}}</ol>
</details>

Report issue: {{ .RepositoryURL }}/issues
