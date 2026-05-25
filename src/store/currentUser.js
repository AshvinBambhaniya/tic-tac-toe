import { ref } from 'vue'

export const currentUser = ref("X")

export const setCurrentUser = () => {
    (currentUser.value === "X") ? currentUser.value = "O" : currentUser.value = "X";
}

export const resetCurrentUser = () => {
    currentUser.value = "X";
}