package trello

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	trelloAPI = "https://api.trello.com/1"
)

// Client for accessing the Trello API.
type Client struct {
	apiKey string
	token  string
}

// NewClient instance.
func NewClient(key, token string) Client {
	return Client{key, token}
}

// Board data.
type Board struct {
	ID          string
	Name        string
	Description string `json:"desc"`
}

// GetBoards of user.
func (cli *Client) GetBoards() ([]Board, error) {
	resp, err := http.Get(fmt.Sprintf("%s/members/me/boards?key=%s&token=%s", trelloAPI, cli.apiKey, cli.token))
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("trello API: Got %d - %s", resp.StatusCode, body)
	}

	var data []Board
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// List data.
type List struct {
	ID   string
	Name string
}

// GetLists from board.
func (cli *Client) GetLists(board string) ([]List, error) {
	resp, err := http.Get(fmt.Sprintf("%s/boards/%s/lists?key=%s&token=%s", trelloAPI, board, cli.apiKey, cli.token))
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("trello API: Got %d - %s", resp.StatusCode, body)
	}

	var data []List
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Card data.
type Card struct {
	ID          string
	Name        string
	Description string   `json:"desc"`
	Labels      []string `json:"idLabels"`
}

// GetCards from list.
func (cli *Client) GetCards(list string) ([]Card, error) {
	resp, err := http.Get(fmt.Sprintf("%s/lists/%s/cards?key=%s&token=%s", trelloAPI, list, cli.apiKey, cli.token))
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("trello API: Got %d - %s", resp.StatusCode, body)
	}

	var data []Card
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// NewCard for list.
func (cli *Client) NewCard(list, name, desc string, labels []string) (Card, error) {
	reqBody := map[string]interface{}{}

	reqBody["idList"] = list
	reqBody["name"] = name
	reqBody["desc"] = desc
	reqBody["idLabels"] = labels

	if len(labels) == 0 || labels == nil {
		reqBody["idLabels"] = []string{}
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return Card{}, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/cards?key=%s&token=%s", trelloAPI, cli.apiKey, cli.token), bytes.NewReader(jsonBody))
	if err != nil {
		return Card{}, err
	}
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
		return Card{}, fmt.Errorf("trello API: Got %d - %s", resp.StatusCode, body)
	}

	var data Card
	err = json.Unmarshal(body, &data)
	if err != nil {
		return Card{}, err
	}

	return data, nil
}

// DeleteCard from list.
func (cli *Client) DeleteCard(card string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/cards/%s?key=%s&token=%s", trelloAPI, card, cli.apiKey, cli.token), nil)
	if err != nil {
		return err
	}

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
		return fmt.Errorf("trello API: Got %d - %s", resp.StatusCode, body)
	}

	return nil
}

// UpdateCard data.
func (cli *Client) UpdateCard(card, name, desc string, labels []string) (Card, error) {
	reqBody := map[string]interface{}{}

	reqBody["name"] = name
	reqBody["desc"] = desc
	reqBody["idLabels"] = labels

	if len(labels) == 0 || labels == nil {
		reqBody["idLabels"] = []string{}
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return Card{}, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/cards/%s?key=%s&token=%s", trelloAPI, card, cli.apiKey, cli.token), bytes.NewReader(jsonBody))
	if err != nil {
		return Card{}, err
	}
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
		return Card{}, fmt.Errorf("trello API: Got %d - %s", resp.StatusCode, body)
	}

	var data Card
	err = json.Unmarshal(body, &data)
	if err != nil {
		return Card{}, err
	}

	return data, nil
}
