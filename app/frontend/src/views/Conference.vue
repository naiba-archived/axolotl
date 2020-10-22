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
        v-if="localStream"
        :muted="true"
        :nickname="user.nickname"
        :stream="localStream"
        :offset="0"
      />
      <div v-for="(item, index) in peersArray" v-bind:key="index">
        <Hello
          v-if="item.peer._remoteStreams.length"
          :stream="item.peer._remoteStreams[0]"
          :nickname="item.name"
          :offset="(index+1) * 5"
        />
      </div>
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
import { MonacoBinding } from "y-monaco";
import { executeCode, fetchCodeList } from "../api/code";
import Clipboard from "clipboard";
import halfmoon from "halfmoon";
import YWS from "../utils/y-ws";

export default Vue.extend({
  name: "Conference",
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
      peers: {} as any,
      localStream: undefined as any
    };
  },
  computed: {
    ...mapState({
      darkMode: "darkMode",
      user: "user"
    }),
    peersArray() {
      var arr = [] as any;
      Object.keys(this.peers).forEach((k: any) => {
        arr.push({
          name: k,
          peer: this.peers[k]
        });
      });
      return arr;
    }
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
    this.localStream.getTracks().forEach((track: any) => track.stop());
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
        language: "go"
      }
    );

    // init websocket
    this.ws = new WebSocket(
      (window.location.protocol == "http:" ? "ws" : "wss") +
        "://" +
        window.location.host +
        "/ws/" +
        this.$router.currentRoute.params.id
    );
    const ydocument = new Y.Doc();
    const type = ydocument.getText("monaco");
    const yws = new YWS(ydocument, (data: Uint8Array) => {
      if (this.ws) this.ws.send(data);
    });
    const binding = new MonacoBinding(
      type,
      this.editor.getModel(),
      new Set([this.editor]),
      yws.awareness
    );
    this.ws.onopen = async (e: any) => {
      yws.onOpen();
      this.localStream = await navigator.mediaDevices.getUserMedia({
        video: {
          width: 160,
          height: 90
        },
        audio: true
      });
      console.log("local stream", this.localStream);
    };
    this.ws.onclose = (e: any) => {
      this.ws = null;
    };
    this.ws.onmessage = (e: MessageEvent) => {
      if (e.data instanceof Blob) {
        yws.onMessage(e.data);
        return;
      }
      const data = JSON.parse(e.data);
      let signal;
      switch (data.type) {
        case 0:
          signal = JSON.parse(data.data);
          if (this.peers[data.from]) {
            this.peers[data.from].signal(signal);
          }
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

        case 3:
          if (this.peers[data.from]) {
            this.peers[data.from].streams.forEach((stream: any) => {
              stream.getTracks().forEach((track: any) => track.stop());
            });
            this.peers[data.from].destroy();
          }
          const peer = new Peer({
            stream: this.localStream
          });
          peer.on("connect", () => {
            console.log("passive peer connected", data.from, this.peers);
          });
          peer.on("error", (err: any) => {
            Vue.delete(this.peers, data.from);
            console.log(
              "passive peer disconnected",
              data.from,
              err,
              this.peers
            );
          });
          peer.on("stream", (stream: any) => {
            console.log("passive peer stream", data.from, stream, this.peers);
          });
          peer.on("signal", (signal: any) => {
            this.ws.send(
              JSON.stringify({
                type: 0,
                to: data.from,
                data: JSON.stringify(signal)
              })
            );
          });
          Vue.set(this.peers, data.from, peer);
          console.log("create passive peer", this.peers);
          break;

        case 4:
          // 设置所选语言
          this.lang = data.data.lang;
          monaco.editor.setModelLanguage(
            this.editor.getModel(),
            this.lang.split(":")[0]
          );
          // 创建各个 Peer
          if (data.data.user) {
            for (let i = 0; i < data.data.user.length; i++) {
              const fromUser = data.data.user[i];
              const peer = new Peer({
                initiator: true
              });
              peer.on("connect", (conn: any) => {
                peer.addStream(this.localStream);
                console.log("active peer connect", fromUser, conn, peer);
                Vue.set(this.peers, fromUser, peer);
              });
              peer.on("stream", (stream: any) => {
                console.log("active peer stream", fromUser, stream, this.peers);
              });
              peer.on("error", (err: any) => {
                Vue.delete(this.peers, fromUser);
                console.log(
                  "active peer disconnected",
                  fromUser,
                  err,
                  this.peers
                );
              });
              peer.on("signal", (signal: any) => {
                this.ws.send(
                  JSON.stringify({
                    type: 0,
                    to: fromUser,
                    data: JSON.stringify(signal)
                  })
                );
              });
              Vue.set(this.peers, fromUser, peer);
              console.log("create active peer", this.peers);
            }
          }
          break;

        default:
          break;
      }
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