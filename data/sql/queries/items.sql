-- name: CreateItem :one
INSERT INTO items (
    board_id, title, description
) VALUES (
    ?, ?, ?
)
RETURNING id, board_id, title, description, completed, created_at, last_updated_at;

-- name: UpdateItemByID :one
UPDATE items
SET
    title = ?,
    description = ?,
    completed = ?,
    last_updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: DeleteItemByID :exec
DELETE FROM items
WHERE id = ?;

-- name: GetItemByID :one
SELECT
*
FROM items
WHERE items.id = ?;

-- name: ListItemsByBoardID :many
SELECT
*
FROM items
WHERE items.board_id = ?
ORDER BY items.created_at;
