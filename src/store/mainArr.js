import { isWin } from "../composable/gameStatus.js";
import { currentUser } from "./currentUser.js";
import { ref } from "vue";
import { activeBox } from "./activeBox.js";
import { setGameOver } from "./gameStatusStore.js";

export const mainArr = ref(Array(9).fill(""));

export const resetMainArr = () => {
    mainArr.value = Array(9).fill("");
}

export const addToarrayOfMain = (num) => {
    if (mainArr.value[num] !== "") {
        return;
    }
    mainArr.value[num] = currentUser.value;
    const result = isWin(mainArr.value);

    if (result) {
        setGameOver(result);
    }
}

export const isAlreadyWinBySomeone = (index) => {
    if (mainArr.value[index] !== "") {
        return true;
    }
    return false
}