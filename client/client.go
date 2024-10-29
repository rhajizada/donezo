package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/rhajizada/donezo/internal/handler"
	"github.com/rhajizada/donezo/internal/repository"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	APIKey     string
}

// New initializes and returns an API client.
func New(baseURL, apiKey string, timeout time.Duration) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: timeout,
		},
		APIKey: apiKey,
	}
}

// NewRequest returns a new request with prepared authorization headers
func (c *Client) NewRequest(method string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if c.APIKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))
	} else {
		return nil, errors.New("missing bearer token")
	}
	return req, nil
}

// Healthy returns service health status
func (c *Client) Healthy() error {
	reqURL, err := url.Parse(c.BaseURL + "/healthz")
	if err != nil {
		return err
	}

	req, err := c.NewRequest("GET", reqURL.String(), nil)
	if err != nil {
		return err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err
	}

	return nil
}

// ListBoards lists all the boards
func (c *Client) ListBoards() (*[]repository.Board, error) {
	reqURL, err := url.Parse(c.BaseURL + "/api/boards")
	if err != nil {
		return nil, err
	}

	req, err := c.NewRequest("GET", reqURL.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed with status code: %d", resp.StatusCode)
	}

	var boards []repository.Board
	err = json.NewDecoder(resp.Body).Decode(&boards)
	if err != nil {
		return nil, fmt.Errorf("error decoding JSON: %v", err)
	}
	return &boards, nil
}

// CreateBoard creates a new board with the given name
func (c *Client) CreateBoard(boardName string) (*repository.Board, error) {
	reqURL, err := url.Parse(c.BaseURL + "/api/boards")
	if err != nil {
		return nil, err
	}

	bodyStruct := handler.BoardRequest{
		Name: boardName,
	}
	bodyData, err := json.Marshal(bodyStruct)
	body := bytes.NewReader(bodyData)
	if err != nil {
		return nil, fmt.Errorf("error marshalling struct to JSON: %v", err)
	}

	req, err := c.NewRequest("POST", reqURL.String(), body)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed with status code: %d", resp.StatusCode)
	}

	var board repository.Board
	err = json.NewDecoder(resp.Body).Decode(&board)
	if err != nil {
		return nil, fmt.Errorf("error decoding JSON: %v", err)
	}

	return &board, nil
}

// UpdateBoard updates specified board
func (c *Client) UpdateBoard(board *repository.Board) (*repository.Board, error) {
	boardID := strconv.Itoa(int(board.ID))
	reqURL, err := url.Parse(c.BaseURL + "/api/boards/" + boardID)
	if err != nil {
		return nil, err
	}

	bodyStruct := handler.BoardRequest{
		Name: board.Name,
	}
	bodyData, err := json.Marshal(bodyStruct)
	body := bytes.NewReader(bodyData)
	if err != nil {
		return nil, fmt.Errorf("error marshalling struct to JSON: %v", err)
	}

	req, err := c.NewRequest("PUT", reqURL.String(), body)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed with status code: %d", resp.StatusCode)
	}

	var updatedBoard repository.Board
	err = json.NewDecoder(resp.Body).Decode(&updatedBoard)
	if err != nil {
		return nil, fmt.Errorf("error decoding JSON: %v", err)
	}

	return &updatedBoard, nil
}

// DeleteBoard deletes specidied board
func (c *Client) DeleteBoard(board *repository.Board) error {
	boardID := strconv.Itoa(int(board.ID))

	reqURL, err := url.Parse(c.BaseURL + "/api/boards/" + boardID)
	if err != nil {
		return err
	}

	req, err := c.NewRequest("DELETE", reqURL.String(), nil)
	if err != nil {
		return err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed with status code: %d", resp.StatusCode)
	}

	return nil
}

// ListItems lists all items in the specified board
func (c *Client) ListItems(board *repository.Board) (*[]repository.Item, error) {
	boardID := strconv.Itoa(int(board.ID))

	reqURL, err := url.Parse(c.BaseURL + "/api/boards/" + boardID + "/items")
	if err != nil {
		return nil, err
	}

	req, err := c.NewRequest("GET", reqURL.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed with status code: %d", resp.StatusCode)
	}

	var items []repository.Item
	err = json.NewDecoder(resp.Body).Decode(&items)
	if err != nil {
		return nil, fmt.Errorf("error decoding JSON: %v", err)
	}
	return &items, nil
}

// AddItem creates a new item in the specified board
func (c *Client) AddItem(board *repository.Board, title string, description string) (*repository.Item, error) {
	boardID := strconv.Itoa(int(board.ID))

	reqURL, err := url.Parse(c.BaseURL + "/api/boards/" + boardID + "/items")
	if err != nil {
		return nil, err
	}

	bodyStruct := handler.CreateItemRequest{
		Title:       title,
		Description: description,
	}
	bodyData, err := json.Marshal(bodyStruct)
	body := bytes.NewReader(bodyData)
	if err != nil {
		return nil, fmt.Errorf("error marshalling struct to JSON: %v", err)
	}

	req, err := c.NewRequest("POST", reqURL.String(), body)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed with status code: %d", resp.StatusCode)
	}

	var item repository.Item
	err = json.NewDecoder(resp.Body).Decode(&item)
	if err != nil {
		return nil, fmt.Errorf("error decoding JSON: %v", err)
	}

	return &item, nil
}

// UpdateItem updates specified item
func (c *Client) UpdateItem(item *repository.Item) (*repository.Item, error) {
	boardID := strconv.Itoa(int(item.BoardID))
	itemID := strconv.Itoa(int(item.ID))
	reqURL, err := url.Parse(c.BaseURL + "/api/boards/" + boardID + "/items/" + itemID)
	if err != nil {
		return nil, err
	}

	bodyStruct := handler.UpdateItemRequest{
		CreateItemRequest: handler.CreateItemRequest{
			Title:       item.Title,
			Description: item.Description,
		},
		Completed: item.Completed,
	}
	bodyData, err := json.Marshal(bodyStruct)
	body := bytes.NewReader(bodyData)
	if err != nil {
		return nil, fmt.Errorf("error marshalling struct to JSON: %v", err)
	}

	req, err := c.NewRequest("PUT", reqURL.String(), body)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed with status code: %d", resp.StatusCode)
	}

	var updatedItem repository.Item
	err = json.NewDecoder(resp.Body).Decode(&updatedItem)
	if err != nil {
		return nil, fmt.Errorf("error decoding JSON: %v", err)
	}

	return &updatedItem, nil
}

// DeleteItem deletes specified item
func (c *Client) DeleteItem(item *repository.Item) error {
	boardID := strconv.Itoa(int(item.BoardID))
	itemID := strconv.Itoa(int(item.ID))

	reqURL, err := url.Parse(c.BaseURL + "/api/boards/" + boardID + "/items/" + itemID)
	if err != nil {
		return err
	}

	req, err := c.NewRequest("DELETE", reqURL.String(), nil)
	if err != nil {
		return err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed with status code: %d", resp.StatusCode)
	}

	return nil
}
