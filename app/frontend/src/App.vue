<template>
  <div class="page-wrapper with-navbar">
    <Navbar />
    <div class="content-wrapper">
      <div class="container">
        <router-view />
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import halfmoon from "halfmoon";
import Navbar from "@/components/Navbar.vue";

export default Vue.extend({
  components: {
    Navbar
  },
  mounted() {
    halfmoon.onDOMContentLoaded();
    document.title = "Hello Engineer";
  }
});
</script>