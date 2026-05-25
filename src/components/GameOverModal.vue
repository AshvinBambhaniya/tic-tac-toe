<template>
  <div v-if="isGameOver" class="modal-overlay">
    <div class="glass-modal">
      <h2 class="title">{{ gameDraw ? "It's a Draw!" : `${gameWinner} Wins!` }}</h2>
      <div class="message">
        {{ gameDraw ? "Well played both sides!" : "Congratulations to the winner!" }}
      </div>
      <button class="reset-btn" @click="$emit('reset')">Play Again</button>
    </div>
  </div>
</template>

<script setup>
import { isGameOver, gameWinner, gameDraw } from '../store/gameStatusStore.js';

defineEmits(['reset']);
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(8px);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.glass-modal {
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 20px;
  padding: 3rem;
  text-align: center;
  box-shadow: 0 8px 32px 0 rgba(31, 38, 135, 0.37);
  backdrop-filter: blur(4px);
  -webkit-backdrop-filter: blur(4px);
  max-width: 90%;
}

.title {
  color: white;
  font-size: 3rem;
  margin-bottom: 1rem;
  font-weight: 800;
  text-shadow: 0 0 10px rgba(255, 255, 255, 0.5);
}

.message {
  color: rgba(255, 255, 255, 0.8);
  font-size: 1.2rem;
  margin-bottom: 2rem;
}

.reset-btn {
  background: white;
  color: #042743;
  border: none;
  padding: 0.8rem 2.5rem;
  font-size: 1.1rem;
  font-weight: 700;
  border-radius: 50px;
  cursor: pointer;
  transition: transform 0.2s, box-shadow 0.2s;
}

.reset-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(255, 255, 255, 0.4);
}

.reset-btn:active {
  transform: translateY(0);
}
</style>
