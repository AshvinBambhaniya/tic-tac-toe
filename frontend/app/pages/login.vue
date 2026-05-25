<script setup lang="ts">
import AppInput from '~/components/ui/AppInput.vue';
import AppButton from '~/components/ui/AppButton.vue';

const { login } = useAuth();

const form = reactive({
  email: '',
  password: ''
});

const loading = ref(false);
const error = ref('');

definePageMeta({
  middleware: 'auth'
});

const handleLogin = async () => {
  loading.value = true;
  error.value = '';
  try {
    await login(form);
    navigateTo('/');
  } catch (err: any) {
    error.value = err.data?.message || 'Invalid credentials. Please try again.';
  } finally {
    loading.value = false;
  }
};
</script>

<template>
  <div class="min-h-screen flex items-center justify-center p-4">
    <div class="w-full max-w-[450px] glass-effect rounded-[32px] p-10 md:p-8 border border-white/10 shadow-2xl">
      <div class="text-center mb-10">
        <h1 class="text-[2.5rem] font-extrabold mb-2 leading-tight">
          Welcome <span class="bg-gradient-to-r from-accent-x to-accent-o bg-clip-text text-transparent">Back</span>
        </h1>
        <p class="text-white/40 tracking-wide uppercase text-sm">Sign in to start your match</p>
      </div>

      <form @submit.prevent="handleLogin" class="flex flex-col gap-6">
        <AppInput
          v-model="form.email"
          label="Email Address"
          type="email"
          placeholder="Enter your email"
          required
        />
        
        <AppInput
          v-model="form.password"
          label="Password"
          type="password"
          placeholder="••••••••"
          required
        />

        <div v-if="error" class="bg-red-500/10 border border-red-500/20 text-red-500 p-4 rounded-xl text-sm text-center">
          {{ error }}
        </div>

        <AppButton type="submit" :loading="loading">
          Sign In
        </AppButton>
      </form>

      <div class="mt-8 text-center">
        <p class="text-white/40 text-sm">
          Don't have an account? 
          <NuxtLink to="/register" class="text-accent-x font-bold hover:underline ml-1">Create one</NuxtLink>
        </p>
      </div>
    </div>
  </div>
</template>
