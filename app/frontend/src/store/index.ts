import Vue from "vue";
import Vuex from "vuex";
import createPersistedState from "vuex-persistedstate";
import { fetchUser, logout } from '@/api/user';
import router from '@/router';

Vue.use(Vuex);

export enum Authority {
  User,
}

export default new Vuex.Store({
  state: {
    user: {
      id: Number,
      nickname: String,
      githubId: Number,
      authority: Authority,
    }
  },
  mutations: {
    update(state, payload) {
      Object.assign(state, payload);
    },
  },
  actions: {
    async fetchUser({ commit }) {
      const user = await fetchUser()
      commit('update', { user })
    },
    async logout({ commit }) {
      await logout();
      commit('update', { user: {} })
      if (router.currentRoute.path != "/") router.push("/");
    }
  },
  modules: {},
  plugins: [createPersistedState()],
});
