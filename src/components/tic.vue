<template>
  <div class="sub-grid-container" :class="{ 
    'active-grid': isGridActive,
    'won-grid': win 
  }">
    <div v-for="(item, index) in arr" :key="index" 
         class="cell" 
         :class="{ 'cell-occupied': item !== '' }"
         @click="addturn(index)">
      <transition name="pop">
        <span v-if="item" :class="item === 'X' ? 'x-mark' : 'o-mark'">{{ item }}</span>
      </transition>
    </div>
    
    <!-- Win Overlay -->
    <transition name="fade">
      <div v-if="win" class="win-overlay">
        <span :class="winner === 'X' ? 'x-mark' : 'o-mark'">{{ winner }}</span>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, defineProps, computed } from "vue";
import { currentUser, setCurrentUser } from '../store/currentUser.js'
import { activeBox, setActiveBox } from '../store/activeBox.js'
import { addToarrayOfMain, isAlreadyWinBySomeone } from "../store/mainArr.js";
import { isDraw, isWin } from "../composable/gameStatus.js";
import { isGameOver } from '../store/gameStatusStore.js';

const arr = ref(Array(9).fill(""));
const win = ref(false)
const winner = ref("")

const props = defineProps({
  ticId: Number,
})

const isGridActive = computed(() => {
  if (isGameOver.value) return false;
  return !win.value && (activeBox.value === 9 || activeBox.value === props.ticId);
});

function addturn(i) {
  if (isGridActive.value) {
    if (arr.value[i] !== "" || win.value) {
      return;
    }
    arr.value[i] = currentUser.value;

    const result = isWin(arr.value);
    if (result) {
      winner.value = result;
      win.value = true
      setActiveBox(9)
      addToarrayOfMain(props.ticId)
      setCurrentUser()
      return;
    }

    const draw = isDraw(arr.value)
    if (draw) {
      resetArr()
      setCurrentUser()
      setActiveBox(9)
      return
    }

    setCurrentUser()
    if (isAlreadyWinBySomeone(i)) {
      setActiveBox(9)
      return
    }
    setActiveBox(i)
  }
}

function resetArr() {
  arr.value = Array(9).fill("");
}
</script>

<style scoped>
.sub-grid-container {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  grid-template-rows: repeat(3, 1fr);
  gap: 6px;
  padding: 10px;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.15);
  border-radius: 16px;
  transition: all 0.4s cubic-bezier(0.175, 0.885, 0.32, 1.275);
  position: relative;
  width: 100%;
  height: 100%;
}

.active-grid {
  background: rgba(255, 255, 255, 0.1);
  border-color: var(--accent-x);
  box-shadow: 0 0 25px rgba(255, 255, 255, 0.1);
  transform: scale(1.03);
  z-index: 10;
}

.cell {
  aspect-ratio: 1;
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.5rem;
  font-weight: 800;
  cursor: pointer;
  transition: all 0.2s ease;
  user-select: none;
}

.cell:hover:not(.cell-occupied) {
  background: rgba(255, 255, 255, 0.15);
  border-color: rgba(255, 255, 255, 0.2);
}

.x-mark {
  color: var(--accent-x);
  text-shadow: 0 0 15px rgba(0, 242, 255, 0.6);
}

.o-mark {
  color: var(--accent-o);
  text-shadow: 0 0 15px rgba(255, 0, 234, 0.6);
}

.win-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.8);
  backdrop-filter: blur(6px);
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 20;
  pointer-events: none;
}

.win-overlay span {
  font-size: 4rem;
  opacity: 1;
}

/* Animations */
.pop-enter-active {
  animation: pop-in 0.3s cubic-bezier(0.175, 0.885, 0.32, 1.275);
}

@keyframes pop-in {
  0% { transform: scale(0); opacity: 0; }
  100% { transform: scale(1.2); opacity: 1; }
  100% { transform: scale(1); opacity: 1; }
}

.fade-enter-active, .fade-leave-active {
  transition: opacity 0.5s ease;
}

.fade-enter-from, .fade-leave-to {
  opacity: 0;
}

@media (max-width: 991px) {
  .sub-grid-container { padding: 6px; gap: 4px; }
  .cell { font-size: 1.1rem; }
  .win-overlay span { font-size: 2.5rem; }
}
</style>
