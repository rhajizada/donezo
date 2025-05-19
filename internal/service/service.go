package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/rhajizada/donezo/internal/repository"
)

type Service struct {
	Repo *repository.Queries
}

func New(r *repository.Queries) *Service {
	return &Service{
		Repo: r,
	}
}

// unmarshalTags converts the interface returned from sqlc for the tags field
// into a slice of strings by performing the proper type assertions.
func unmarshalTags(v interface{}) []string {
	var tags []string
	switch t := v.(type) {
	case []byte:
		if err := json.Unmarshal(t, &tags); err != nil {
			return []string{}
		}
	case string:
		if err := json.Unmarshal([]byte(t), &tags); err != nil {
			return []string{}
		}
	default:
		return []string{}
	}
	return tags
}

// Board-related functions.
func (s *Service) ListBoards(ctx context.Context) (*[]Board, error) {
	data, err := s.Repo.ListBoards(ctx)
	if err != nil {
		return nil, err
	}
	boards := make([]Board, len(data))
	for i, b := range data {
		boards[i] = Board{b}
	}
	return &boards, nil
}

func (s *Service) CreateBoard(ctx context.Context, boardName string) (*Board, error) {
	data, err := s.Repo.CreateBoard(ctx, boardName)
	if err != nil {
		return nil, err
	}
	return &Board{data}, nil
}

func (s *Service) UpdateBoard(ctx context.Context, board *Board) (*Board, error) {
	params := repository.UpdateBoardByIDParams{
		Name: board.Name,
		ID:   board.ID,
	}
	data, err := s.Repo.UpdateBoardByID(ctx, params)
	if err != nil {
		return nil, err
	}
	return &Board{data}, nil
}

func (s *Service) DeleteBoard(ctx context.Context, board *Board) error {
	return s.Repo.DeleteBoardByID(ctx, board.ID)
}

// ListItemsByBoard uses the aggregated JSON query and unmarshals the tags.
func (s *Service) ListItemsByBoard(ctx context.Context, board *Board) (*[]Item, error) {
	data, err := s.Repo.ListItemsByBoardID(ctx, board.ID)
	if err != nil {
		return nil, err
	}
	items := make([]Item, len(data))
	for i, v := range data {
		tags := unmarshalTags(v.Tags)
		items[i] = Item{
			Item: repository.Item{
				ID:            v.ID,
				BoardID:       v.BoardID,
				Title:         v.Title,
				Description:   v.Description,
				Completed:     v.Completed,
				CreatedAt:     v.CreatedAt,
				LastUpdatedAt: v.LastUpdatedAt,
			},
			Tags: tags,
		}
	}
	return &items, nil
}

// ListItemsByTag uses the updated aggregated query and unmarshals the tags.
func (s *Service) ListItemsByTag(ctx context.Context, tag string) (*[]Item, error) {
	data, err := s.Repo.ListItemsByTag(ctx, tag)
	if err != nil {
		return nil, err
	}
	items := make([]Item, len(data))
	for i, v := range data {
		tags := unmarshalTags(v.Tags)
		items[i] = Item{
			Item: repository.Item{
				ID:            v.ID,
				BoardID:       v.BoardID,
				Title:         v.Title,
				Description:   v.Description,
				Completed:     v.Completed,
				CreatedAt:     v.CreatedAt,
				LastUpdatedAt: v.LastUpdatedAt,
			},
			Tags: tags,
		}
	}
	return &items, nil
}

func (s *Service) CreateItem(ctx context.Context, board *Board, title string, description string) (*Item, error) {
	params := repository.CreateItemParams{
		BoardID:     board.ID,
		Title:       title,
		Description: description,
	}
	data, err := s.Repo.CreateItem(ctx, params)
	if err != nil {
		return nil, err
	}
	return &Item{
		Item: data,
		Tags: []string{},
	}, nil
}

func (s *Service) UpdateItem(ctx context.Context, item *Item) (*Item, error) {
	params := repository.UpdateItemByIDParams{
		Title:       item.Title,
		Description: item.Description,
		Completed:   item.Completed,
		ID:          item.ID,
	}

	emptyTags := false
	for _, tag := range item.Tags {
		if len(tag) == 0 {
			emptyTags = true
			break
		}
	}
	if emptyTags {
		return nil, errors.New("tag must not be empty")
	}

	data, err := s.Repo.UpdateItemByID(ctx, params)
	if err != nil {
		return nil, err
	}

	// Synchronize tags.
	existingTags := s.listTagsByItemID(ctx, data.ID)
	existingTagsMap := make(map[string]struct{}, len(existingTags))
	for _, t := range existingTags {
		existingTagsMap[t] = struct{}{}
	}
	newTagsMap := make(map[string]struct{}, len(item.Tags))
	for _, t := range item.Tags {
		newTagsMap[t] = struct{}{}
	}

	// Remove tags that are not in the updated item.
	for _, t := range existingTags {
		if _, found := newTagsMap[t]; !found {
			err := s.Repo.RemoveTagFromItemByID(ctx, repository.RemoveTagFromItemByIDParams{
				ItemID: data.ID,
				Tag:    t,
			})
			if err != nil {
				return nil, err
			}
		}
	}

	// Add new tags that are missing in the database.
	for _, t := range item.Tags {
		if _, found := existingTagsMap[t]; !found {
			err := s.Repo.AddTagToItemByID(ctx, repository.AddTagToItemByIDParams{
				ItemID: data.ID,
				Tag:    t,
			})
			if err != nil {
				return nil, err
			}
		}
	}

	return &Item{
		Item: data,
		Tags: item.Tags, // Return the updated tags.
	}, nil
}

func (s *Service) DeleteItem(ctx context.Context, item *Item) error {
	return s.Repo.DeleteItemByID(ctx, item.ID)
}

// Private helper to get tags for an item.
func (s *Service) listTagsByItemID(ctx context.Context, itemID int64) []string {
	tags, err := s.Repo.ListTagsByItemID(ctx, itemID)
	if err != nil {
		return []string{}
	}
	return tags
}

// Tag-related functions.
func (s *Service) ListTags(ctx context.Context) ([]string, error) {
	return s.Repo.ListTags(ctx)
}

func (s *Service) DeleteTag(ctx context.Context, tag string) error {
	return s.Repo.DeleteTag(ctx, tag)
}

func (s *Service) CountItemsByTag(ctx context.Context, tag string) (int64, error) {
	return s.Repo.CountItemsByTag(ctx, tag)
}

func ItemsToMarkdown(header string, items []Item) string {
	var md []string
	md = append(md, fmt.Sprintf("### %s", header))
	for _, v := range items {
		status := " "
		if v.Completed {
			status = "X"
		}
		md = append(md, fmt.Sprintf("- [%s] **%s**", status, v.Title))
		md = append(md, fmt.Sprintf("\t- %s", v.Description))
	}
	return strings.Join(md, "\n")
}
