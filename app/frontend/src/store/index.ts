import Vue from "vue";
import Vuex from "vuex";
import createPersistedState from "vuex-persistedstate";
import { fetchUser, logout as logoutReq } from '@/api/user';
import router from '@/router';
import halfmoon from "halfmoon";

Vue.use(Vuex);

export enum Authority {
  User,
}

export default new Vuex.Store({
  state: {
    darkMode: halfmoon.readCookie("darkModeOn") == "yes",
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
      try {
        const user = await fetchUser()
        commit('update', { user })
      } catch (error) {
        console.log('fetchUser', error)
        commit('update', { user: {} })
        if (router.currentRoute.path != "/") router.push("/");
      }
    },
    async logout({ commit }) {
      await logoutReq();
      commit('update', { user: {} })
      if (router.currentRoute.path != "/") router.push("/");
    },
    async toggleDarkMode({ commit }) {
      halfmoon.toggleDarkMode();
      commit('update', { darkMode: halfmoon.readCookie("darkModeOn") == "yes", })
    }
  },
  modules: {},
  plugins: [createPersistedState()],
});
