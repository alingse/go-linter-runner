package runner

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/carlmjohnson/requests"
)

type pantryStorage struct {
	pantryID     string
	basketPrefix string
}

func newPantryStorage(pantryID, basketPrefix string) *pantryStorage {
	return &pantryStorage{
		pantryID:     pantryID,
		basketPrefix: basketPrefix,
	}
}

func (ps *pantryStorage) encodeRepo(repo string) string {
	return base64.StdEncoding.EncodeToString([]byte(repo))
}

func (ps *pantryStorage) decodeRepo(repo string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(repo)
	if err != nil {
		return "", err
	}
	repo = string(data)
	if !utf8.ValidString(repo) {
		return "", errors.New("invalid repo name")
	}
	return repo, nil
}

func (ps *pantryStorage) buildRepoURL(repo string) string {
	basketName := ps.basketPrefix + ps.encodeRepo(repo)
	return fmt.Sprintf("https://getpantry.cloud/apiv1/pantry/%s/basket/%s", ps.pantryID, basketName)
}

func (ps *pantryStorage) SetRepoOutput(ctx context.Context, repo string, payload map[string]interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}
	url := ps.buildRepoURL(repo)
	err = requests.
		URL(url).
		BodyBytes(data).
		ContentType("application/json").
		Post().
		Fetch(ctx)
	if err != nil {
		return fmt.Errorf("failed to save payload: %w", err)
	}
	return nil
}

type GetReposResponse struct {
	Baskets []struct {
		Name string `json:"name"`
	} `json:"baskets"`
}

func (ps *pantryStorage) GetRepos(ctx context.Context) ([]string, error) {
	url := fmt.Sprintf("https://getpantry.cloud/apiv1/pantry/%s", ps.pantryID)

	var response GetReposResponse
	err := requests.
		URL(url).
		ToJSON(&response).
		Fetch(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list baskets: %w", err)
	}

	var repos []string
	for _, basket := range response.Baskets {
		if strings.HasPrefix(basket.Name, ps.basketPrefix) {
			repo := strings.TrimPrefix(basket.Name, ps.basketPrefix)
			repo, err := ps.decodeRepo(repo)
			if err != nil {
				return nil, err
			}
			repos = append(repos, repo)
		}
	}
	return repos, nil
}

func (ps *pantryStorage) GetRepoOutput(ctx context.Context, repo string) (map[string]interface{}, error) {
	url := ps.buildRepoURL(repo)
	var payload map[string]interface{}
	err := requests.
		URL(url).
		ToJSON(&payload).
		Fetch(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get payload: %w", err)
	}

	return payload, nil
}

func (ps *pantryStorage) DeleteRepo(ctx context.Context, repo string) error {
	url := ps.buildRepoURL(repo)
	err := requests.
		URL(url).
		Delete().
		Fetch(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete repo: %w", err)
	}
	return nil
}
