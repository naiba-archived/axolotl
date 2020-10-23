import Vue from "vue";
import Vuex from "vuex";
import createPersistedState from "vuex-persistedstate";
import { fetchUser, logout as logoutReq } from "@/api/user";
import router from "@/router";
import halfmoon from "halfmoon";
import { getConfig } from '@/api/config';

Vue.use(Vuex);

export enum Authority {
  User
}

export default new Vuex.Store({
  state: {
    darkMode: halfmoon.readCookie("darkModeOn") == "yes",
    user: {
      id: Number,
      nickname: String,
      githubId: Number,
      authority: Authority
    },
    config: {},
  },
  mutations: {
    update(state, payload) {
      Object.assign(state, payload);
    }
  },
  actions: {
    async fetchConfig({ commit }) {
      const config = await getConfig();
      commit("update", { config });
    },
    async fetchUser({ commit }) {
      try {
        const user = await fetchUser();
        commit("update", { user });
      } catch (error) {
        commit('update', { user: {} })
        if (window.location.pathname != "/") {
          localStorage.setItem("returnURL", window.location.href)
          router.push("/");
        }
      }
    },
    async logout({ commit }) {
      await logoutReq();
      commit("update", { user: {} });
      if (router.currentRoute.path != "/") router.push("/");
    },
    async toggleDarkMode({ commit }) {
      halfmoon.toggleDarkMode();
      commit("update", {
        darkMode: halfmoon.readCookie("darkModeOn") == "yes"
      });
    }
  },
  modules: {},
  plugins: [createPersistedState()]
});
