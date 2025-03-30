package submit

type SubmitConfig struct {
	Workflow    string `json:"workflow"`
	WorkflowRef string `json:"workflowRef"`
	Source      string `json:"source"`
	RepoCount   int    `json:"repo_count"`
}
