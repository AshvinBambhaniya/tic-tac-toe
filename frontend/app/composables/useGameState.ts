// State
export const useActiveBox = () => useState<number>('activeBox', () => 9);
export const useCurrentUser = () => useState<string>('currentUser', () => 'X');
export const useMainArr = () => useState<string[]>('mainArr', () => Array(9).fill(''));
export const useIsGameOver = () => useState<boolean>('isGameOver', () => false);
export const useGameWinner = () => useState<string | null>('gameWinner', () => null);
export const useGameDraw = () => useState<boolean>('gameDraw', () => false);
export const useIsOpponentDisconnected = () => useState<boolean>('isOpponentDisconnected', () => false);

export const BOT_USER_ID = "00000000-0000-0000-0000-000000000001";
export const useIsAIGame = () => useState<boolean>('isAIGame', () => false);

// New states for multiplayer
export const useGameID = () => useState<string | null>('gameID', () => null);
export const usePlayerSymbol = () => useState<string | null>('playerSymbol', () => null); // 'X' or 'O'
export const useAllMoves = () => useState<any[]>('allMoves', () => []);
export const useSubGridResults = () => useState<any[]>('subGridResults', () => []);

// Timer States
export const useGameMode = () => useState<string>('gameMode', () => 'normal');
export const useTimeBankX = () => useState<number | null>('timeBankX', () => null);
export const useTimeBankO = () => useState<number | null>('timeBankO', () => null);
export const useMissedTurnsX = () => useState<number>('missedTurnsX', () => 0);
export const useMissedTurnsO = () => useState<number>('missedTurnsO', () => 0);
export const useLastMoveAt = () => useState<string | null>('lastMoveAt', () => null);

export const useGameState = () => {
    const activeBox = useActiveBox();
    const currentUser = useCurrentUser();
    const mainArr = useMainArr();
    const isGameOver = useIsGameOver();
    const gameWinner = useGameWinner();
    const gameDraw = useGameDraw();
    const isOpponentDisconnected = useIsOpponentDisconnected();
    const isAIGame = useIsAIGame();
    const gameID = useGameID();
    const playerSymbol = usePlayerSymbol();
    const allMoves = useAllMoves();
    const subGridResults = useSubGridResults();
    
    const gameMode = useGameMode();
    const timeBankX = useTimeBankX();
    const timeBankO = useTimeBankO();
    const missedTurnsX = useMissedTurnsX();
    const missedTurnsO = useMissedTurnsO();
    const lastMoveAt = useLastMoveAt();

    // In server-authoritative mode, we just update state based on server messages
    const updateFromServer = (payload: any) => {
        const { game, moves, results } = payload;
        const { authUser } = useAuth();
        
        activeBox.value = game.active_sub_grid;
        currentUser.value = game.current_turn;
        isGameOver.value = game.status !== 'ongoing';
        
        const winnerId = game.winner_id?.toLowerCase();
        const playerXId = game.player_x_id?.toLowerCase();
        const playerOId = game.player_o_id?.toLowerCase();
        const currentUserId = authUser.value?.id?.toLowerCase();

        gameWinner.value = winnerId ? (winnerId === playerXId ? 'X' : 'O') : null;
        gameDraw.value = game.status === 'draw';
        isOpponentDisconnected.value = false;
        isAIGame.value = playerOId === BOT_USER_ID;
        gameID.value = game.id;

        // Timer mapping
        gameMode.value = game.game_mode;
        timeBankX.value = game.time_bank_x;
        timeBankO.value = game.time_bank_o;
        missedTurnsX.value = game.missed_turns_x;
        missedTurnsO.value = game.missed_turns_o;
        lastMoveAt.value = game.last_move_at;

        // Set player symbol locally
        if (currentUserId) {
            if (playerXId === currentUserId) {
                playerSymbol.value = 'X';
            } else if (playerOId === currentUserId) {
                playerSymbol.value = 'O';
            }
        }
        
        if (moves) allMoves.value = moves;
        if (results) subGridResults.value = results;
        
        // Update mainArr based on results
        const newMainArr = Array(9).fill('');
        results?.forEach((res: any) => {
            newMainArr[res.grid_index] = res.winner_symbol;
        });
        mainArr.value = newMainArr;
    };

    const resetFullGame = () => {
        activeBox.value = 9;
        currentUser.value = 'X';
        mainArr.value = Array(9).fill('');
        isGameOver.value = false;
        gameWinner.value = null;
        gameDraw.value = false;
        isOpponentDisconnected.value = false;
        isAIGame.value = false;
        gameID.value = null;
        playerSymbol.value = null;
        allMoves.value = [];
        subGridResults.value = [];
        
        gameMode.value = 'normal';
        timeBankX.value = null;
        timeBankO.value = null;
        missedTurnsX.value = 0;
        missedTurnsO.value = 0;
        lastMoveAt.value = null;
    };

    return {
        activeBox,
        currentUser,
        mainArr,
        isGameOver,
        gameWinner,
        gameDraw,
        isOpponentDisconnected,
        isAIGame,
        gameID,
        playerSymbol,
        allMoves,
        subGridResults,
        gameMode,
        timeBankX,
        timeBankO,
        missedTurnsX,
        missedTurnsO,
        lastMoveAt,
        updateFromServer,
        resetFullGame
    };
};
