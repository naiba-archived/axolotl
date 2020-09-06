import axios from 'axios';
import router from '@/router';
import store from '@/store';

axios.interceptors.request.use(function (config) {
    return {
        ...config,
        params: {
            ...config.params,
        },
    };
});

axios.interceptors.response.use(function (response) {
    if (response.status !== 200) {
        throw new Error(response.statusText);
    }
    const { code, msg } = response.data;
    if (!code) return response.data.data || {}; // 请求正常
    if (code == 10001) { // 未授权
        store.dispatch('logout')
        return {};
    }
    throw new Error(msg);
});

export default axios;