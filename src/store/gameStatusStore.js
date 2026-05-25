import { ref } from 'vue';

export const isGameOver = ref(false);
export const gameWinner = ref(null);
export const gameDraw = ref(false);

export const setGameOver = (winner) => {
    isGameOver.value = true;
    gameWinner.value = winner;
    gameDraw.value = !winner;
};

export const resetGameStatus = () => {
    isGameOver.value = false;
    gameWinner.value = null;
    gameDraw.value = false;
};
