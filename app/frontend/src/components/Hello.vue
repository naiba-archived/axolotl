<template>
  <div v-drag>
    <video ref="cam" :muted="muted" autoplay></video>
    <p>X</p>
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
  > video {
    background-color: whitesmoke;
    width: 320px;
    height: 180px;
    z-index: 9999;
    border-radius: 1rem;
    cursor: move;
  }
  > p {
    color: red;
    font-size: 30px;
    font-weight: bolder;
    position: absolute;
    bottom: -2.5rem;
    left: 1rem;
  }
}
</style>
