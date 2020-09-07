<template>
  <div>
    <div id="vditor"></div>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import Vditor from "vditor";
import "../../node_modules/vditor/src/assets/scss/index.scss";
import { mapState } from "vuex";

export default Vue.extend({
  name: "Meeting",
  data() {
    return {
      vditor: {} as any
    };
  },
  computed: {
    ...mapState({
      darkMode: "darkMode"
    })
  },
  watch: {
    darkMode() {
      console.log("darkMode", this.darkMode);
      this.vditor.setTheme(this.darkMode ? "dark" : "classic");
    }
  },
  mounted() {
    this.vditor = new Vditor("vditor", {
      height: window.innerHeight - 50,
      mode: "wysiwyg",
      theme: this.darkMode ? "dark" : "classic",
      preview: {
        hljs: {
          lineNumber: true,
          style: "dracula"
        }
      },
      toolbar: [
        "emoji",
        "headings",
        "bold",
        "italic",
        "strike",
        "link",
        "|",
        "list",
        "ordered-list",
        "check",
        "outdent",
        "indent",
        "|",
        "quote",
        "line",
        "code",
        "inline-code",
        "insert-before",
        "insert-after",
        "|",
        "upload",
        "record",
        "table",
        "|",
        "undo",
        "redo",
        "|",
        "fullscreen",
        // "edit-mode",
        {
          name: "more",
          toolbar: [
            "both",
            "code-theme",
            "content-theme",
            "export",
            "outline",
            "preview"
            // "devtools",
            // "info",
            // "help"
          ]
        }
      ]
    });
  }
});
</script>
<style lang="scss">
#vditor .vditor-reset {
  color: unset;
}
</style>