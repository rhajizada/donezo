-- +goose Up
-- +goose StatementBegin
PRAGMA foreign_keys = ON;

CREATE TABLE boards (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    board_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (board_id) REFERENCES boards(id) ON DELETE CASCADE
);

CREATE TRIGGER update_board_last_updated
AFTER UPDATE ON boards
BEGIN
    UPDATE boards SET last_updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER update_item_last_updated
AFTER UPDATE ON items
BEGIN
    UPDATE items SET last_updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER update_board_after_item_insert
AFTER INSERT ON items
BEGIN
    UPDATE boards
    SET last_updated_at = CURRENT_TIMESTAMP
    WHERE id = NEW.board_id;
END;

CREATE TRIGGER update_board_after_item_update
AFTER UPDATE ON items
BEGIN
    UPDATE boards
    SET last_updated_at = CURRENT_TIMESTAMP
    WHERE id = NEW.board_id;
END;

CREATE TRIGGER update_board_after_item_delete
AFTER DELETE ON items
BEGIN
    UPDATE boards
    SET last_updated_at = CURRENT_TIMESTAMP
    WHERE id = OLD.board_id;
END;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS items;
DROP TABLE IF EXISTS boards;
-- +goose StatementEnd
