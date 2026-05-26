<script setup lang="ts">
const { authUser, logout } = useAuth();
const router = useRouter();
const config = useRuntimeConfig();
const { connect, disconnect } = useMultiplayer();
const { $apiFetch } = useApiFetch();

definePageMeta({
  middleware: 'auth'
});

const gameIdToJoin = ref('');
const isMatching = ref(false);
const activeGames = ref<any[]>([]);

onMounted(async () => {
  try {
    const { data } = await $apiFetch<any>('/api/v1/games/active');
    activeGames.value = data || [];
  } catch (err) {
    console.error('Failed to fetch active games:', err);
  }
});

const createGame = async () => {
  try {
    const data = await $apiFetch<any>('/api/v1/games', {
      method: 'POST'
    });
    if (data) {
      router.push(`/game/${data.data.id}`);
    }
  } catch (err) {
    console.error('Failed to create game:', err);
  }
};

const joinGame = () => {
  if (gameIdToJoin.value) {
    router.push(`/game/${gameIdToJoin.value}`);
  }
};

const startMatchmaking = async () => {
  isMatching.value = true;
  
  try {
    // 1. Connect to lobby socket first and WAIT for it to be open
    await connect('lobby', (newGameId) => {
      disconnect();
      router.push(`/game/${newGameId}`);
    });

    // 2. Only after socket is open, tell the server to add us to the queue
    await $apiFetch('/api/v1/games/matchmake', {
      method: 'POST'
    });
  } catch (err) {
    console.error('Matchmaking failed:', err);
    isMatching.value = false;
    disconnect();
  }
};

onUnmounted(() => {
  disconnect();
});

useHead({
  title: 'Ultimate Tic-Tac-Toe - Lobby',
})
</script>

<template>
  <div class="min-h-screen flex items-center justify-center p-8 relative overflow-hidden bg-[#0a0a0c]">
    <!-- Background Accents -->
    <div class="absolute top-[-10%] left-[-10%] w-[40%] h-[40%] bg-accent-x/10 rounded-full blur-[120px] pointer-events-none"></div>
    <div class="absolute bottom-[-10%] right-[-10%] w-[40%] h-[40%] bg-accent-o/10 rounded-full blur-[120px] pointer-events-none"></div>

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

    <div class="max-w-2xl w-full z-10">
      <div class="text-center mb-16">
        <h1 class="text-6xl font-black mb-4 tracking-tighter">
          Ultimate <span class="bg-gradient-to-r from-accent-x to-accent-o bg-clip-text text-transparent">Tic-Tac-Toe</span>
        </h1>
        <p class="text-white/40 text-xl font-medium">Multiplayer Edition</p>
      </div>

      <div class="grid grid-cols-1 gap-6">
        <!-- Rejoin Active Games -->
        <div v-if="activeGames.length > 0" class="glass-effect rounded-[32px] p-8 border border-accent-x/20 space-y-4">
          <h2 class="text-xs font-black uppercase tracking-[0.3em] text-accent-x text-center mb-6">Active Matches Found</h2>
          <div v-for="game in activeGames" :key="game.id" class="flex items-center justify-between bg-white/5 p-4 rounded-2xl border border-white/5 hover:border-accent-x/30 transition-all">
            <div class="flex flex-col">
              <span class="text-[0.6rem] text-white/40 uppercase tracking-widest font-bold">Match ID</span>
              <span class="text-sm font-mono text-white/80">{{ game.id.slice(0, 8) }}...</span>
            </div>
            <button 
              @click="router.push(`/game/${game.id}`)"
              class="bg-accent-x/10 hover:bg-accent-x text-accent-x hover:text-black px-6 py-2 rounded-xl text-xs font-black transition-all border border-accent-x/20"
            >
              REJOIN
            </button>
          </div>
        </div>

        <!-- Play Options -->
        <div class="glass-effect rounded-[32px] p-8 border border-white/10 space-y-8">
          <div class="space-y-4">
            <button 
              @click="createGame"
              class="w-full bg-accent-x text-black py-5 rounded-2xl font-black text-xl hover:scale-[1.02] transition-all duration-300 shadow-[0_0_30px_rgba(0,242,255,0.3)] flex items-center justify-center gap-3"
            >
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><path d="M5 12h14"/><path d="M12 5v14"/></svg>
              CREATE PRIVATE MATCH
            </button>
            <p class="text-white/30 text-center text-sm">Create a game and invite a friend with the code</p>
          </div>

          <div class="relative py-4 flex items-center">
            <div class="flex-grow border-t border-white/10"></div>
            <span class="flex-shrink mx-4 text-white/20 font-bold uppercase tracking-widest text-xs">OR</span>
            <div class="flex-grow border-t border-white/10"></div>
          </div>

          <div class="space-y-4">
            <div class="flex gap-3">
              <input 
                v-model="gameIdToJoin"
                type="text" 
                placeholder="ENTER GAME CODE" 
                class="flex-1 bg-white/5 border border-white/10 rounded-2xl px-6 py-5 text-white font-bold tracking-widest focus:border-accent-o outline-none transition-all"
              />
              <button 
                @click="joinGame"
                class="bg-accent-o text-black px-8 rounded-2xl font-black transition-all hover:scale-[1.05]"
              >
                JOIN
              </button>
            </div>
          </div>

          <div class="pt-4">
            <button 
              @click="startMatchmaking"
              :disabled="isMatching"
              class="w-full border-2 border-white/10 text-white py-5 rounded-2xl font-black text-xl hover:bg-white/5 hover:border-white/30 transition-all flex items-center justify-center gap-3"
              :class="{'opacity-50 cursor-not-allowed': isMatching}"
            >
              <svg v-if="!isMatching" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 12a9 9 0 0 1-9 9m9-9a9 9 0 0 0-9-9m9 9H3m9 9a9 9 0 0 1-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9"/></svg>
              <span v-if="isMatching" class="animate-pulse">FINDING OPPONENT...</span>
              <span v-else>QUICK MATCH (ONLINE)</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.glass-effect {
  background: rgba(255, 255, 255, 0.03);
  backdrop-filter: blur(20px);
}
</style>
