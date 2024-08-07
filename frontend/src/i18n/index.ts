import { createI18n } from "vue-i18n";
import en from './en_US'
import zh from "./zh_CN"

// 默认读取本地存储语言设置
const defaultLocale = localStorage.getItem('language') || 'zh'

const i18n = createI18n({
  locale: defaultLocale,// 默认语言
  fallbackLocale: 'en',// 不存在默认则为英文
  allowComposition: true,// 允许组合式api
  globalInjection: true, // 全局注册$t方法
  silentTranslationWarn: true, // 去掉警告
  messages: {
    en, // 标识:配置对象
    zh
  },
})
export default i18n