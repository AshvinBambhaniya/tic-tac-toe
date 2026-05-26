export const useApiFetch = () => {
  const config = useRuntimeConfig();
  const baseURL = import.meta.server ? config.apiUrl : config.public.apiUrl;

  // $apiFetch is for client-side calls (methods, event handlers)
  const $apiFetch = $fetch.create({
    baseURL,
    credentials: 'include',
    onResponseError({ response }) {
      console.error('API Error:', response._data);
    }
  });

  // apiFetch is a wrapper around useFetch for SSR-friendly calls in setup
  const apiFetch = (path: string, options: any = {}) => {
    return useFetch(path, {
      baseURL: `${baseURL}/api/v1`,
      ...options,
      credentials: 'include',
    });
  };

  return {
    $apiFetch,
    apiFetch
  };
};
