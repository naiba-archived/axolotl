<template>
  <div>
    <Navbar />
    <router-view />
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
  async mounted() {
    halfmoon.onDOMContentLoaded();
    await this.$store.dispatch("fetchUser");
  }
});
</script>