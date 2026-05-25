<script setup lang="ts">
const { isGameOver, gameWinner, gameDraw } = useGameState();

defineEmits(['reset']);
</script>

<template>
  <div v-if="isGameOver" class="fixed inset-0 bg-black/60 backdrop-blur-md flex justify-center items-center z-[1000]">
    <div class="bg-white/10 border border-white/20 rounded-[20px] p-12 text-center shadow-[0_8px_32px_0_rgba(31,38,135,0.37)] backdrop-blur-sm max-w-[90%]">
      <h2 class="text-white text-[3rem] mb-4 font-extrabold drop-shadow-[0_0_10px_rgba(255,255,255,0.5)]">
        {{ gameDraw ? "It's a Draw!" : `${gameWinner} Wins!` }}
      </h2>
      <p class="text-white/80 text-xl mb-8">
        {{ gameDraw ? "Well played both sides!" : "Congratulations to the winner!" }}
      </p>
      <button 
        class="bg-white text-[#042743] border-none py-3.5 px-10 text-[1.1rem] font-bold rounded-[50px] cursor-pointer transition-all duration-200 hover:-translate-y-0.5 hover:shadow-[0_5px_15px_rgba(255,255,255,0.4)] active:translate-y-0"
        @click="$emit('reset')"
      >
        Play Again
      </button>
    </div>
  </div>
</template>
