-- +migrate Up
CREATE TABLE IF NOT EXISTS moves (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    player_id UUID REFERENCES users(id) ON DELETE SET NULL,
    sub_grid_index SMALLINT NOT NULL CHECK (sub_grid_index BETWEEN 0 AND 8),
    cell_index SMALLINT NOT NULL CHECK (cell_index BETWEEN 0 AND 8),
    symbol CHAR(1) NOT NULL CHECK (symbol IN ('X', 'O')),
    move_order INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(game_id, sub_grid_index, cell_index),
    UNIQUE(game_id, move_order)
);

CREATE INDEX idx_moves_game_id ON moves(game_id);
CREATE INDEX idx_moves_game_order ON moves(game_id, move_order);

-- +migrate Down
DROP TABLE IF EXISTS moves;
