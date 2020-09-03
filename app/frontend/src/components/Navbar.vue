<template>
  <nav class="navbar">
    <div class="container">
      <a href="#" class="navbar-brand">
        <img src="/favicon.png" alt="Hello Engineer" />
        Hello Engineer
      </a>
      <ul class="navbar-nav d-none d-md-flex">
        <li class="nav-item active">
          <a href="#" class="nav-link">Docs</a>
        </li>
        <li class="nav-item">
          <a href="#" class="nav-link">Products</a>
        </li>
      </ul>
      <form class="form-inline d-none d-md-flex ml-auto" action="..." method="...">
        <input type="text" class="form-control" placeholder="Email address" required="required" />
        <button class="btn btn-primary" type="submit">
          <font-awesome-icon :icon="['fab', 'github']" />&nbsp;Sign up
        </button>
      </form>
      <div class="navbar-content d-md-none ml-auto">
        <div class="dropdown with-arrow">
          <button
            class="btn"
            data-toggle="dropdown"
            type="button"
            id="navbar-dropdown-toggle-btn-1"
          >
            Menu
            <i class="fa fa-angle-down" aria-hidden="true"></i>
          </button>
          <div
            class="dropdown-menu dropdown-menu-right w-200"
            aria-labelledby="navbar-dropdown-toggle-btn-1"
          >
            <a href="#" class="dropdown-item">Docs</a>
            <a href="#" class="dropdown-item">Products</a>
            <div class="dropdown-divider"></div>
            <div class="dropdown-content">
              <form action="..." method="...">
                <div class="form-group">
                  <input
                    type="text"
                    class="form-control"
                    placeholder="Email address"
                    required="required"
                  />
                </div>
                <button class="btn btn-primary btn-block" type="submit">Sign up</button>
              </form>
            </div>
          </div>
        </div>
      </div>
    </div>
  </nav>
</template>

<script lang="ts">
import Vue from "vue";

export default Vue.extend({
  name: "Navbar"
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
}
</style>