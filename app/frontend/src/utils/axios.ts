import axios from "axios";

axios.interceptors.request.use(function (config) {
  return {
    ...config,
    params: {
      ...config.params
    }
  };
});

axios.interceptors.response.use(function (response) {
  if (response.status !== 200) {
    throw new Error(response.statusText);
  }
  const { code } = response.data;
  if (!code) return response.data.data || {}; // 请求正常
  throw new Error(JSON.stringify(response.data));
});

export default axios;
