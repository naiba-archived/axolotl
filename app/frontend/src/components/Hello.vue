<template>
  <div v-drag>
    <video ref="cam" :muted="muted" autoplay></video>
  </div>
</template>

<script lang="ts">
import Vue from "vue";

export default Vue.extend({
  name: "Hello",
  props: {
    stream: MediaStream,
    muted: Boolean
  },
  mounted() {
    (this.$refs.cam as any).srcObject = this.stream;
  }
});
</script>
<style lang="scss" scoped>
div {
  position: absolute;
  top: 8rem;
  left: 8rem;
  z-index: 9999;

  > video {
    background-color: whitesmoke;
    width: 320px;
    height: 180px;
    border-radius: 1rem;
    cursor: move;
  }
}
</style>
