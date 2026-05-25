<template>
  <div class="bg-primary-bg min-h-screen text-white font-sans selection:bg-accent-x/30 selection:text-white overflow-x-hidden">
    <!-- Animated Background -->
    <div class="fixed inset-0 z-0 pointer-events-none overflow-hidden">
      <div class="absolute top-[-10%] left-[-10%] w-[40%] h-[40%] bg-accent-x/10 blur-[120px] rounded-full animate-pulse"></div>
      <div class="absolute bottom-[-10%] right-[-10%] w-[40%] h-[40%] bg-accent-o/10 blur-[120px] rounded-full animate-pulse" style="animation-delay: 2s"></div>
    </div>

    <!-- Main Content -->
    <div class="relative z-10">
      <NuxtPage />
    </div>
  </div>
</template>

<style>
.bg-primary-bg {
  background-color: #05070a;
}
</style>
