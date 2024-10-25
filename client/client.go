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

	// Add headers, including an API key if needed.
	req.Header.Set("Content-Type", "application/json")
	if c.APIKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))
	} else {
		return nil, errors.New("missing bearer token")
	}
	return req, nil
}

// ListBoards lists all the boards
func (c *Client) ListBoards() ([]repository.Board, error) {
	// Build the full URL with query parameters.
	reqURL, err := url.Parse(c.BaseURL + "/api/boards")
	if err != nil {
		return nil, err
	}
	req, err := c.NewRequest("GET", reqURL.String(), nil)
	if err != nil {
		return nil, err
	}

	// Perform the request.
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Handle non-200 HTTP responses.
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed with status code: %d", resp.StatusCode)
	}

	// Read and return the response body.
	var boards []repository.Board
	err = json.NewDecoder(resp.Body).Decode(&boards)
	if err != nil {
		return nil, fmt.Errorf("error decoding JSON: %v", err)
	}
	return boards, nil
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
func (c *Client) UpdateBoard(board repository.Board) (*repository.Board, error) {
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

	var boards repository.Board
	err = json.NewDecoder(resp.Body).Decode(&boards)
	if err != nil {
		return nil, fmt.Errorf("error decoding JSON: %v", err)
	}

	return &board, nil
}

// DeleteBoard deletes specidied board
func (c *Client) DeleteBoard(board repository.Board) error {
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
