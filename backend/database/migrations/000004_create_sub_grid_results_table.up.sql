-- +migrate Up
CREATE TABLE IF NOT EXISTS sub_grid_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    grid_index SMALLINT NOT NULL CHECK (grid_index BETWEEN 0 AND 8),
    winner_symbol CHAR(1) NOT NULL CHECK (winner_symbol IN ('X', 'O', 'D')),
    won_at_move_id UUID REFERENCES moves(id) ON DELETE SET NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(game_id, grid_index)
);

CREATE INDEX idx_sub_grid_results_game ON sub_grid_results(game_id);

-- +migrate Down
DROP TABLE IF EXISTS sub_grid_results;
