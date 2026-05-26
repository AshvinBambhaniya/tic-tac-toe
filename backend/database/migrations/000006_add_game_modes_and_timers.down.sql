-- +migrate Down
ALTER TABLE games DROP COLUMN game_mode;
ALTER TABLE games DROP COLUMN time_bank_x;
ALTER TABLE games DROP COLUMN time_bank_o;
ALTER TABLE games DROP COLUMN missed_turns_x;
ALTER TABLE games DROP COLUMN missed_turns_o;
ALTER TABLE games DROP COLUMN last_move_at;
