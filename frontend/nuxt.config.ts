// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },

  future: {
    compatibilityVersion: 4,
  },

  modules: [
    '@nuxt/eslint',
    '@nuxtjs/color-mode',
    '@nuxtjs/tailwindcss'
  ],

  css: ['~/assets/css/main.css'],

  runtimeConfig: {
    apiUrl: process.env.NUXT_API_URL || 'http://backend:3000',

    public: {
      apiUrl: process.env.NUXT_PUBLIC_API_URL || 'http://localhost:3000',
    }
  },

  tailwindcss: {
    exposeConfig: true,
    viewer: true,
  }
})
