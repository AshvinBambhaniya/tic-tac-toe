<script setup lang="ts">
import confetti from 'canvas-confetti';

const { isGameOver, gameWinner, gameDraw } = useGameState();
const isDismissed = ref(false);

defineEmits(['reset']);

const triggerConfetti = () => {
  if (gameDraw.value) return;

  const count = 200;
  const defaults = {
    origin: { y: 0.7 },
    zIndex: 2000,
  };

  function fire(particleRatio: number, opts: any) {
    confetti({
      ...defaults,
      ...opts,
      particleCount: Math.floor(count * particleRatio),
    });
  }

  fire(0.25, { spread: 26, startVelocity: 55 });
  fire(0.2, { spread: 60 });
  fire(0.35, { spread: 100, decay: 0.91, scalar: 0.8 });
  fire(0.1, { spread: 120, startVelocity: 25, decay: 0.92, scalar: 1.2 });
  fire(0.1, { spread: 120, startVelocity: 45 });
};

watch(isGameOver, (val) => {
  if (val) {
    isDismissed.value = false;
    setTimeout(triggerConfetti, 300);
  }
}, { immediate: true });
</script>

<template>
  <transition 
    enter-active-class="transition duration-700 ease-out"
    enter-from-class="transform translate-y-full opacity-0"
    enter-to-class="transform translate-y-0 opacity-100"
    leave-active-class="transition duration-500 ease-in"
    leave-from-class="transform translate-y-0 opacity-100"
    leave-to-class="transform translate-y-full opacity-0"
  >
    <div v-if="isGameOver && !isDismissed" class="fixed bottom-8 left-0 right-0 z-[1000] px-4 pointer-events-none">
      <div class="max-w-xl mx-auto glass-effect border border-white/20 rounded-[32px] p-8 text-center shadow-[0_20px_50px_rgba(0,0,0,0.5)] pointer-events-auto relative overflow-hidden group">
        <!-- Background Glow -->
        <div class="absolute inset-0 bg-gradient-to-t from-accent-x/5 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-700"></div>

        <div class="relative z-10">
          <div class="flex items-center justify-center gap-4 mb-2">
             <div class="w-12 h-[2px] bg-gradient-to-r from-transparent to-white/20"></div>
             <span class="text-[0.7rem] font-black uppercase tracking-[0.4em] text-white/40">Match Concluded</span>
             <div class="w-12 h-[2px] bg-gradient-to-l from-transparent to-white/20"></div>
          </div>

          <h2 class="text-4xl md:text-5xl font-black mb-2 tracking-tighter"
              :class="gameDraw ? 'text-white' : (gameWinner === 'X' ? 'text-accent-x' : 'text-accent-o')">
            {{ gameDraw ? "IT'S A DRAW" : `PLAYER ${gameWinner} WINS` }}
          </h2>
          
          <p class="text-white/60 font-medium mb-8 text-sm">
            {{ gameDraw ? "A perfectly balanced strategic battle!" : "Ultimate dominance achieved on the meta-board." }}
          </p>

          <div class="flex flex-col sm:flex-row gap-4 justify-center">
            <button 
              @click="$emit('reset')"
              class="bg-white text-black px-8 py-4 rounded-2xl font-black text-xs uppercase tracking-widest hover:scale-105 transition-all shadow-xl"
            >
              PLAY AGAIN
            </button>
            <button 
              @click="isDismissed = true"
              class="bg-white/5 border border-white/10 text-white px-8 py-4 rounded-2xl font-black text-xs uppercase tracking-widest hover:bg-white/10 transition-all"
            >
              REVIEW BOARD
            </button>
          </div>
        </div>
      </div>
    </div>
  </transition>
</template>

<style scoped>
.glass-effect {
  background: rgba(15, 15, 20, 0.8);
  backdrop-filter: blur(25px);
}
</style>
