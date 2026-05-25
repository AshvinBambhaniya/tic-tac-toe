<script setup lang="ts">
const props = defineProps<{
  ticId: number;
}>();

const { 
  activeBox, 
  currentUser, 
  isGameOver, 
  toggleUser, 
  addToMainBoard, 
  isSubGridWon 
} = useGameState();

const arr = ref<string[]>(Array(9).fill(""));
const win = ref(false);
const winner = ref("");

const isGridActive = computed(() => {
  if (isGameOver.value) return false;
  // Box is active if:
  // 1. It's not already won
  // 2. AND (the global activeBox targets this specific grid OR any grid is allowed [activeBox === 9])
  return !win.value && (activeBox.value === 9 || activeBox.value === props.ticId);
});

const isDimmed = computed(() => {
  // Dim the box if the game is still going, it's not won yet, but it's not the active box
  return !isGameOver.value && !win.value && activeBox.value !== 9 && activeBox.value !== props.ticId;
});

function addturn(i: number) {
  if (isGridActive.value) {
    if (arr.value[i] !== "" || win.value) {
      return;
    }
    arr.value[i] = currentUser.value;

    const result = isWin(arr.value);
    if (result) {
      winner.value = result;
      win.value = true;
      
      // Update global board status
      addToMainBoard(props.ticId);
      
      // Determine next active box:
      // If the clicked cell (i) points to a sub-grid that is already won globally, set activeBox to 9
      if (isSubGridWon(i)) {
        activeBox.value = 9;
      } else {
        activeBox.value = i;
      }
      
      toggleUser();
      return;
    }

    const draw = isDraw(arr.value);
    if (draw) {
      resetArr();
      toggleUser();
      activeBox.value = 9;
      return;
    }

    // Normal move: toggle user and set next active box
    toggleUser();
    
    if (isSubGridWon(i)) {
      activeBox.value = 9;
    } else {
      activeBox.value = i;
    }
  }
}

function resetArr() {
  arr.value = Array(9).fill("");
}
</script>

<template>
  <div 
    class="grid grid-cols-3 grid-rows-3 gap-1.5 p-2.5 rounded-2xl transition-all duration-500 ease-out relative w-full h-full lg:gap-1.5 lg:p-2.5 md:gap-1 md:p-1.5"
    :class="[
      isGridActive 
        ? 'bg-white/10 border-2 border-accent-x scale-[1.03] z-10 shadow-[0_0_40px_rgba(0,242,255,0.2)]' 
        : 'bg-white/5 border border-white/10',
      isDimmed ? 'opacity-30 grayscale-[0.7]' : 'opacity-100 grayscale-0',
      win ? 'border-transparent' : ''
    ]"
  >
    <div 
      v-for="(item, index) in arr" 
      :key="index" 
      class="aspect-square bg-white/10 border border-white/5 rounded-lg flex items-center justify-center text-2xl font-extrabold cursor-pointer transition-all duration-200 select-none hover:bg-white/20 hover:border-white/20 md:text-lg"
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
