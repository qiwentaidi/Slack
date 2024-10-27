import { createApp, App } from "vue";
import AppComponent from "./App.vue";
import "./style/style.css";
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
//引入依赖和语言
import hljs from "highlight.js/lib/core";
import hljsVuePlugin from "@highlightjs/vue-plugin";
//按需引入语言
import bash from "highlight.js/lib/languages/bash";
import yaml from "highlight.js/lib/languages/yaml";
import http from "highlight.js/lib/languages/http";
import '@imengyu/vue3-context-menu/lib/vue3-context-menu.css'
import ContextMenu from '@imengyu/vue3-context-menu'

hljs.registerLanguage("bash", bash);
hljs.registerLanguage("yaml", yaml);
hljs.registerLanguage("http", http);

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

// 使得comonent 可以正确渲染el-icon
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

app.use(router).use(i18n).use(hljsVuePlugin).use(ContextMenu).mount("#app");
