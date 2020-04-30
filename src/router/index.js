import Vue from "vue";
import VueRouter from "vue-router";

import About from "../views/About.vue";
import Home from "../views/Home.vue";
import ContractList from "../components/ContractList.vue";

Vue.use(VueRouter);

const routes = [
  { path: "/", name: "Home", component: Home },
  {
    path: "/about",
    name: "About",
    component: About
  },
  { path: "/contracts", name: "Contracts", component: ContractList }
];

const router = new VueRouter({
  mode: "history",
  base: process.env.BASE_URL,
  routes
});

export default router;
