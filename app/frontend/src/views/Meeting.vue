<template>
  <div class="container">
    <div class="row">
      <div class="col-6">
        <div id="editor"></div>
      </div>
      <div class="col-6">
        <textarea class="form-control" disabled readonly></textarea>
      </div>
    </div>
    <Hello
      v-if="selfPeer.streams"
      :muted="true"
      :stream="selfPeer.streams[0]"
    />
    <Hello
      v-for="(peer, index) in peers"
      :stream="peer.streams[0]"
      v-bind:key="index"
    />
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import { mapState } from "vuex";
import Hello from "@/components/Hello.vue";
import Peer from "simple-peer";
import * as Y from "yjs";
import { WebrtcProvider } from "y-webrtc";
import * as monaco from "monaco-editor";
import { MonacoBinding } from "y-monaco";

export default Vue.extend({
  name: "Meeting",
  components: {
    Hello
  },
  data() {
    return {
      editor: {} as any,
      peers: [] as any,
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
      monaco.editor.setTheme(this.darkMode ? "vs-dark" : "vs");
    }
  },
  mounted() {
    const ydocument = new Y.Doc();
    const provider = new WebrtcProvider(
      this.$router.currentRoute.params.id,
      ydocument,
      {
        password: "optional-room-password"
      } as any
    );
    const type = ydocument.getText("monaco");
    this.editor = monaco.editor.create(
      document.getElementById("editor") || new HTMLElement(),
      {
        value: "",
        fontSize: 18,
        theme: this.darkMode ? "vs-dark" : "vs",
        language: "php"
      }
    );
    const monacoBinding = new MonacoBinding(
      type,
      this.editor.getModel(),
      new Set([this.editor]),
      provider.awareness
    );

    // init websocket
    const ws = new WebSocket(
      (window.location.protocol == "http:" ? "ws" : "wss") +
        "://" +
        window.location.host +
        "/ws/1"
    );
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
        const peer = new Peer({
          initiator: true,
          stream: data
        });
        peer.on("close", () => {
          console.log("onClose", peer);
        });
        peer.on("error", (err: any) => {
          console.log("error", peer, err);
        });
        this.peers.push(peer);
      });
    };
    ws.onmessage = (e: any) => {
      const data = JSON.parse(e.data);
      let signal;
      switch (data.type) {
        case 0:
          signal = JSON.parse(data.data);
          console.log("onSignal", signal);
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
<style lang="scss" scoped>
#editor {
  height: calc(90vh);
}
textarea {
  height: 100%;
  resize: none;
  overflow-y: scroll;
}
</style>