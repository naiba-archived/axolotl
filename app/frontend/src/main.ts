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

// drag element
Vue.directive('drag', {
  bind: function (el) {
    let odiv = el;   //获取当前元素
    const minSpace = 50;
    odiv.onmousedown = (e) => {
      //算出鼠标相对元素的位置
      let disX = e.clientX - odiv.offsetLeft;
      let disY = e.clientY - odiv.offsetTop;
      document.onmousemove = (e) => {
        //用鼠标的位置减去鼠标相对元素的位置，得到元素的位置
        let left = e.clientX - disX;
        let top = e.clientY - disY;

        if (top < 0 || left < 0 || top > window.innerHeight - minSpace || left > window.innerWidth - minSpace) return;

        //移动当前元素
        odiv.style.left = left + 'px';
        odiv.style.top = top + 'px';
      };
      document.onmouseup = (e) => {
        document.onmousemove = null;
        document.onmouseup = null;
      };
    };
  }
})

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount("#app");
