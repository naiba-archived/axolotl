<template>
  <div class="page-wrapper">
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
              <button
                class="btn clipboard"
                :data-clipboard-text="conferenceLink"
              >
                Share Conference Link
              </button>
              &nbsp;
              <button @click="execute" class="btn">Execute</button>
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
        v-for="(stream, index) in streams"
        :stream="stream"
        v-bind:key="index"
      />
    </div>
    <div class="sticky-alerts"></div>
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
import Clipboard from "clipboard";
import halfmoon from "halfmoon";

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
      conferenceLink: "",
      log: "Waiting for execution",
      langs: {} as any,
      editor: {} as any,
      streams: [] as any,
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
  beforeDestroy() {
    this.selfPeer.streams.forEach((stream: any) => {
      stream.getTracks().forEach((track: any) => track.stop());
    });
    this.selfPeer.destroy();
  },
  async mounted() {
    halfmoon.onDOMContentLoaded();
    this.conferenceLink = window.location.href;
    const clip = new Clipboard("button.clipboard");
    clip.on("success", function(e: any) {
      halfmoon.initStickyAlert({
        content:
          "The conference link has been copied, please send it to the participants.",
        title: "Successful copied",
        alertType: "alert-success"
      });
      e.clearSelection();
    });

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
        password: window.location.host
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
      const stream = await navigator.mediaDevices.getUserMedia({
        video: {
          width: 160,
          height: 90
        },
        audio: true
      });
      this.selfPeer = new Peer({
        initiator: true,
        stream: stream
      });
      this.selfPeer.on("signal", (data: any) => {
        this.ws.send(JSON.stringify({ type: 0, data: JSON.stringify(data) }));
      });
      this.selfPeer.on("stream", (stream: any) => {
        stream.oninactive = () => {
          for (let i = 0; i < this.streams.length; i++) {
            if (this.streams[i].id == stream.id) {
              this.streams.splice(i,1);
              return;
            }
          }
        };
        this.streams.push(stream);
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

    // after everything
    try {
      this.langs = await fetchCodeList();
    } catch (error) {
      console.log(error);
    }
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