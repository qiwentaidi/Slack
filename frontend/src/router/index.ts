import { createRouter, createWebHashHistory, RouteRecordRaw } from "vue-router";
// 懒加载
export const routes: Array<RouteRecordRaw> = [
  { 
    path: "/", 
    component: () => import("../views/Home.vue"),
    meta: {
      title: "Home",
    },
  },
  {
    path: "/Permeation/Webscan",
    component: () => import("../views/Permeation/Webscan.vue"),
    meta: {
      title: "Webscan",
    },
  },
  {
    path: "/Permeation/Portscan",
    component: () => import("../views/Permeation/Portscan.vue"),
    meta: {
      title: "Portscan",
    },
  },
  {
    path: "/Permeation/Dirsearch",
    component: () => import("../views/Permeation/Dirsearch.vue"),
    meta: {
      title: "Dirsearch",
    },
  },
  {
    path: "/Permeation/Jsfinder",
    component: () => import("../views/Permeation/Jsfinder.vue"),
    meta: {
      title: "Jsfinder",
    },
  },
  {
    path: "/Permeation/Exploitation",
    component: () => import("../views/Permeation/Exploitation.vue"),
    meta: {
      title: "Jsfinder",
    },
  },
  { path: "/Asset/Company", 
    component: () => import("../views/Asset/Company.vue"),
    meta: {
      title: "Asset",
    },
  },
  {
    path: "/Asset/Subdomain",
    component: () => import("../views/Asset/Subdomain.vue"),
    meta: {
      title: "Subdomain",
    },
  },
  {
    path: "/Asset/Ipdomain",
    component: () => import("../views/Asset/Ipdomain.vue"),
    meta: {
      title: "Ipdomain",
    },
  },
  {
    path: "/SpaceEngine/Fofa",
    component: () => import("../views/SpaceEngine/FOFA.vue"),
    meta: {
      title: "FOFA",
    },
  },
  {
    path: "/SpaceEngine/Hunter",
    component: () => import("../views/SpaceEngine/Hunter.vue"),
    meta: {
      title: "Hunter",
    },
  },
  {
    path: "/SpaceEngine/AgentPool",
    component: () => import("../views/SpaceEngine/AgentPool.vue"),
    meta: {
      title: "AgentPool",
    },
  },
  { path: "/Tools/Codec", 
    component: () => import("../views/Tools/Codec.vue"),
    meta: {
      title: "Codec",
    }
  },
  {
    path: "/Tools/System",
    component: () => import("../views/Tools/System.vue"),
    meta: {
      title: "System",
    },
  },
  { 
    path: "/Tools/DataHanding", 
    component: () => import("../views/Tools/DataHanding.vue"),
    meta: {
      title: "Fscan",
    }
  },
  { 
    path: "/Tools/Memo", 
    component: () => import("../views/Tools/Memo.vue"),
    meta: {
      title: "Memo",
    }
  },
  {
    path: "/Tools/Thinkdict",
    component: () => import("../views/Tools/Thinkdict.vue"),
    meta: {
      title: "Thinkdict",
    },
  },
  { 
    path: "/Tools/AKSK", 
    component: () => import("../views/Tools/AKSK.vue"),
    meta: {
      title: "AKSK",
    }
  },
  { 
    path: "/Settings", 
    component: () => import("../views/Settings.vue"),
    meta: {
      title: "Settings",
    }
  },
];

const router = createRouter({
  history: createWebHashHistory(),
  routes,
});

export default router;
