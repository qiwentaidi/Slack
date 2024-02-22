import { createRouter, createWebHashHistory, RouteRecordRaw } from 'vue-router';
// 懒加载
export const routes: Array<RouteRecordRaw> = [
  { path: '/', component: () => import('../views/Home.vue') },
  { path: '/Permeation/Webscan', component: () => import('../views/Permeation/Webscan.vue') },
  { path: '/Permeation/Portscan', component: () => import('../views/Permeation/Portscan.vue') },
  { path: '/Permeation/Dirsearch', component: () => import('../views/Permeation/Dirsearch.vue') },
  { path: '/Permeation/Pocdetail', component: () => import('../views/Permeation/Pocdetail.vue') },
  { path: '/Asset/Asset', component: () => import('../views/Asset/Asset.vue') },
  { path: '/Asset/Subdomain', component: () => import('../views/Asset/Subdomain.vue') },
  { path: '/Asset/Ipdomain', component: () => import('../views/Asset/Ipdomain.vue') },
  { path: '/SpaceEngine/Fofa', component: () => import('../views/SpaceEngine/Fofa.vue') },
  { path: '/SpaceEngine/Hunter', component: () => import('../views/SpaceEngine/Hunter.vue') },
  { path: '/SpaceEngine/AgentPool', component: () => import('../views/SpaceEngine/AgentPool.vue') },
  { path: '/Tools/Codec', component: () => import('../views/Tools/Codec.vue') },
  { path: '/Tools/System', component: () => import('../views/Tools/System.vue') },
  { path: '/Tools/Fscan', component: () => import('../views/Tools/Fscan.vue') },
  { path: '/Tools/Reverse', component: () => import('../views/Tools/Reverse.vue') },
  { path: '/Tools/Thinkdict', component: () => import('../views/Tools/Thinkdict.vue') },
  { path: '/Tools/Wxappid', component: () => import('../views/Tools/Wxappid.vue') },
  { path: '/Settings', component: () => import('../views/Settings.vue') },
];

const router = createRouter({
    history: createWebHashHistory(),
    routes
})

export default router;