package internal

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/k0kubun/pp"
)

var (
	APIURL = "https://api.todoist.com/sync/v9/sync"
)

type Request struct {
	Type   string                 `json:"type"`
	UUID   string                 `json:"uuid"`
	TempID string                 `json:"temp_id"`
	Args   map[string]interface{} `json:"args"`
}

func (r Request) MarshalJSON() ([]byte, error) {
	if r.TempID == "" {
		return json.Marshal(map[string]interface{}{
			"type": r.Type,
			"uuid": r.UUID,
			"args": r.Args,
		})
	} else {
		return json.Marshal(map[string]interface{}{
			"type":    r.Type,
			"uuid":    r.UUID,
			"args":    r.Args,
			"temp_id": r.TempID,
		})
	}
}

type TodoistAPI struct {
	Token     string
	logger    *slog.Logger
	synctoken string
	backlog   []Request
}

type Response struct {
	SyncToken string    `json:"sync_token"`
	Tasks     []Task    `json:"items"`
	Projects  []Project `json:"projects"`
}

// NewTodoist creates a new Todoist client
func NewDispatcher(token string) *TodoistAPI {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	return &TodoistAPI{Token: token, logger: logger, synctoken: "*"}
}

func (t *TodoistAPI) Commit() error {
	form := url.Values{}
	commands, err := json.Marshal(t.backlog)
	if err != nil {
		return err
	}
	pp.Println(t.backlog)
	// pp.Println(commands)
	form.Add("commands", string(commands))
	req, err := http.NewRequest("POST", APIURL, strings.NewReader(form.Encode()))
	if err != nil {
		t.logger.Error(err.Error())
		return err
	}

	req.Header.Add("Authorization", "Bearer "+t.Token)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.logger.Error(err.Error())
		return err
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.logger.Error(err.Error())
		return err
	}
	pp.Println(string(bodyBytes))
	str := resp.Status
	if resp.StatusCode != 200 {
		t.logger.Error("Error: " + str)
	}
	t.logger.Info("Success: " + str)
	t.backlog = nil
	return nil

}

func (t *TodoistAPI) update(Type string, args map[string]interface{}) error {
	t.create(Type, args, "")
	return nil
}

func (t *TodoistAPI) create(Type string, args map[string]interface{}, TempID string) error {
	t.backlog = append(t.backlog, Request{Type: Type, UUID: uuid.New().String(), Args: args, TempID: TempID})
	return nil
}

func (t *TodoistAPI) Sync() (*Response, error) {
	form := url.Values{}
	form.Add("sync_token", t.synctoken)
	form.Add("resource_types", `["all"]`)

	req, err := http.NewRequest("POST", APIURL, strings.NewReader(form.Encode()))
	if err != nil {
		t.logger.Error(err.Error())
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+t.Token)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.logger.Error(err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	str := resp.Status
	if resp.StatusCode != 200 {
		t.logger.Error("Error: " + str)
		return nil, err
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.logger.Error(err.Error())
		return nil, err
	}

	var jsonResponse map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &jsonResponse); err != nil {
		t.logger.Error(err.Error())
		return nil, err
	}
	var response Response
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		t.logger.Error(err.Error())
		return nil, err
	}
	t.synctoken = response.SyncToken
	pp.Print(response)
	// pp.Println(jsonResponse)
	for _, item := range jsonResponse["items"].([]interface{}) {
		task := item.(map[string]interface{})
		if task["due"] != nil {
			pp.Println(task)
		}
		// pp.Println(task)
	}
	return &response, nil
}
