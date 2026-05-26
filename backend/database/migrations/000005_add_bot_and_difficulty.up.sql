-- +migrate Up
ALTER TABLE games ADD COLUMN difficulty SMALLINT DEFAULT 0;

INSERT INTO users (id, first_name, last_name, email, password, roles)
VALUES ('00000000-0000-0000-0000-000000000001', 'DeepBlue', 'Bot', 'bot@ultimate-tictactoe.com', NULL, 'bot')
ON CONFLICT (id) DO NOTHING;
