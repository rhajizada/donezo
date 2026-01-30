-- name: CreateBoard :one
INSERT INTO boards (
  name
) VALUES (
  ?
)
RETURNING *;

-- name: ListBoards :many
SELECT * FROM boards
ORDER BY id;

-- name: GetBoardByID :one
SELECT * FROM boards
WHERE id = ? LIMIT 1;

-- name: UpdateBoardByID :one
UPDATE boards
SET name = ?,
last_updated_at = CURRENT_TIMESTAMP
WHERE boards.id = ?
RETURNING *;

-- name: DeleteBoardByID :exec
DELETE FROM boards
WHERE id = ?;
