-- +migrate Up
CREATE TABLE IF NOT EXISTS games (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_x_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    player_o_id UUID REFERENCES users(id) ON DELETE SET NULL,
    current_turn CHAR(1) NOT NULL DEFAULT 'X' CHECK (current_turn IN ('X', 'O')),
    active_sub_grid SMALLINT NOT NULL DEFAULT 9 CHECK (active_sub_grid BETWEEN 0 AND 9),
    winner_id UUID REFERENCES users(id) ON DELETE SET NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'ongoing' CHECK (status IN ('ongoing', 'finished', 'draw')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_games_player_x ON games(player_x_id);
CREATE INDEX idx_games_player_o ON games(player_o_id);
CREATE INDEX idx_games_status ON games(status);
CREATE INDEX idx_games_winner ON games(winner_id);

-- +migrate Down
DROP TABLE IF EXISTS games;
