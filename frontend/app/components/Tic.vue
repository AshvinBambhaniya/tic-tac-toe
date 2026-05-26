<script setup lang="ts">
const props = defineProps<{
  ticId: number;
}>();

const { 
  activeBox, 
  currentUser, 
  isGameOver,
  allMoves,
  subGridResults,
  playerSymbol
} = useGameState();

const { sendMove } = useMultiplayer();

const isPlayerTurn = computed(() => {
  return playerSymbol.value === currentUser.value;
});

// Derived local array based on allMoves
const arr = computed(() => {
  const localArr = Array(9).fill("");
  allMoves.value.forEach(move => {
    if (move.sub_grid_index === props.ticId) {
      localArr[move.cell_index] = move.symbol;
    }
  });
  return localArr;
});

const subGridResult = computed(() => {
  return subGridResults.value.find(res => res.grid_index === props.ticId);
});

const win = computed(() => !!subGridResult.value && subGridResult.value.winner_symbol !== 'D');
const winner = computed(() => subGridResult.value?.winner_symbol || "");

const isGridActive = computed(() => {
  if (isGameOver.value) return false;
  // Box is active if it's not already won and is the target box
  return !subGridResult.value && (activeBox.value === 9 || activeBox.value === props.ticId);
});

const isDimmed = computed(() => {
  return !isGameOver.value && !subGridResult.value && activeBox.value !== 9 && activeBox.value !== props.ticId;
});

function addturn(i: number) {
  console.log("Attempting to add turn at index:", i);
  if (isGridActive.value && isPlayerTurn.value) {
    if (arr.value[i] !== "" || subGridResult.value) {
      return;
    }
    // In multiplayer, we send the move to the server
    sendMove(props.ticId, i);
  }
}
</script>

<template>
  <div 
    class="grid grid-cols-3 grid-rows-3 gap-1.5 p-2.5 rounded-2xl transition-all duration-500 ease-out relative w-full h-full lg:gap-1.5 lg:p-2.5 md:gap-1 md:p-1.5"
    :class="[
      isGridActive 
        ? 'bg-white/10 border-2 border-accent-x scale-[1.03] z-10 shadow-[0_0_40px_rgba(0,242,255,0.2)]' 
        : 'bg-white/5 border border-white/10',
      isDimmed || (!isPlayerTurn && !isGameOver) ? 'opacity-30 grayscale-[0.7]' : 'opacity-100 grayscale-0',
      win ? 'border-transparent' : ''
    ]"
  >
    <div 
      v-for="(item, index) in arr" 
      :key="index" 
      class="aspect-square bg-white/10 border border-white/5 rounded-lg flex items-center justify-center text-2xl font-extrabold transition-all duration-200 select-none md:text-lg"
      :class="[
        item ? '' : (isGridActive && isPlayerTurn ? 'cursor-pointer hover:bg-white/20 hover:border-white/20' : 'cursor-not-allowed opacity-50')
      ]"
      @click="addturn(index)"
    >
      <transition 
        enter-active-class="animate-pop-in"
      >
        <span 
          v-if="item" 
          :class="item === 'X' ? 'text-accent-x drop-shadow-[0_0_15px_rgba(0,242,255,0.6)]' : 'text-accent-o drop-shadow-[0_0_15px_rgba(255,0,234,0.6)]'"
        >
          {{ item }}
        </span>
      </transition>
    </div>
    
    <!-- Win Overlay -->
    <transition 
      enter-active-class="transition-opacity duration-500"
      enter-from-class="opacity-0"
      leave-active-class="transition-opacity duration-500"
      leave-to-class="opacity-0"
    >
      <div v-if="win" class="absolute inset-0 bg-black/80 backdrop-blur-md rounded-2xl flex items-center justify-center z-20 pointer-events-none border-2" :class="winner === 'X' ? 'border-accent-x/50' : 'border-accent-o/50'">
        <span 
          class="text-[4rem] md:text-[2.5rem]"
          :class="winner === 'X' ? 'text-accent-x drop-shadow-[0_0_20px_rgba(0,242,255,0.8)]' : 'text-accent-o drop-shadow-[0_0_20px_rgba(255,0,234,0.8)]'"
        >
          {{ winner }}
        </span>
      </div>
    </transition>
  </div>
</template>
