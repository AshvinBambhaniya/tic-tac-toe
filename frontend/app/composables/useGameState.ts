// State
export const useActiveBox = () => useState<number>('activeBox', () => 9);
export const useCurrentUser = () => useState<string>('currentUser', () => 'X');
export const useMainArr = () => useState<string[]>('mainArr', () => Array(9).fill(''));
export const useIsGameOver = () => useState<boolean>('isGameOver', () => false);
export const useGameWinner = () => useState<string | null>('gameWinner', () => null);
export const useGameDraw = () => useState<boolean>('gameDraw', () => false);

export const useGameState = () => {
    const activeBox = useActiveBox();
    const currentUser = useCurrentUser();
    const mainArr = useMainArr();
    const isGameOver = useIsGameOver();
    const gameWinner = useGameWinner();
    const gameDraw = useGameDraw();

    const toggleUser = () => {
        currentUser.value = currentUser.value === 'X' ? 'O' : 'X';
    };

    const setGameOver = (winner: string | null) => {
        isGameOver.value = true;
        gameWinner.value = winner;
        gameDraw.value = !winner;
    };

    const addToMainBoard = (index: number) => {
        if (mainArr.value[index] !== "") return;
        
        mainArr.value[index] = currentUser.value;
        const result = isWin(mainArr.value);
        if (result) {
            setGameOver(result);
        }
    };

    const isSubGridWon = (index: number) => {
        return mainArr.value[index] !== "";
    };

    const resetFullGame = () => {
        activeBox.value = 9;
        currentUser.value = 'X';
        mainArr.value = Array(9).fill('');
        isGameOver.value = false;
        gameWinner.value = null;
        gameDraw.value = false;
    };

    return {
        activeBox,
        currentUser,
        mainArr,
        isGameOver,
        gameWinner,
        gameDraw,
        toggleUser,
        setGameOver,
        addToMainBoard,
        isSubGridWon,
        resetFullGame
    };
};
