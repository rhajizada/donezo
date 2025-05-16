package service

import (
	"github.com/rhajizada/donezo/internal/repository"
)

type Board struct {
	repository.Board
}

type Item struct {
	repository.Item
	Tags []string
}
