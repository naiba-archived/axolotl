import Vue from "vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";

// font-awesome
import { library } from '@fortawesome/fontawesome-svg-core'
import { faLanguage, faAngleDown, faMoon, faSun } from '@fortawesome/free-solid-svg-icons'
import { faGithub } from '@fortawesome/free-brands-svg-icons'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
library.add(faLanguage, faAngleDown, faGithub, faMoon, faSun)
Vue.component('font-awesome-icon', FontAwesomeIcon)

// halfmoon
require("halfmoon/css/halfmoon.min.css");

Vue.config.productionTip = false;

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount("#app");
