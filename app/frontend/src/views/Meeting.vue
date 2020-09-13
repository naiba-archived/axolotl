<template>
  <div>
    <div id="vditor"></div>
    <Hello v-if="selfPeer.streams" :stream="selfPeer.streams[0]" />
    <Hello v-for="(peer, index) in peers" :stream="peer.streams[0]" v-bind:key="index" />
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
      selfPeer: {} as any
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
    const ws = new WebSocket((window.location.protocol == 'http:'?'ws':'wss')+"://" +window.location.host+"/ws/1");
    console.log(ws, navigator.mediaDevices.getUserMedia);
    ws.onopen = async (e: any) => {
      const stream = await navigator.mediaDevices.getUserMedia({
        video: {
          width: 160,
          height: 90
        },
        audio: true
      });
      console.log("stream", stream);
      this.selfPeer = new Peer({
        initiator: true,
        stream: stream
      });
      console.log("selfPeer", this.selfPeer);
      this.selfPeer.on("signal", (data: any) => {
        ws.send(JSON.stringify({ type: 0, data: JSON.stringify(data) }));
      });
      this.selfPeer.on("connect", (data: any) => {
        console.log("onConnect", data);
      });
      this.selfPeer.on("stream", (data: any) => {
        console.log("onStream", data);
      });
    };
    ws.onmessage = (e: any) => {
      const data = JSON.parse(e.data);
      let signal;
      switch (data.type) {
        case 0:
          signal = JSON.parse(data.data);
          console.log("remote signal", signal);
          this.selfPeer.signal(signal);
          break;

        default:
          break;
      }
    };

    // on close warning
    window.onbeforeunload = function(e: any) {
      const ee = e || window.event;
      if (ee) {
        ee.returnValue = "ATTENTION REQUIRED";
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
