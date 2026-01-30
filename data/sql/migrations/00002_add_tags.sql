-- +goose Up
-- +goose StatementBegin
PRAGMA foreign_keys = ON;

CREATE TABLE tags (
    item_id INTEGER NOT NULL,
    tag TEXT NOT NULL,
    PRIMARY KEY (item_id, tag),
    FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE
);

CREATE TRIGGER update_item_last_updated_on_tag_insert
AFTER INSERT ON tags
BEGIN
    UPDATE items
    SET last_updated_at = CURRENT_TIMESTAMP
    WHERE id = NEW.item_id;
END;

CREATE TRIGGER update_item_last_updated_on_tag_delete
AFTER DELETE ON tags
BEGIN
    UPDATE items
    SET last_updated_at = CURRENT_TIMESTAMP
    WHERE id = OLD.item_id;
END;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tags;
-- +goose StatementEnd
