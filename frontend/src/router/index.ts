import { createRouter, createWebHashHistory, RouteRecordRaw } from 'vue-router';
import Home from '../views/Home.vue'
import Webscan from '../views/Permeation/Webscan.vue'
import Portscan from '../views/Permeation/Portscan.vue'
import Dirsearch from '../views/Permeation/Dirsearch.vue'
import Fofa from '../views/SpaceEngine/Fofa.vue'
import Hunter from '../views/SpaceEngine/Hunter.vue'
import Setting from '../views/Settings.vue'
// 懒加载
export const routes: Array<RouteRecordRaw> = [
  { path: '/', component: Home },
  { path: '/Permeation/Webscan', component: Webscan },
  { path: '/Permeation/Portscan', component: Portscan },
  { path: '/Permeation/Dirsearch', component: Dirsearch },
  { path: '/Permeation/Pocdetail', component: () => import('../views/Permeation/Pocdetail.vue') },
  { path: '/Asset/Asset', component: () => import('../views/Asset/Asset.vue') },
  { path: '/Asset/Subdomain', component: () => import('../views/Asset/Subdomain.vue') },
  { path: '/Asset/Ipdomain', component: () => import('../views/Asset/Ipdomain.vue') },
  { path: '/SpaceEngine/Fofa', component: Fofa },
  { path: '/SpaceEngine/Hunter', component: Hunter },
  { path: '/SpaceEngine/AgentPool', component: () => import('../views/SpaceEngine/AgentPool.vue') },
  { path: '/Tools/Codec', component: () => import('../views/Tools/Codec.vue') },
  { path: '/Tools/System', component: () => import('../views/Tools/System.vue') },
  { path: '/Tools/Fscan', component: () => import('../views/Tools/Fscan.vue') },
  { path: '/Tools/Reverse', component: () => import('../views/Tools/Reverse.vue') },
  { path: '/Tools/Thinkdict', component: () => import('../views/Tools/Thinkdict.vue') },
  { path: '/Tools/Wxappid', component: () => import('../views/Tools/Wxappid.vue') },
  { path: '/Settings', component: Setting },
];

const router = createRouter({
    history: createWebHashHistory(),
    routes
})

export default router;