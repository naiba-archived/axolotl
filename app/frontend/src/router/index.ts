import { getConfig } from '@/api/config';
import Vue from "vue";
import VueRouter, { RouteConfig } from "vue-router";
import Home from "../views/Home.vue";

Vue.use(VueRouter);

const routes: Array<RouteConfig> = [
  {
    path: "/",
    name: "Home",
    meta: {
      // title: "Home"
    },
    component: Home
  },
  {
    path: "/conference/:id",
    name: "Conference",
    meta: {
      title: "Conference"
    },
    component: () => import("../views/Conference.vue")
  }
];

const router = new VueRouter({
  mode: "history",
  base: process.env.BASE_URL,
  routes
});

router.beforeEach(async function (from, to, next) {
  const config = await getConfig()
  document.title = (from.meta.title ? from.meta.title + " | " : "") + config.name;
  if (!from.meta.title) {
    document.title += " - " + config.desc
  }
  next();
});

export default router;
