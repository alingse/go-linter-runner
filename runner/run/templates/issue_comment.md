{{ $repoName := repoName .RepositoryURL }}
{{ $repoInfo := .RepoInfo }}

Run `{{.Linter}}` on Repo: {{.RepositoryURL}}

{{if $repoInfo}}
<div style="border: 1px solid #e1e4e8; border-radius: 6px; padding: 16px; margin: 16px 0;">
  <h3 style="margin-top: 0;">{{ $repoName }}</h3>

  {{if $repoInfo.IsArchived}}
  <div style="background-color: #f6f8fa; border-left: 4px solid #d4d72d; padding: 8px 12px; margin-bottom: 12px;">
    ⚠️ This repository is archived
  </div>
  {{end}}

  <div style="display: flex; gap: 16px; margin-bottom: 12px;">
    <div>
      <span style="font-weight: 600;">⭐ Stars:</span> {{ formatCount $repoInfo.StargazerCount }}
    </div>
    <div>
      <span style="font-weight: 600;">🍴 Forks:</span> {{ formatCount $repoInfo.ForkCount }}
    </div>
  </div>

  <div style="margin-bottom: 12px;">
    {{if $repoInfo.IsFork}}
    <div>🔀 This is a fork repository</div>
    {{end}}
    {{if $repoInfo.IsEmpty}}
    <div>📭 This is an empty repository</div>
    {{end}}
  </div>

  <div style="display: flex; gap: 16px;">
    <div>
      <span style="font-weight: 600;">最后推送:</span>
      {{if isOldDate $repoInfo.PushedAt}}
      <span style="background-color: #f6f8fa; padding: 2px 4px; border-radius: 4px;">
        ⚠️ {{ formatDate $repoInfo.PushedAt }}
      </span>
      {{else}}
      {{ formatDate $repoInfo.PushedAt }}
      {{end}}
    </div>
    <div>
      <span style="font-weight: 600;">最后更新:</span>
      {{if isOldDate $repoInfo.UpdatedAt}}
      <span style="background-color: #f6f8fa; padding: 2px 4px; border-radius: 4px;">
        ⚠️ {{ formatDate $repoInfo.UpdatedAt }}
      </span>
      {{else}}
      {{ formatDate $repoInfo.UpdatedAt }}
      {{end}}
    </div>
  </div>
</div>
{{else}}
<div style="border: 1px solid #e1e4e8; border-radius: 6px; padding: 16px; margin: 16px 0;">
  <h3 style="margin-top: 0;">{{ $repoName }}</h3>
  <div>⚠️ Failed to get repository details</div>
</div>
{{end}}

Got total {{len .Lines}} lines output in action: {{ .GithubActionLink }}

<details open>
<summary>Click to expand details</summary>
<ol>{{range $index, $line := .Lines}}
<li>{{$line}}</li>
{{- end}}</ol>
</details>

Report issue: {{ .RepositoryURL }}/issues
