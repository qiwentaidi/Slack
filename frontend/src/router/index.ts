import { createRouter, createWebHashHistory, RouteRecordRaw } from "vue-router";
import Webscan from "../views/Permeation/Webscan.vue";
import Portcan from "../views/Permeation/Portscan.vue";
import Crack from "../views/Permeation/Crack.vue";

// 懒加载
export const routes: Array<RouteRecordRaw> = [
  { 
    path: "/", 
    component: () => import("../views/Home.vue"),
  },
  {
    path: "/Permeation/Webscan",
    component: Webscan,
  },
  {
    path: "/Permeation/Portscan",
    component: Portcan,
  },
  {
    path: "/Permeation/Crack",
    component: Crack,
  },
  {
    path: "/Permeation/Dirsearch",
    component: () => import("../views/Permeation/Dirsearch.vue"),
  },
  {
    path: "/Permeation/Jsfinder",
    component: () => import("../views/Permeation/Jsfinder.vue"),
  },
  {
    path: "/Permeation/Exploitation",
    component: () => import("../views/Permeation/Exploitation.vue"),
  },
  { path: "/Asset/Company", 
    component: () => import("../views/Asset/Company.vue"),
  },
  {
    path: "/Asset/Subdomain",
    component: () => import("../views/Asset/Subdomain.vue"),
  },
  {
    path: "/Asset/Ipdomain",
    component: () => import("../views/Asset/Ipdomain.vue"),
  },
  {
    path: "/SpaceEngine/Fofa",
    component: () => import("../views/SpaceEngine/FOFA.vue"),
  },
  {
    path: "/SpaceEngine/Hunter",
    component: () => import("../views/SpaceEngine/Hunter.vue"),
  },
  {
    path: "/SpaceEngine/AgentPool",
    component: () => import("../views/SpaceEngine/AgentPool.vue"),
  },
  { path: "/Tools/Codec", 
    component: () => import("../views/Tools/Codec.vue"),
  },
  {
    path: "/Tools/System",
    component: () => import("../views/Tools/System.vue"),
  },
  { 
    path: "/Tools/DataHanding", 
    component: () => import("../views/Tools/DataHanding.vue"),
  },
  { 
    path: "/Tools/Memo", 
    component: () => import("../views/Tools/Memo.vue"),
  },
  { 
    path: "/Tools/Reverse", 
    component: () => import("../views/Tools/Reverse.vue"),
  },
  {
    path: "/Tools/Thinkdict",
    component: () => import("../views/Tools/Thinkdict.vue"),
  },
  { 
    path: "/Tools/AKSK", 
    component: () => import("../views/Tools/AKSK.vue"),
  },
  { 
    path: "/Guard/Guarantee", 
    component: () => import("../views/Guard/Guarantee.vue"),
  },
  { 
    path: "/Settings", 
    component: () => import("../views/Settings.vue"),
  },
];

const router = createRouter({
  history: createWebHashHistory(),
  routes,
});

export default router;
