<template>
  <div>
    <div id="vditor"></div>
    <Hello v-for="(peer,index) in peers" :stream="peer.streams[0]" v-bind:key="index" />
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import Vditor from "vditor";
import "../../node_modules/vditor/src/assets/scss/index.scss";
import { mapState } from "vuex";
import vditorConfig from "../utils/vditor";
import Hello from "@/components/Hello.vue";
import Peer from "simple-peer";

export default Vue.extend({
  name: "Meeting",
  components: {
    Hello
  },
  data() {
    return {
      vditor: {} as any,
      peers: [],
    };
  },
  computed: {
    ...mapState({
      darkMode: "darkMode"
    })
  },
  watch: {
    darkMode() {
      this.vditor.setTheme(this.darkMode ? "dark" : "classic");
    }
  },
  mounted() {
    // init vditor
    vditorConfig.theme = this.darkMode ? "dark" : "classic";
    this.vditor = new Vditor("vditor", vditorConfig);

    // init websocket
    const ws = new WebSocket("ws://localhost/ws/1");
    ws.onopen = async (e: any) => {
      console.log("onopen", e);
      const stream = await navigator.mediaDevices.getUserMedia({
        video: {
          width: 160,
          height: 90
        },
        audio: true
      });
      const peer = new Peer({
        initiator: true,
        trickle: false,
        stream: stream
      });
      ws.send(JSON.stringify({ type: 0, data: JSON.stringify(this.peer1) }));
    };
    ws.onmessage = function(e: any) {
      console.log("onmessage", e);
    };

    // on close warning
    window.onbeforeunload = function(e: any) {
      var e = e || window.event;
      if (e) {
        e.returnValue = "ATTENTION REQUIRED";
      }
      return "ATTENTION REQUIRED";
    };
  }
});
</script>
<style lang="scss">
#vditor .vditor-reset {
  color: unset;
}
</style>