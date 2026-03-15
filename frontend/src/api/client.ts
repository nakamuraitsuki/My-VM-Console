import axios from 'axios';
import camelcaseKeys from 'camelcase-keys';
import snakecaseKeys from 'snakecase-keys';

// cf. https://axios-http.com/ja/docs/intro
export const apiClient = axios.create({
  baseURL: import.meta.env.VITE_BACKEND_URL || '',
  withCredentials: true,
  timeout: 30000,
})

// -- Request Interceptor --
apiClient.interceptors.request.use((config) => {
  // FormDataは触らない
  if (config.data && !(config.data instanceof FormData)) {
    // cf. https://github.com/bendrucker/snakecase-keys
    config.data = snakecaseKeys(config.data, { deep: true });
  }
  return config;
})

// -- Response Interceptor --
apiClient.interceptors.response.use(
  (response) => {
    // JSON以外触らない
    if (response.data && response.headers['content-type']?.includes('application/json')) {
      // cf. https://github.com/sindresorhus/camelcase-keys
      response.data = camelcaseKeys(response.data, { deep: true });
    }
    return response;
  },
  (error) => {
    const message = error.response?.data?.message || error.message;
    console.error(`[API Error] ${error.config.method?.toUpperCase()} ${error.config.url}`, message);

    // その他エラー時のリダイレクト処理などを追記
    return Promise.reject(error);
  }
);

