
-- name: CreateItem :one
INSERT INTO items (
    board_id, title, description
) VALUES (
    ?, ?, ?
)
RETURNING id, board_id, title, description, state_id, created_at, last_updated_at;

-- name: CreateItemWithState :one
INSERT INTO items (
    board_id, title, description, state_id
) VALUES (
    ?, ?, ?, ?
)
RETURNING id, board_id, title, description, state_id, created_at, last_updated_at;

-- name: UpdateItemByID :one
UPDATE items
SET
    title = COALESCE(?, title),
    description = COALESCE(?, description),
    state_id = COALESCE(?, state_id),
    last_updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING id, board_id, title, description, state_id, created_at, last_updated_at;

-- name: DeleteItem :exec
DELETE FROM items
WHERE id = ?;

-- name: GetItemByID :one
SELECT
    items.id,
    items.board_id,
    items.title,
    items.description,
    items.state_id,
    items.created_at,
    items.last_updated_at,
    json_extract(states, '$') AS "state"
FROM items
JOIN states ON items.state_id = states.id
WHERE items.id = ?;

-- name: ListItemsByBoardID :many
SELECT
    items.id,
    items.board_id,
    items.title,
    items.description,
    items.state_id,
    items.created_at,
    items.last_updated_at,
    json_extract(states, '$') AS "state"
FROM items
JOIN states ON items.state_id = states.id
WHERE items.board_id = ?
ORDER BY items.created_at;

