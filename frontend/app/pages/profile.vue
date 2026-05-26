<script setup lang="ts">
const { authUser, logout } = useAuth();
const { $apiFetch } = useApiFetch();
const router = useRouter();

definePageMeta({
  middleware: 'auth'
});

const profileData = ref<any>(null);
const isLoading = ref(true);

onMounted(async () => {
  try {
    const { data } = await $apiFetch<any>('/api/v1/games/profile');
    profileData.value = data;
  } catch (err) {
    console.error('Failed to fetch profile:', err);
  } finally {
    isLoading.value = false;
  }
});

const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  });
};

const getResultLabel = (game: any) => {
  if (game.status === 'draw') return { text: 'DRAW', class: 'text-white/40 bg-white/5' };
  
  const isWinner = game.winner_id === authUser.value?.id;
  if (isWinner) return { text: 'WIN', class: 'text-accent-x bg-accent-x/10 shadow-[0_0_15px_rgba(0,242,255,0.1)]' };
  return { text: 'LOSS', class: 'text-accent-o bg-accent-o/10 shadow-[0_0_15px_rgba(255,0,234,0.1)]' };
};

const getOpponentName = (game: any) => {
  const isPlayerX = game.player_x_id === authUser.value?.id;
  return isPlayerX ? (game.player_o_name || 'Waiting...') : game.player_x_name;
};

useHead({
  title: 'Ultimate Tic-Tac-Toe - Profile',
})
</script>

<template>
  <div class="min-h-screen p-8 relative overflow-hidden bg-[#05070a]">
    <!-- Background Accents -->
    <div class="absolute top-[-10%] left-[-10%] w-[40%] h-[40%] bg-accent-x/5 rounded-full blur-[120px] pointer-events-none"></div>
    <div class="absolute bottom-[-10%] right-[-10%] w-[40%] h-[40%] bg-accent-o/5 rounded-full blur-[120px] pointer-events-none"></div>

    <div class="max-w-6xl mx-auto z-10 relative">
      <!-- Navigation -->
      <div class="mb-12 flex items-center justify-between">
        <NuxtLink to="/" class="flex items-center gap-2 text-white/50 hover:text-white transition-all font-bold uppercase tracking-widest text-xs group">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round" class="group-hover:-translate-x-1 transition-transform"><path d="m15 18-6-6 6-6"/></svg>
          Back to Lobby
        </NuxtLink>

        <div v-if="authUser" class="flex items-center gap-4 glass-effect p-1.5 pr-4 rounded-full border border-white/10">
          <div class="w-8 h-8 rounded-full bg-accent-x flex items-center justify-center font-black text-black text-xs">
            {{ authUser.first_name[0] }}
          </div>
          <span class="text-sm font-bold text-white/80">{{ authUser.first_name }} {{ authUser.last_name }}</span>
          <button @click="logout" class="p-1.5 text-white/30 hover:text-white transition-colors">
             <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/><polyline points="16 17 21 12 16 7"/><line x1="21" y1="12" x2="9" y2="12"/></svg>
          </button>
        </div>
      </div>

      <div v-if="isLoading" class="flex flex-col items-center justify-center py-32">
        <div class="w-12 h-12 border-4 border-accent-x/20 border-t-accent-x rounded-full animate-spin mb-4"></div>
        <p class="text-white/40 font-bold uppercase tracking-widest text-xs">Loading Profile...</p>
      </div>

      <div v-else-if="profileData" class="space-y-12">
        <!-- Profile Header -->
        <div class="text-center">
          <h1 class="text-5xl font-black mb-4 tracking-tighter text-white">Player <span class="bg-gradient-to-r from-accent-x to-accent-o bg-clip-text text-transparent">Profile</span></h1>
          <p class="text-white/40 font-medium">Tracking your ultimate performance</p>
        </div>

        <!-- Stats Grid -->
        <div class="grid grid-cols-2 md:grid-cols-5 gap-6">
          <div class="glass-effect p-6 rounded-[24px] border border-white/5 text-center">
            <p class="text-[0.6rem] font-black uppercase tracking-[0.2em] text-white/30 mb-2">Total Matches</p>
            <p class="text-3xl font-black text-white">{{ profileData.total_games }}</p>
          </div>
          <div class="glass-effect p-6 rounded-[24px] border border-white/5 text-center">
            <p class="text-[0.6rem] font-black uppercase tracking-[0.2em] text-white/30 mb-2">Wins</p>
            <p class="text-3xl font-black text-accent-x">{{ profileData.wins }}</p>
          </div>
          <div class="glass-effect p-6 rounded-[24px] border border-white/5 text-center">
            <p class="text-[0.6rem] font-black uppercase tracking-[0.2em] text-white/30 mb-2">Losses</p>
            <p class="text-3xl font-black text-accent-o">{{ profileData.losses }}</p>
          </div>
          <div class="glass-effect p-6 rounded-[24px] border border-white/5 text-center">
            <p class="text-[0.6rem] font-black uppercase tracking-[0.2em] text-white/30 mb-2">Draws</p>
            <p class="text-3xl font-black text-white/60">{{ profileData.draws }}</p>
          </div>
          <div class="glass-effect p-6 rounded-[24px] border border-white/5 text-center col-span-2 md:col-span-1">
            <p class="text-[0.6rem] font-black uppercase tracking-[0.2em] text-white/30 mb-2">Win Rate</p>
            <p class="text-3xl font-black text-white">{{ profileData.win_rate.toFixed(1) }}%</p>
          </div>
        </div>

        <!-- Match History -->
        <div class="glass-effect rounded-[32px] border border-white/10 overflow-hidden">
          <div class="p-8 border-b border-white/5 flex items-center justify-between">
             <h2 class="text-sm font-black uppercase tracking-[0.3em] text-white">Match History</h2>
             <span class="text-[0.6rem] font-black text-white/20 uppercase tracking-widest">{{ profileData.history.length }} matches recorded</span>
          </div>

          <div v-if="profileData.history.length === 0" class="py-20 text-center">
             <p class="text-white/20 font-bold uppercase tracking-widest text-xs">No matches played yet.</p>
             <NuxtLink to="/" class="mt-4 inline-block text-accent-x hover:underline font-bold text-xs">Go play your first match!</NuxtLink>
          </div>

          <div v-else class="divide-y divide-white/5">
            <div v-for="game in profileData.history" :key="game.id" class="p-6 md:p-8 hover:bg-white/[0.02] transition-colors group">
              <div class="flex flex-col md:flex-row md:items-center gap-6">
                <!-- Status Badge -->
                <div class="flex-none">
                  <div :class="getResultLabel(game).class" class="px-4 py-2 rounded-xl text-xs font-black tracking-widest text-center w-24">
                    {{ getResultLabel(game).text }}
                  </div>
                </div>

                <!-- Opponent Info -->
                <div class="flex-1 min-w-0">
                  <p class="text-[0.6rem] font-black uppercase tracking-widest text-white/30 mb-1">Opponent</p>
                  <p class="text-lg font-bold text-white truncate group-hover:text-accent-x transition-colors">{{ getOpponentName(game) }}</p>
                </div>

                <!-- Match Details -->
                <div class="flex-none text-left md:text-right">
                  <p class="text-[0.6rem] font-black uppercase tracking-widest text-white/30 mb-1">Date Played</p>
                  <p class="text-sm font-medium text-white/60">{{ formatDate(game.updated_at) }}</p>
                </div>

                <!-- Action -->
                <div class="flex-none">
                  <NuxtLink 
                    :to="`/game/${game.id}`"
                    class="bg-white/5 hover:bg-white/10 text-white font-bold py-3 px-6 rounded-xl text-xs transition-all border border-white/5 hover:border-white/20 flex items-center gap-2"
                  >
                    REVIEW MATCH
                    <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><path d="M5 12h14"/><path d="m12 5 7 7-7 7"/></svg>
                  </NuxtLink>
                </div>
              </div>
            </div>
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
