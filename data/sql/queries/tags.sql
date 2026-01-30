-- name: ListTags :many
SELECT DISTINCT tag FROM tags ORDER BY tag;

-- name: ListTagsByItemID :many
SELECT tag FROM tags
WHERE item_id = ?
ORDER BY tag;

-- name: ListItemsByTag :many
SELECT
    i.id,
    i.board_id,
    i.title,
    i.description,
    i.completed,
    i.created_at,
    i.last_updated_at,
    COALESCE(json_group_array(t2.tag), '[]') AS tags
FROM items i
JOIN tags t ON i.id = t.item_id
LEFT JOIN tags t2 ON i.id = t2.item_id
WHERE t.tag = ?
GROUP BY i.id
ORDER BY i.created_at;

-- name: CountItemsByTag :one
SELECT COUNT(*)
FROM items i
JOIN tags t ON i.id = t.item_id
WHERE t.tag = ?;

-- name: AddTagToItemByID :exec
INSERT INTO tags (item_id, tag)
VALUES (?, ?)
ON CONFLICT(item_id, tag) DO NOTHING;

-- name: RemoveTagFromItemByID :exec
DELETE FROM tags
WHERE item_id = ? AND tag = ?;

-- name: DeleteTag :exec
DELETE FROM tags
WHERE tag = ?;
