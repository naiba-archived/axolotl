<template>
  <div class="container h-600">
    <div class="row h-full">
      <div class="col"></div>
      <div class="col align-self-center">
        <img width="500rem" src="axolotl.jpg" />
        <p class="font-smaller text-center">The curious Axolotl species of salamander. Public domain photo</p>
        <button
          @click="startConference"
          :disabled="!user.id"
          class="btn btn-lg btn-primary btn-block"
        >
          Start Conference
        </button>
      </div>
      <div class="col"></div>
    </div>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import { mapState } from "vuex";

export default Vue.extend({
  name: "Home",
  components: {},
  computed: {
    ...mapState({ user: "user" })
  },
  methods: {
    makeid(length: number) {
      var result = "";
      var characters =
        "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
      var charactersLength = characters.length;
      for (var i = 0; i < length; i++) {
        result += characters.charAt(
          Math.floor(Math.random() * charactersLength)
        );
      }
      return result;
    },
    startConference() {
      this.$router.push("/meeting/HE-" + this.makeid(8));
    }
  }
});
</script>
<style lang="scss" scoped>
</style>