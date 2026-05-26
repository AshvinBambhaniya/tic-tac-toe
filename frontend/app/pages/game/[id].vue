<script setup lang="ts">
const route = useRoute();
const gameId = route.params.id as string;

const { 
  currentUser, 
  mainArr, 
  activeBox, 
  resetFullGame,
  isGameOver,
  gameWinner,
  gameDraw,
  playerSymbol,
  isOpponentDisconnected,
  updateFromServer
} = useGameState();

const { connect, disconnect, sendForfeit } = useMultiplayer();
const { authUser, logout } = useAuth();
const { $apiFetch } = useApiFetch();

const isReviewMode = ref(false);
const winnerName = ref('');

definePageMeta({
  middleware: 'auth'
});

onMounted(async () => {
  resetFullGame();
  
  try {
    // 1. Fetch initial state via REST
    const { data } = await $apiFetch<any>(`/api/v1/games/${gameId}`);
    if (data) {
      updateFromServer(data);
      
      // 2. If game is ongoing, connect WebSockets
      if (data.game.status === 'ongoing') {
        connect(gameId);
      } else {
        // 3. If finished, it's Review Mode
        isReviewMode.value = true;
        if (data.game.winner_id) {
           const isXWinner = data.game.winner_id === data.game.player_x_id;
           winnerName.value = isXWinner ? 'Player X' : 'Player O';
        }
      }
    }
  } catch (err) {
    console.error('Failed to load game state:', err);
  }
});

onUnmounted(() => {
  disconnect();
});

const handleReset = () => {
  if (!isGameOver.value) {
    if (confirm('Are you sure you want to quit? This will forfeit the match.')) {
      sendForfeit();
    } else {
      return;
    }
  }
  useRouter().push('/');
};

useHead({
  title: 'Ultimate Tic-Tac-Toe - Match',
})
</script>

<template>
  <div class="min-h-screen flex items-center justify-center p-8 md:p-4 relative">
    <!-- Disconnection Notification -->
    <transition 
      enter-active-class="transition duration-500 ease-out"
      enter-from-class="transform -translate-y-full opacity-0"
      enter-to-class="transform translate-y-0 opacity-100"
      leave-active-class="transition duration-300 ease-in"
      leave-from-class="transform translate-y-0 opacity-100"
      leave-to-class="transform -translate-y-full opacity-0"
    >
      <div v-if="isOpponentDisconnected && !isGameOver" class="fixed top-0 left-0 right-0 z-[100] p-4 flex justify-center">
        <div class="glass-effect border border-accent-o/50 px-8 py-4 rounded-full flex items-center gap-4 shadow-[0_0_30px_rgba(255,0,234,0.2)]">
          <div class="w-3 h-3 bg-accent-o rounded-full animate-ping"></div>
          <span class="text-white font-bold tracking-wide">Opponent disconnected. Waiting for them to rejoin...</span>
          <button @click="handleReset" class="ml-4 text-[0.7rem] uppercase tracking-widest bg-white/10 hover:bg-white/20 px-4 py-2 rounded-full transition-all">Quit Match</button>
        </div>
      </div>
    </transition>

    <!-- Header / Nav -->
    <div class="absolute top-8 left-8 z-50">
      <NuxtLink to="/" class="flex items-center gap-2 text-white/50 hover:text-white transition-all font-bold uppercase tracking-widest text-xs">
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><path d="m15 18-6-6 6-6"/></svg>
        Back to Lobby
      </NuxtLink>
    </div>

    <!-- User Profile -->
    <div v-if="authUser" class="absolute top-8 right-8 z-50">
      <div class="flex items-center gap-4 glass-effect p-2 pr-4 rounded-full border border-white/10 hover:border-white/30 transition-all duration-300">
        <div class="w-10 h-10 rounded-full bg-accent-x flex items-center justify-center font-bold text-black">
          {{ authUser.first_name[0] }}
        </div>
        <div class="flex flex-col mr-4">
          <span class="text-sm font-bold text-white">{{ authUser.first_name }} {{ authUser.last_name }}</span>
          <span class="text-[0.7rem] text-white/50 uppercase tracking-wider">{{ authUser.roles }}</span>
        </div>
        <button 
          @click="logout"
          class="p-2 rounded-full hover:bg-white/10 text-white/50 hover:text-white transition-all duration-300"
        >
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/><polyline points="16 17 21 12 16 7"/><line x1="21" y1="12" x2="9" y2="12"/></svg>
        </button>
      </div>
    </div>

    <div class="flex flex-col gap-16 max-w-[1400px] w-full xl:flex-row xl:items-center xl:gap-16">
      <!-- Left Section: Info Panel -->
      <div class="flex-none w-full max-w-[380px] mx-auto glass-effect rounded-[32px] p-12 md:p-8">
        <div class="mb-8">
           <p class="text-[0.7rem] uppercase tracking-[3px] text-white/30 mb-1 font-bold">MATCH ID</p>
           <p class="text-[0.8rem] font-mono text-white/60 truncate">{{ gameId }}</p>
        </div>

        <h1 class="text-[3rem] font-extrabold mb-10 leading-[1.1] text-left md:text-[2.2rem]">
          Ultimate <br><span class="bg-gradient-to-r from-accent-x to-accent-o bg-clip-text text-transparent">Tic-Tac-Toe</span>
        </h1>
        
        <div class="mb-10">
          <div v-if="!isReviewMode"
            class="py-[1.2rem] rounded-[20px] font-extrabold text-[1.5rem] text-center transition-all duration-400 text-black relative overflow-hidden"
            :class="currentUser === 'X' ? 'bg-accent-x turn-indicator-shadow-x' : 'bg-accent-o turn-indicator-shadow-o'"
          >
            {{ currentUser }}'s Turn
            <div v-if="(currentUser === 'X' && playerSymbol === 'X') || (currentUser === 'O' && playerSymbol === 'O')" class="absolute top-0 right-0 bg-black/10 px-3 py-1 text-[0.6rem] uppercase tracking-tighter">Your Turn</div>
          </div>
          <div v-else
            class="py-[1.2rem] rounded-[20px] font-extrabold text-[1.5rem] text-center transition-all duration-400 text-white relative overflow-hidden bg-white/10"
          >
            <span class="text-white/40 uppercase tracking-widest text-[0.6rem] block mb-1">Final Result</span>
            <span v-if="gameDraw" class="text-white">MATCH DRAW</span>
            <span v-else :class="gameWinner === 'X' ? 'text-accent-x' : 'text-accent-o'">{{ winnerName }} WINS</span>
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
          @click="isReviewMode ? $router.push('/profile') : handleReset()"
        >
          {{ isReviewMode ? 'Back to Profile' : 'Quit Match' }}
        </button>
      </div>

      <!-- Right Section: Game Board -->
      <div class="flex-1 flex justify-center items-center w-full">
        <div class="bg-white/5 border-2 border-white/10 rounded-[32px] p-6 w-full max-w-[850px] aspect-square flex md:max-w-[600px] md:p-3 md:rounded-[20px]">
          <div class="grid grid-cols-3 grid-rows-3 gap-5 w-full h-full md:gap-3">
            <Tic v-for="i in 9" :key="i-1" :ticId="i-1" />
          </div>
        </div>
      </div>
    </div>

    <GameOverModal v-if="isGameOver && !isReviewMode" @reset="handleReset" />
  </div>
</template>
