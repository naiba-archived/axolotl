<template>
  <div class="container">
    <div class="row">
      <div class="col-12">
        <div class="code-runner row justify-content-between">
          <div class="col-4">
            <select v-model="lang" @change="chooseLang" class="form-control">
              <option value="" selected="selected" disabled="disabled">
                Select Programming Language
              </option>
              <option v-for="(k, v) in langs" :key="v" :value="v">
                {{ v }}
              </option>
            </select>
          </div>
          <div class="col-4">
            <button @click="execute" class="btn" type="button">Execute</button>
          </div>
        </div>
      </div>
      <div class="col-6">
        <div id="editor"></div>
      </div>
      <div class="col-6">
        <textarea v-model="log" class="form-control" disabled readonly>
        </textarea>
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
import * as monaco from "monaco-editor";
import * as Y from "yjs";
import { WebrtcProvider } from "y-webrtc";
import { MonacoBinding } from "y-monaco";
import { executeCode, fetchCodeList } from "../api/code";

export default Vue.extend({
  name: "Meeting",
  components: {
    Hello
  },
  data() {
    return {
      executing: false,
      ws: {} as any,
      lang: "",
      log: "Waiting for execution",
      langs: {} as any,
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
  methods: {
    async execute() {
      if (this.executing) {
        return;
      }
      try {
        this.executing = true;
        let out = "";
        await executeCode(
          this.$router.currentRoute.params.id,
          this.editor.getValue(),
          this.lang
        );
      } catch (error) {
        console.log("run code", error);
      } finally {
        this.executing = false;
      }
    },
    chooseLang() {
      monaco.editor.setModelLanguage(
        this.editor.getModel(),
        this.lang.split(":")[0]
      );
      if (this.ws) {
        this.ws.send(JSON.stringify({ type: 1, data: this.lang }));
      }
      if (this.editor.getValue().trim() == "") {
        this.editor.setValue(this.langs[this.lang].template);
      }
    }
  },
  async mounted() {
    this.langs = await fetchCodeList();

    this.editor = monaco.editor.create(
      document.getElementById("editor") || new HTMLElement(),
      {
        value: "",
        fontSize: 18,
        theme: this.darkMode ? "vs-dark" : "vs",
        language: "php"
      }
    );
    const ydocument = new Y.Doc();
    const provider = new WebrtcProvider(
      this.$router.currentRoute.params.id,
      ydocument,
      {
        password: "optional-room-password"
      } as any
    );
    const type = ydocument.getText("monaco");
    new MonacoBinding(
      type,
      this.editor.getModel(),
      new Set([this.editor]),
      provider.awareness
    );

    // init websocket
    this.ws = new WebSocket(
      (window.location.protocol == "http:" ? "ws" : "wss") +
        "://" +
        window.location.host +
        "/ws/" +
        this.$router.currentRoute.params.id
    );
    this.ws.onopen = async (e: any) => {
      return;
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
        this.ws.send(JSON.stringify({ type: 0, data: JSON.stringify(data) }));
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
    this.ws.onclose = (e: any) => {
      console.log("onclose", e);
      this.ws = null;
    };
    this.ws.onmessage = (e: any) => {
      const data = JSON.parse(e.data);
      let signal;
      switch (data.type) {
        case 0:
          signal = JSON.parse(data.data);
          console.log("onSignal", signal);
          this.selfPeer.signal(signal);
          break;

        case 1:
          this.lang = data.data;
          monaco.editor.setModelLanguage(
            this.editor.getModel(),
            this.lang.split(":")[0]
          );
          break;

        case 2:
          let out = "";
          if (data.data[0] == "{") {
            out = JSON.parse(data.data).out;
          } else {
            out = data.data;
          }
          this.log = out + "\n" + this.log;
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
  height: calc(85vh);
}
textarea {
  height: 100%;
  resize: none;
  overflow-y: scroll;
}
.code-runner {
  padding-top: 1rem;
  padding-bottom: 1rem;
  .col-4:last-child {
    text-align: right;
  }
}
</style>