export interface User {
  id: string;
  first_name: string;
  last_name: string;
  email: string;
  roles: string;
}

export const useAuth = () => {
  const authUser = useState<User | null>('auth_user', () => null);
  const { $apiFetch } = useApiFetch();

  const isLoggedIn = computed(() => !!authUser.value);

  const fetchMe = async () => {
    try {
      const { data } = await $apiFetch<{ data: User }>('/api/v1/users/me', {
        headers: useRequestHeaders(['cookie'])
      });
      authUser.value = data;
    } catch (error) {
      authUser.value = null;
    }
  };

  const login = async (credentials: any) => {
    const { data } = await $apiFetch<{ data: User }>('/api/v1/login', {
      method: 'POST',
      body: credentials
    });
    authUser.value = data;
    return data;
  };

  const register = async (userData: any) => {
    const { data } = await $apiFetch<{ data: User }>('/api/v1/users', {
      method: 'POST',
      body: userData
    });
    return data;
  };

  const logout = async () => {
    await $apiFetch('/api/v1/logout');
    authUser.value = null;
    navigateTo('/login');
  };

  return {
    authUser,
    isLoggedIn,
    fetchMe,
    login,
    register,
    logout
  };
};
