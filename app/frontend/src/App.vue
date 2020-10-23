<template>
  <div>
    <Navbar />
    <router-view />
    <nav class="navbar">
      <div class="container">
        <ul class="navbar-nav ml-auto">
          <li class="nav-item"></li>
        </ul>
        <span class="navbar-text">
          {{ gitHash }} · &copy; 2020 {{ config.name }}, All rights
          reserved</span
        >
      </div>
    </nav>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import halfmoon from "halfmoon";
import Navbar from "@/components/Navbar.vue";
import { mapState } from "vuex";

export default Vue.extend({
  components: {
    Navbar
  },
  data() {
    return {
      gitHash: process.env.VUE_APP_GIT_HASH
    };
  },
  computed: mapState({ user: "user", config: "config" }),
  async mounted() {
    await this.$store.dispatch("fetchConfig");
    halfmoon.onDOMContentLoaded();
    await this.$store.dispatch("fetchUser");
    const returnURL = localStorage.getItem("returnURL");
    if (this.user.id && returnURL) {
      localStorage.removeItem("returnURL");
      window.location.href = returnURL;
    }
  }
});
</script>
