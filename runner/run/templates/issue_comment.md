Run `{{.Linter}}` on Repo: {{.RepositoryURL}}

Got total {{len .Lines}} lines output in action: {{ .GithubActionLink }}

<details>
<summary>Expand</summary>
<ol>{{range $index, $line := .Lines}}
<li>{{$line}}</li>
{{- end}}</ol>
</details>

Report issue: {{ .RepositoryURL }}/issues
