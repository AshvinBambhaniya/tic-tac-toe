export const useApiFetch = () => {
  const config = useRuntimeConfig();

  const baseURL = import.meta.server ? config.apiUrl : config.public.apiUrl;

  const $apiFetch = $fetch.create({
    baseURL,
    credentials: 'include',
    onResponseError({ response }) {
      // Handle global errors here if needed
      console.error('API Error:', response._data);
    }
  });

  return {
    $apiFetch
  };
};
