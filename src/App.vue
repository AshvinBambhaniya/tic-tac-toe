<template>
  <div class="main-container">
    <div class="content-wrapper">
      <!-- Left Section: Info Panel -->
      <div class="info-panel glass-card">
        <h1 class="game-title">Ultimate <br><span>Tic-Tac-Toe</span></h1>
        
        <div class="status-box">
          <div class="turn-indicator" :class="{ 'turn-o': currentUser === 'O' }">
            <span class="player-label">{{ currentUser }}'s Turn</span>
          </div>
        </div>

        <div class="main-board-preview">
          <p class="preview-label">Global Board Status</p>
          <div class="preview-grid">
            <div v-for="(item, index) in mainArr" :key="index" 
                 class="preview-cell"
                 :class="{ 'x-won': item === 'X', 'o-won': item === 'O' }">
              {{ item }}
            </div>
          </div>
        </div>

        <div class="game-rules">
          <p class="rules-title">Active Zone</p>
          <div class="zone-indicator">
            {{ activeBox === 9 ? 'Anywhere' : `Grid #${activeBox + 1}` }}
          </div>
        </div>

        <button class="reset-btn-outline" @click="handleReset">Restart Match</button>
      </div>

      <!-- Right Section: Game Board -->
      <div class="game-board-container">
        <div class="ultimate-board-wrapper">
          <div class="ultimate-grid" :key="gameKey">
            <Tic v-for="i in 9" :key="i-1" :ticId="i-1" />
          </div>
        </div>
      </div>
    </div>

    <GameOverModal @reset="handleReset" />
  </div>
</template>

<script setup>
import { ref } from 'vue';
import Tic from "./components/tic.vue";
import GameOverModal from "./components/GameOverModal.vue";
import { currentUser, resetCurrentUser } from './store/currentUser.js'
import { mainArr, resetMainArr } from "./store/mainArr.js";
import { activeBox, resetActiveBox } from "./store/activeBox.js";
import { resetGameStatus } from "./store/gameStatusStore.js";

const gameKey = ref(0);

const handleReset = () => {
  resetMainArr();
  resetCurrentUser();
  resetActiveBox();
  resetGameStatus();
  gameKey.value++; 
};
</script>

<style>
@import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;600;800&display=swap');

:root {
  --primary-bg: #042743;
  --glass-bg: rgba(255, 255, 255, 0.05);
  --glass-border: rgba(255, 255, 255, 0.15);
  --accent-x: #00f2ff;
  --accent-o: #ff00ea;
  --text-main: #ffffff;
}

* {
  box-sizing: border-box;
}

body {
  margin: 0;
  padding: 0;
  font-family: 'Inter', sans-serif;
  background: radial-gradient(circle at top right, #0f172a, #1e1b4b, #312e81);
  background-attachment: fixed;
  min-height: 100vh;
  color: var(--text-main);
  overflow-x: hidden;
}

.main-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 2rem;
}

.content-wrapper {
  display: flex;
  gap: 4rem;
  max-width: 1400px;
  width: 100%;
  flex-direction: column;
}

@media (min-width: 1200px) {
  .content-wrapper {
    flex-direction: row;
    align-items: center;
  }
}

.glass-card {
  background: rgba(255, 255, 255, 0.03);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid var(--glass-border);
  border-radius: 32px;
  padding: 3rem;
  box-shadow: 0 20px 50px rgba(0, 0, 0, 0.5);
}

.info-panel {
  flex: 0 0 auto;
  width: 100%;
  max-width: 380px;
  margin: 0 auto;
}

.game-title {
  font-size: 3rem;
  font-weight: 800;
  margin-bottom: 2.5rem;
  line-height: 1.1;
  text-align: left;
}

.game-title span {
  background: linear-gradient(to right, var(--accent-x), var(--accent-o));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.status-box {
  margin-bottom: 2.5rem;
}

.turn-indicator {
  background: var(--accent-x);
  color: #000;
  padding: 1.2rem;
  border-radius: 20px;
  font-weight: 800;
  font-size: 1.5rem;
  text-align: center;
  box-shadow: 0 0 30px rgba(0, 242, 255, 0.3);
  transition: all 0.4s ease;
}

.turn-indicator.turn-o {
  background: var(--accent-o);
  box-shadow: 0 0 30px rgba(255, 0, 234, 0.3);
}

.main-board-preview {
  margin-bottom: 2.5rem;
}

.preview-label {
  font-size: 0.9rem;
  text-transform: uppercase;
  letter-spacing: 2px;
  color: rgba(255, 255, 255, 0.4);
  margin-bottom: 1.2rem;
}

.preview-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 10px;
  width: 130px;
}

.preview-cell {
  aspect-ratio: 1;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 800;
}

.x-won { color: var(--accent-x); }
.o-won { color: var(--accent-o); }

.zone-indicator {
  font-size: 1.2rem;
  font-weight: 600;
  color: #fff;
  background: rgba(255, 255, 255, 0.1);
  padding: 0.8rem 1.5rem;
  border-radius: 12px;
  display: inline-block;
  margin-bottom: 2.5rem;
}

.game-board-container {
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
}

.ultimate-board-wrapper {
  background: rgba(255, 255, 255, 0.02);
  border: 2px solid rgba(255, 255, 255, 0.1);
  border-radius: 32px;
  padding: 24px;
  width: 100%;
  max-width: 850px;
  aspect-ratio: 1;
  display: flex;
}

.ultimate-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  grid-template-rows: repeat(3, 1fr);
  gap: 20px;
  width: 100%;
  height: 100%;
}

.reset-btn-outline {
  background: transparent;
  border: 2px solid rgba(255, 255, 255, 0.2);
  color: white;
  padding: 1.2rem;
  border-radius: 20px;
  font-weight: 700;
  font-size: 1.1rem;
  cursor: pointer;
  transition: all 0.3s;
  width: 100%;
}

.reset-btn-outline:hover {
  background: rgba(255, 255, 255, 0.05);
  border-color: white;
}

@media (max-width: 991px) {
  .ultimate-board-wrapper {
    max-width: 600px;
    padding: 12px;
  }
  .ultimate-grid {
    gap: 12px;
  }
}

@media (max-width: 576px) {
  .main-container { padding: 1rem; }
  .game-title { font-size: 2.2rem; }
  .glass-card { padding: 2rem; }
  .ultimate-board-wrapper { border-radius: 20px; }
}
</style>
