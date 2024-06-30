package runner

import (
	"context"
	"log"

	"github.com/alingse/go-linter-runner/runner/submit"
)

const (
	DefaultSource string = `top.txt`
	DefaultCount  int64  = 2000
)

func RunSubmit(sourceFile string, repoCount int64, workflowName string) {
	repos, err := submit.ReadSubmitRepos(sourceFile, repoCount)
	if err != nil {
		log.Fatalf("read submit source file failed %s %+v", sourceFile, err)
		return
	}
	if len(repos) == 0 {
		log.Fatalf("read submit source file got empty %s", sourceFile)
		return
	}
	// Submit
	ctx := context.Background()
	err = submit.SumitActions(ctx, workflowName, repos)
	if err != nil {
		log.Fatalf("submit repos failed with %+v", err)
		return
	}
}
