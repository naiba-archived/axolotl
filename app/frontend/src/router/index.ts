import Vue from "vue";
import VueRouter, { RouteConfig } from "vue-router";
import Home from "../views/Home.vue";
import { nextTick } from 'vue/types/umd';

Vue.use(VueRouter);

const routes: Array<RouteConfig> = [
  {
    path: "/",
    name: "Home",
    meta: {
      title: 'Home',
    },
    component: Home
  },
  {
    path: "/meeting/:id",
    name: "Meeting",
    meta: {
      title: "Meeting"
    },
    component: () =>
      import("../views/Meeting.vue")
  }
];

const router = new VueRouter({
  mode: "history",
  base: process.env.BASE_URL,
  routes
});

router.beforeEach(function (from, to, next) {
  document.title = (from.meta.title ? from.meta.title + ' | ' : '') + 'Hello Engineer'
  next()
})

export default router;
