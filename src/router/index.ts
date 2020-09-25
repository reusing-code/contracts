import Vue from 'vue';
import VueRouter, {RouteConfig} from 'vue-router';

import About from '../views/About.vue';
import Contracts from '../views/Contracts.vue';
import Home from '../views/Home.vue';

Vue.use(VueRouter);

const routes: Array<RouteConfig> = [
  {
    path: '/',
    name: 'Home',
    component: Home,
  },
  {
    path: '/about',
    name: 'About',
    component: About,
  },
  {path: '/contracts', name: 'Contracts', component: Contracts},
];

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes,
});

export default router;
