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
RETURNING id, board_id, title, description, completed, created_at, last_updated_at;

-- name: DeleteItemByID :exec
DELETE FROM items
WHERE id = ?;

-- name: GetItemByID :one
SELECT
    i.id,
    i.board_id,
    i.title,
    i.description,
    i.completed,
    i.created_at,
    i.last_updated_at,
    COALESCE(json_group_array(t.tag), '[]') AS tags
FROM items i
LEFT JOIN tags t ON i.id = t.item_id
WHERE i.id = ?
GROUP BY i.id;

-- name: ListItemsByBoardID :many
SELECT
    i.id,
    i.board_id,
    i.title,
    i.description,
    i.completed,
    i.created_at,
    i.last_updated_at,
    COALESCE(json_group_array(t.tag), '[]') AS tags
FROM items i
LEFT JOIN tags t ON i.id = t.item_id
WHERE i.board_id = ?
GROUP BY i.id
ORDER BY i.created_at;
