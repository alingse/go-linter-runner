package runner

import (
	"context"
	"testing"
	"time"
)

const (
	testPantryID string = `f78bc885-593a-41c6-98ed-d445fdb5bb7f`
)

func TestPantryStorage(t *testing.T) {
	ps := NewPantryStorage[map[string]any](testPantryID, "runner-local-testing")
	var repos = []string{
		"https://github.com/alingse/sundrylint",
		"https://github.com/alingse/asasalint",
	}
	var err error
	var ctx = context.Background()
	for _, repo := range repos {
		var payload = map[string]any{
			"output": "hello",
			"repo":   repo,
		}
		err = ps.SetRepoOutput(ctx, repo, payload)
		time.Sleep(1 * time.Second)
		if err != nil {
			t.Errorf("call SetRepoOutput failed repo %s %+v", repo, err)
			t.Fail()
		}
		repos2, err := ps.GetRepos(ctx)
		time.Sleep(1 * time.Second)
		if err != nil {
			t.Errorf("call GetRepos failed repo %s %+v", repo, err)
			t.Fail()
		}
		t.Logf("call GetRepos got repos= %+v", repos2)
		payload2, err := ps.GetRepoOutput(ctx, repo)
		time.Sleep(1 * time.Second)
		if err != nil {
			t.Errorf("call GetRepoOutput failed repo %s %+v", repo, err)
			t.Fail()
		}
		t.Logf("call GetRepoOutput got payload= %+v", payload2)
	}
	for _, repo := range repos {
		err = ps.DeleteRepo(ctx, repo)
		time.Sleep(1 * time.Second)
		if err != nil {
			t.Errorf("call DeleteRepo failed repo %s %+v", repo, err)
			t.Fail()
		}
	}
	repos2, err := ps.GetRepos(ctx)
	time.Sleep(1 * time.Second)
	if err != nil {
		t.Errorf("call GetRepos failed  %+v", err)
		t.Fail()
	}
	if len(repos2) != 0 {
		t.Errorf("call GetRepos not empty %+v", repos2)
		t.Fail()
	}
}
