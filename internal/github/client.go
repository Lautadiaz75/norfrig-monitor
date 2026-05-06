package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const baseURL = "https://api.github.com"

type Client struct {
	token      string
	httpClient *http.Client
}

func NewClient(token string) *Client {
	return &Client{
		token: token,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// WorkflowRun representa un run de GitHub Actions.
type WorkflowRun struct {
	ID         int64      `json:"id"`
	Name       string     `json:"name"`
	Status     string     `json:"status"`
	Conclusion string     `json:"conclusion"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	HeadCommit HeadCommit `json:"head_commit"`
}

type HeadCommit struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

type runsResponse struct {
	WorkflowRuns []WorkflowRun `json:"workflow_runs"`
}

// GetLatestRun devuelve el run más reciente del repositorio dado.
func (c *Client) GetLatestRun(owner, repo string) (WorkflowRun, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/actions/runs?per_page=1", baseURL, owner, repo)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return WorkflowRun{}, fmt.Errorf("creando request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return WorkflowRun{}, fmt.Errorf("llamando a GitHub: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return WorkflowRun{}, fmt.Errorf("token inválido o sin permisos")
	}
	if resp.StatusCode != http.StatusOK {
		return WorkflowRun{}, fmt.Errorf("GitHub respondió con estado %d", resp.StatusCode)
	}

	var result runsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return WorkflowRun{}, fmt.Errorf("parseando respuesta: %w", err)
	}

	if len(result.WorkflowRuns) == 0 {
		return WorkflowRun{}, fmt.Errorf("no hay runs en este repositorio")
	}

	return result.WorkflowRuns[0], nil
}
