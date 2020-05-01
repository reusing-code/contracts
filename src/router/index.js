import Vue from "vue";
import VueRouter from "vue-router";

import About from "../views/About.vue";
import Home from "../views/Home.vue";
import Contracts from "../views/Contracts.vue";
import EditTest from "../components/EditTest.vue";

Vue.use(VueRouter);

const routes = [
  { path: "/", name: "Home", component: Home },
  {
    path: "/about",
    name: "About",
    component: About,
  },
  { path: "/contracts", name: "Contracts", component: Contracts },
  { path: "/test", name: "Test", component: EditTest },
];

const router = new VueRouter({
  mode: "history",
  base: process.env.BASE_URL,
  routes,
});

export default router;
