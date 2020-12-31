package deck

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Client for accessing the Trello API.
type Client struct {
	api      string
	username string
	password string
}

// NewClient instance.
func NewClient(api, username, password string) Client {
	return Client{api, username, password}
}

// Board data.
type Board struct {
	ID    int
	Title string
}

// GetBoards of user.
func (cli *Client) GetBoards() ([]Board, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/boards", cli.api), nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(cli.username, cli.password)
	req.Header.Add("OCS-APIRequest", "true")
	req.Header.Add("Content-Type", "application/json")

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		var errMsg struct{ message string }
		err = json.Unmarshal(body, &errMsg)
		if err != nil {
			return nil, fmt.Errorf("trello API: Got %d - %s", resp.StatusCode, resp.Status)
		}
		return nil, fmt.Errorf("trello API: Got %d - %s", resp.StatusCode, errMsg.message)
	}

	var data []Board
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Stack data.
type Stack struct {
	ID    int
	Title string
	Cards []Card
}

// GetStacks from board.
func (cli *Client) GetStacks(board int) ([]Stack, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/boards/%d/stacks", cli.api, board), nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(cli.username, cli.password)
	req.Header.Add("OCS-APIRequest", "true")
	req.Header.Add("Content-Type", "application/json")

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		var errMsg struct{ message string }
		err = json.Unmarshal(body, &errMsg)
		if err != nil {
			return nil, fmt.Errorf("trello API: Got %d - %s", resp.StatusCode, resp.Status)
		}
		return nil, fmt.Errorf("trello API: Got %d - %s", resp.StatusCode, errMsg.message)
	}

	var data []Stack
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// GetStack from board.
func (cli *Client) GetStack(board, stack int) (Stack, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/boards/%d/stacks/%d", cli.api, board, stack), nil)
	if err != nil {
		return Stack{}, err
	}
	req.SetBasicAuth(cli.username, cli.password)
	req.Header.Add("OCS-APIRequest", "true")
	req.Header.Add("Content-Type", "application/json")

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return Stack{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Stack{}, err
	}

	if resp.StatusCode != 200 {
		var errMsg struct{ message string }
		err = json.Unmarshal(body, &errMsg)
		if err != nil {
			return Stack{}, fmt.Errorf("trello API: Got %d - %s", resp.StatusCode, resp.Status)
		}
		return Stack{}, fmt.Errorf("trello API: Got %d - %s", resp.StatusCode, errMsg.message)
	}

	var data Stack
	err = json.Unmarshal(body, &data)
	if err != nil {
		return Stack{}, err
	}

	return data, nil
}

// Card data.
type Card struct {
	ID          int
	Title       string
	Description string
}

// NewCard for stack.
func (cli *Client) NewCard(board, stack int, title, description string) (Card, error) {
	reqBody := map[string]interface{}{}

	reqBody["title"] = title
	reqBody["description"] = description
	reqBody["order"] = 999

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return Card{}, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/boards/%d/stacks/%d/cards", cli.api, board, stack), bytes.NewReader(jsonBody))
	if err != nil {
		return Card{}, err
	}
	req.SetBasicAuth(cli.username, cli.password)
	req.Header.Add("OCS-APIRequest", "true")
	req.Header.Add("Content-Type", "application/json")

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return Card{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Card{}, err
	}

	if resp.StatusCode != 200 {
		var errMsg struct{ message string }
		err = json.Unmarshal(body, &errMsg)
		if err != nil {
			return Card{}, fmt.Errorf("trello API: Got %d - %s", resp.StatusCode, resp.Status)
		}
		return Card{}, fmt.Errorf("trello API: Got %d - %s", resp.StatusCode, errMsg.message)
	}

	var data Card
	err = json.Unmarshal(body, &data)
	if err != nil {
		return Card{}, err
	}

	return data, nil
}

// DeleteCard from stack.
func (cli *Client) DeleteCard(board, stack, card int) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/boards/%d/stacks/%d/cards/%d", cli.api, board, stack, card), nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(cli.username, cli.password)
	req.Header.Add("OCS-APIRequest", "true")
	req.Header.Add("Content-Type", "application/json")

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		var errMsg struct{ message string }
		err = json.Unmarshal(body, &errMsg)
		if err != nil {
			return fmt.Errorf("trello API: Got %d - %s", resp.StatusCode, resp.Status)
		}
		return fmt.Errorf("trello API: Got %d - %s", resp.StatusCode, errMsg.message)
	}

	return nil
}
