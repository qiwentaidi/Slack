import { createRouter, createWebHashHistory, RouteRecordRaw } from "vue-router";

// 在vite中使用 import.meta.glob 动态导入指定目录下的所有vue文件
let modules = import.meta.glob(["../views/**/*.vue"]);
const generateRoutes = () => {
  const routes: Array<RouteRecordRaw> = [];
  for (const path in modules) {
    let routePath = path
      .replace('../views', '')
      .replace('.vue', '');
    if (routePath === "/Tools/CyberChef") {
        continue
    }
    // 指定根路由
    routePath = routePath === "/Home" ? "/" : routePath;

    routes.push({
      path: routePath,
      component: modules[path],
    });
  }
  return routes;
};

const router = createRouter({
  history: createWebHashHistory(),
  routes: generateRoutes(),
});

export default router;
