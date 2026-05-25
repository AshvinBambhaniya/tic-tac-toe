<script setup lang="ts">
const { currentUser, mainArr, activeBox, resetFullGame } = useGameState();

const gameKey = ref(0);

const handleReset = () => {
  resetFullGame();
  gameKey.value++; 
};

useHead({
  title: 'Ultimate Tic-Tac-Toe',
  meta: [
    { name: 'description', content: 'A modern Ultimate Tic-Tac-Toe game built with Nuxt 4 and Tailwind CSS.' }
  ]
})
</script>

<template>
  <div class="min-h-screen flex items-center justify-center p-8 md:p-4">
    <div class="flex flex-col gap-16 max-w-[1400px] w-full xl:flex-row xl:items-center xl:gap-16">
      <!-- Left Section: Info Panel -->
      <div class="flex-none w-full max-w-[380px] mx-auto glass-effect rounded-[32px] p-12 md:p-8">
        <h1 class="text-[3rem] font-extrabold mb-10 leading-[1.1] text-left md:text-[2.2rem]">
          Ultimate <br><span class="bg-gradient-to-r from-accent-x to-accent-o bg-clip-text text-transparent">Tic-Tac-Toe</span>
        </h1>
        
        <div class="mb-10">
          <div 
            class="py-[1.2rem] rounded-[20px] font-extrabold text-[1.5rem] text-center transition-all duration-400 text-black"
            :class="currentUser === 'X' ? 'bg-accent-x turn-indicator-shadow-x' : 'bg-accent-o turn-indicator-shadow-o'"
          >
            {{ currentUser }}'s Turn
          </div>
        </div>

        <div class="mb-10">
          <p class="text-[0.9rem] uppercase tracking-[2px] text-white/40 mb-[1.2rem]">Global Board Status</p>
          <div class="grid grid-cols-3 gap-2.5 w-[130px]">
            <div 
              v-for="(item, index) in mainArr" 
              :key="index" 
              class="aspect-square bg-white/5 border rounded-lg flex items-center justify-center font-extrabold transition-all duration-300"
              :class="[
                item === 'X' ? 'text-accent-x' : item === 'O' ? 'text-accent-o' : '',
                activeBox === index ? 'border-accent-x bg-white/15 scale-110 shadow-[0_0_10px_rgba(0,242,255,0.3)]' : 'border-white/10'
              ]"
            >
              {{ item }}
            </div>
          </div>
        </div>

        <div class="mb-10">
          <p class="text-[0.9rem] uppercase tracking-[2px] text-white/40 mb-2">Active Zone</p>
          <div 
            class="text-[1.2rem] font-semibold text-white bg-white/10 py-[0.8rem] px-[1.5rem] rounded-xl inline-block transition-all duration-300"
            :class="{ 'border border-accent-x shadow-[0_0_15px_rgba(0,242,255,0.2)]': activeBox !== 9 }"
          >
            {{ activeBox === 9 ? 'Anywhere' : `Grid #${activeBox + 1}` }}
          </div>
        </div>

        <button 
          class="bg-transparent border-2 border-white/20 text-white py-[1.2rem] rounded-[20px] font-bold text-[1.1rem] cursor-pointer transition-all duration-300 w-full hover:bg-white/5 hover:border-white"
          @click="handleReset"
        >
          Restart Match
        </button>
      </div>

      <!-- Right Section: Game Board -->
      <div class="flex-1 flex justify-center items-center w-full">
        <div class="bg-white/5 border-2 border-white/10 rounded-[32px] p-6 w-full max-w-[850px] aspect-square flex md:max-w-[600px] md:p-3 md:rounded-[20px]">
          <div class="grid grid-cols-3 grid-rows-3 gap-5 w-full h-full md:gap-3" :key="gameKey">
            <Tic v-for="i in 9" :key="i-1" :ticId="i-1" />
          </div>
        </div>
      </div>
    </div>

    <GameOverModal @reset="handleReset" />
  </div>
</template>
