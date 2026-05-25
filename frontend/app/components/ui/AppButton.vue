<script setup lang="ts">
defineProps<{
  type?: 'button' | 'submit' | 'reset';
  variant?: 'primary' | 'secondary' | 'outline';
  loading?: boolean;
  disabled?: boolean;
}>();
</script>

<template>
  <button
    :type="type || 'button'"
    :disabled="disabled || loading"
    class="relative flex items-center justify-center py-4 px-8 rounded-xl font-bold text-[1.1rem] transition-all duration-300 w-full disabled:opacity-50 disabled:cursor-not-allowed group overflow-hidden"
    :class="[
      variant === 'secondary' ? 'bg-accent-o text-black hover:scale-[1.02] active:scale-95' :
      variant === 'outline' ? 'bg-transparent border-2 border-white/20 text-white hover:bg-white/5 hover:border-white' :
      'bg-accent-x text-black hover:scale-[1.02] active:scale-95'
    ]"
  >
    <div v-if="loading" class="absolute inset-0 flex items-center justify-center bg-inherit">
      <div class="w-6 h-6 border-2 border-black/30 border-t-black rounded-full animate-spin"></div>
    </div>
    <span :class="{ 'opacity-0': loading }">
      <slot />
    </span>
  </button>
</template>
