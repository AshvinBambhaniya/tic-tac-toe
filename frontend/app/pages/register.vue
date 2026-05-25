<script setup lang="ts">
import AppInput from '~/components/ui/AppInput.vue';
import AppButton from '~/components/ui/AppButton.vue';

const { register, login } = useAuth();

const form = reactive({
  first_name: '',
  last_name: '',
  email: '',
  password: '',
  roles: 'player'
});

const loading = ref(false);
const error = ref('');

definePageMeta({
  middleware: 'auth'
});

const handleRegister = async () => {
  loading.value = true;
  error.value = '';
  try {
    await register(form);
    // Automatically login after registration
    await login({
      email: form.email,
      password: form.password
    });
    navigateTo('/');
  } catch (err: any) {
    error.value = err.data?.message || 'Failed to create account. Please try again.';
  } finally {
    loading.value = false;
  }
};
</script>

<template>
  <div class="min-h-screen flex items-center justify-center p-4">
    <div class="w-full max-w-[500px] glass-effect rounded-[32px] p-10 md:p-8 border border-white/10 shadow-2xl">
      <div class="text-center mb-10">
        <h1 class="text-[2.5rem] font-extrabold mb-2 leading-tight">
          Join the <span class="bg-gradient-to-r from-accent-x to-accent-o bg-clip-text text-transparent">Arena</span>
        </h1>
        <p class="text-white/40 tracking-wide uppercase text-sm">Create your account to play</p>
      </div>

      <form @submit.prevent="handleRegister" class="flex flex-col gap-5">
        <div class="grid grid-cols-2 gap-4">
          <AppInput
            v-model="form.first_name"
            label="First Name"
            placeholder="John"
            required
          />
          <AppInput
            v-model="form.last_name"
            label="Last Name"
            placeholder="Doe"
            required
          />
        </div>
        
        <AppInput
          v-model="form.email"
          label="Email Address"
          type="email"
          placeholder="john@example.com"
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
          Create Account
        </AppButton>
      </form>

      <div class="mt-8 text-center">
        <p class="text-white/40 text-sm">
          Already have an account? 
          <NuxtLink to="/login" class="text-accent-o font-bold hover:underline ml-1">Sign In</NuxtLink>
        </p>
      </div>
    </div>
  </div>
</template>
