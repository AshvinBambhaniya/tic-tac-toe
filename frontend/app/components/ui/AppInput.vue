<script setup lang="ts">
defineProps<{
  modelValue: string;
  label?: string;
  type?: string;
  placeholder?: string;
  error?: string;
  required?: boolean;
}>();

defineEmits(['update:modelValue']);
</script>

<template>
  <div class="flex flex-col gap-2 w-full">
    <label v-if="label" class="text-[0.9rem] uppercase tracking-[2px] text-white/40 ml-1">
      {{ label }} {{ required ? '*' : '' }}
    </label>
    <input
      :type="type || 'text'"
      :value="modelValue"
      :placeholder="placeholder"
      :required="required"
      @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
      class="bg-white/5 border-2 border-white/10 rounded-xl px-4 py-3 text-white placeholder:text-white/20 focus:outline-none focus:border-accent-x/50 transition-all duration-300 w-full shadow-inner"
      :class="{ 'border-red-500/50': error }"
    />
    <span v-if="error" class="text-red-500 text-xs mt-1 ml-1">{{ error }}</span>
  </div>
</template>
