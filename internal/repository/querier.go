// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package repository

import (
	"context"
)

type Querier interface {
	CreateBoard(ctx context.Context, name string) (Board, error)
	CreateItem(ctx context.Context, arg CreateItemParams) (Item, error)
	DeleteBoardByID(ctx context.Context, id int64) error
	DeleteItem(ctx context.Context, id int64) error
	GetBoardByID(ctx context.Context, id int64) (Board, error)
	GetItemByID(ctx context.Context, id int64) (Item, error)
	ListBoards(ctx context.Context) ([]Board, error)
	ListItemsByBoardID(ctx context.Context, boardID int64) ([]Item, error)
	UpdateBoardByID(ctx context.Context, arg UpdateBoardByIDParams) (Board, error)
	UpdateItemByID(ctx context.Context, arg UpdateItemByIDParams) (Item, error)
}

var _ Querier = (*Queries)(nil)
