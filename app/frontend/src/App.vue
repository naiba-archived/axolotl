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
import { mapState } from "vuex";

export default Vue.extend({
  components: {
    Navbar
  },
  computed: mapState({ user: "user" }),
  async mounted() {
    halfmoon.onDOMContentLoaded();
    await this.$store.dispatch("fetchUser");
  }
});
</script>