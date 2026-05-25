export default defineNuxtRouteMiddleware((to) => {
  const { isLoggedIn } = useAuth();

  if (!isLoggedIn.value && to.path !== '/login' && to.path !== '/register') {
    return navigateTo('/login');
  }

  if (isLoggedIn.value && (to.path === '/login' || to.path === '/register')) {
    return navigateTo('/');
  }
});
