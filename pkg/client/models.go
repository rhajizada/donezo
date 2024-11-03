package client

import (
	"github.com/rhajizada/donezo/internal/handler"
	"github.com/rhajizada/donezo/internal/repository"
)

type Board struct {
	repository.Board
}

type Item struct {
	repository.Item
}

type Token struct {
	handler.TokenResponse
}
