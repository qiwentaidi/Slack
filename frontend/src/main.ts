import { createApp, App } from "vue";
import AppComponent from "./App.vue";
import "./style.css";
import router from "./router";
import i18n from './i18n/index' //引入配置的语言
import { ElMessage, ElMessageBox, ElNotification } from "element-plus";
import 'element-plus/theme-chalk/dark/css-vars.css'
import "element-plus/theme-chalk/el-message.css";
import "element-plus/theme-chalk/el-message-box.css";
import "element-plus/theme-chalk/el-notification.css";
import 'element-plus/theme-chalk/el-loading.css'
import * as ElementPlusIconsVue from "@element-plus/icons-vue";
import global from "./global";
import "./style/dark.css"
import "./style/light.css"
import { directive } from 'vue3-menus';


let theme = localStorage.getItem('theme') || "light"

global.Theme.value  = theme == "dark" ? true : false

export default (app: App<Element>) => {
  // 全局配置
  app.config.globalProperties.$ELEMENT = {};
  app.use(ElMessage);
  app.use(ElMessageBox);
  app.use(ElNotification);
};

const app = createApp(AppComponent)
app.directive('menus', directive); // 注册指令

// 使得comonent 可以正确渲染el-icon
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

app.use(router).use(i18n).mount("#app");
