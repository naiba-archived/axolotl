<template>
  <div
    :style="'top: ' + (8 + offset) + 'rem;left: ' + (8 + offset) + 'rem;'"
    v-drag
  >
    <p><font-awesome-icon :icon="['fab', 'github']" /> {{ nickname }}</p>
    <video ref="cam" :muted="muted" autoplay></video>
  </div>
</template>

<script lang="ts">
import Vue from "vue";

export default Vue.extend({
  name: "Hello",
  props: {
    nickname: String,
    stream: MediaStream,
    muted: Boolean,
    offset: Number
  },
  mounted() {
    (this.$refs.cam as any).srcObject = this.stream;
  }
});
</script>
<style lang="scss" scoped>
div {
  position: absolute;
  z-index: 9999;

  > video {
    background-color: #f5f5f5;
    width: 320px;
    height: 180px;
    border-radius: 1rem;
    cursor: move;
  }

  > p {
    border: 1px solid black;
    padding-left: 3px;
    padding-right: 3px;
    color: black;
    font-size: 15px;
    font-weight: bolder;
    position: absolute;
    bottom: -0.3rem;
    left: 1.1rem;
  }
}
</style>
