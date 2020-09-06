<template>
  <nav class="navbar">
    <div class="container">
      <a href="#" class="navbar-brand">
        <img src="/favicon.png" alt="Hello Engineer" />
        Hello Engineer
      </a>
      <ul class="navbar-nav d-none d-flex">
        <li class="nav-item active">
          <a href="#" class="nav-link">Docs</a>
        </li>
        <li class="nav-item">
          <a href="#" class="nav-link">Products</a>
        </li>
      </ul>
      <div class="d-none d-flex ml-auto">
        <button v-if="!user.id" @click="login" class="btn btn-primary" type="submit">
          <font-awesome-icon :icon="['fab', 'github']" />&nbsp;Sign in
        </button>
        <div v-if="user.id" class="dropdown toggle-on-hover">
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
            <a href="#" @click="logout" class="dropdown-item">Logout</a>
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
      user: "user"
    })
  },
  methods: {
    login() {
      window.location.href = "/api/oauth2/login";
    },
    logout() {
      this.$store.dispatch("logout");
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