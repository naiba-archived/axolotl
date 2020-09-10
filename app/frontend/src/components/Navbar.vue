<template>
  <nav class="navbar">
    <div class="container">
      <a href="javascript:void(0)" @click="to('/')" class="navbar-brand">
        <img src="/favicon.png" alt="Hello Engineer" />
        Hello Engineer
      </a>
      <ul class="navbar-nav d-none d-flex">
        <li class="nav-item">
          <a
            href="javascript:void(0)"
            @click="wOpen('https://github.com/naiba/helloengineer')"
            class="nav-link"
          >Source Code</a>
        </li>
      </ul>
      <div class="d-none d-flex ml-auto">
        <button class="btn btn-action mr-5" type="button" @click="toggleDarkMode">
          <font-awesome-icon v-if="!darkMode" :icon="['fa', 'moon']" />
          <font-awesome-icon v-if="darkMode" :icon="['fa', 'sun']" />
        </button>
        <button v-if="!user.id" @click="login" class="btn btn-primary" type="submit">
          <font-awesome-icon :icon="['fab', 'github']" />&nbsp;Sign in
        </button>
        <div v-if="user.id" class="dropdown with-arrow toggle-on-hover">
          <button
            class="btn logged-in-btn"
            data-toggle="dropdown"
            type="button"
            id="nav-logged-in"
            aria-haspopup="true"
            aria-expanded="false"
          >
            <font-awesome-icon :icon="['fab', 'github']" />
            <p>{{user.nickname}}</p>
            <font-awesome-icon :icon="['fa', 'angle-down']" class="ml-5" aria-hidden="true" />
          </button>
          <div class="dropdown-menu">
            <a href="javascript:void(0)" @click="logout" class="dropdown-item">Logout</a>
          </div>
        </div>
      </div>
    </div>
  </nav>
</template>

<script lang="ts">
import Vue from "vue";
import { mapState } from "vuex";

export default Vue.extend({
  name: "Navbar",
  computed: {
    ...mapState({
      user: "user",
      darkMode: "darkMode"
    })
  },
  methods: {
    wOpen(url: string) {
      window.open(url);
    },
    login() {
      window.location.href = "/api/oauth2/login";
    },
    logout() {
      this.$store.dispatch("logout");
    },
    to(path: any) {
      if (this.$router.currentRoute.path !== path) {
        this.$router.push(path);
      }
    },
    toggleDarkMode() {
      this.$store.dispatch("toggleDarkMode");
    }
  }
});
</script>

<style lang="scss" scoped>
.navbar {
  padding-left: unset;
  padding-right: unset;
  .container {
    padding-left: unset;
    padding-right: unset;
  }
  .logged-in-btn {
    display: flex;
    justify-content: center;
    align-items: center;
    > p {
      margin-left: 0.4rem;
      margin-right: 0.1rem;
    }
  }
}
</style>