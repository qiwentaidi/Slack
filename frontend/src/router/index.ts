import { createRouter, createWebHashHistory, RouteRecordRaw } from "vue-router";

// 在vite中使用 import.meta.glob 动态导入指定目录下的所有vue文件
let modules = import.meta.glob(["../views/**/*.vue"]);

const generateRoutes = () => {
  const routes: Array<RouteRecordRaw> = [];
  for (const path in modules) {
    const routePath = path
      .replace('../views', '')
      .replace('.vue', '');
    routes.push({
      path: routePath,
      component: modules[path], // 默认懒加载
    });
  }
  return routes;
};

// 手动定义的根路径
const rootRoute: RouteRecordRaw = {
  path: "/",
  component: () => import("../views/Home.vue"),
};

// 合并路由
const routes: Array<RouteRecordRaw> = [rootRoute, ...generateRoutes()];

const router = createRouter({
  history: createWebHashHistory(),
  routes: routes,
});

export default router;
