<script lang="ts" setup>
import { WindowToggleMaximise, WindowGetSize, WindowIsMaximised } from "wailsjs/runtime/runtime";
import global from "@/stores";
import { ref, onMounted, onBeforeUnmount, computed, nextTick } from "vue";
import runnerIcon from "@/assets/icon/apprunner.svg";
import { Search } from "@element-plus/icons-vue";
import { titlebarStyle } from "@/stores/style";
import { routerControl, windowsControl } from "@/stores/options";
import throttle from "lodash/throttle";
import { SaveWindowsScreenSize } from "wailsjs/go/services/Database";
import { useDark, useToggle } from "@vueuse/core";
import { useI18n } from "vue-i18n";
import { useRouter, useRoute } from "vue-router";
import MenuList from "@/router/menu";
import { IsMacOS } from "wailsjs/go/services/File";

const showLogger = ref(false);
const searchText = ref("");
const searchRef = ref();

const syncMaxState = async () => {
  global.temp.isMax = await WindowIsMaximised();
};

onMounted(() => {
  IsMacOS().then(res => {
    global.temp.isMacOS = res
  })
  syncMaxState();
  window.addEventListener("resize", handleResize);
  document.addEventListener("keydown", handleCommandHotkey);
});

onBeforeUnmount(() => {
  window.removeEventListener("resize", handleResize);
  document.removeEventListener("keydown", handleCommandHotkey);
});

const handleResize = () => {
  syncMaxState();
  WindowGetSize().then((size) => {
    throttleUpdate(size.w, size.h);
  });
};

const throttleUpdate = throttle((width: number, height: number) => {
  SaveWindowsScreenSize(width, height);
}, 1000);

function toggleMaximise() {
  WindowToggleMaximise();
  syncMaxState();
}

function setTitle(path: string) {
  switch (path) {
    case "/":
      return "Home";
    default:
      return path.split("/").slice(-1)[0];
  }
}

const isDark = useDark({
  storageKey: "theme",
  valueDark: "dark",
  valueLight: "light",
});

const toggle = useToggle(isDark);
function toggleTheme() {
  global.Theme.value = !global.Theme.value;
  toggle();
}

const { locale, t } = useI18n();
global.Language.value = locale.value;
function changeLanguage(lang: string) {
  localStorage.setItem("language", lang);
  locale.value = lang;
}

const router = useRouter();
const route = useRoute();

const allTools = computed(() => {
  const list: any[] = [];
  MenuList.forEach((group) => {
    if (group.children) {
      group.children.forEach((item) => {
        list.push({
          ...item,
          parentPath: group.path,
          groupName: group.name,
        });
      });
    } else {
      list.push({
        ...group,
        parentPath: "",
        groupName: group.name,
      });
    }
  });
  return list;
});

const matchTools = (keyword: string) => {
  const value = keyword.trim();
  if (!value) return [];
  const isZh = locale.value.toLowerCase().startsWith("zh");
  const lowerKeyword = value.toLowerCase();
  return allTools.value.filter((item) => {
    const localizedName = t(item.name);
    if (isZh) {
      return localizedName.includes(value) || item.name.includes(value);
    }
    return (
      localizedName.toLowerCase().includes(lowerKeyword) ||
      item.name.toLowerCase().includes(lowerKeyword)
    );
  });
};

const commandPalette = computed(() => [
  {
    value: "> zh",
    label: t("titlebar.switchZh", "切换到中文"),
    desc: t("titlebar.switchDesc", "设置界面语言为中文"),
    type: "command",
    action: () => changeLanguage("zh"),
  },
  {
    value: "> en",
    label: t("titlebar.switchEn", "Switch to English"),
    desc: t("titlebar.switchDescEn", "Set interface language to English"),
    type: "command",
    action: () => changeLanguage("en"),
  },
  {
    value: "> theme",
    label: t("titlebar.toggleTheme", "切换主题"),
    desc: t("titlebar.toggleThemeDesc", "在明亮/暗色主题间切换"),
    type: "command",
    action: () => toggleTheme(),
  },
  {
    value: "> console",
    label: t("titlebar.openConsole", "打开日志"),
    desc: t("titlebar.openConsoleDesc", "查看运行日志"),
    type: "command",
    action: () => {
      showLogger.value = true;
    },
  },
]);

const handleSelect = (item: any) => {
  if (item?.type === "command" && typeof item.action === "function") {
    item.action();
    searchText.value = "";
    return;
  }

  if (item?.fullPath) {
    router.push(item.fullPath);
    searchText.value = "";
  }
};

const placeholder = computed(() => setTitle(route.path));

const helperItems = computed(() => [
  {
    value: ">",
    label: t("titlebar.helperCommandHint", "输入 > 以显示可用命令"),
    group: "",
    hintKeys: global.temp.isMacOS ? ["⌘", "⇧", "P"] : ["Ctrl", "Shift", "P"],
    type: "helper",
  },
]);

const defaultMenus = computed(() => {
  const popularPaths = ["/Webscan", "/Dirsearch", "/FOFA", "/CyberChef"];
  return popularPaths
    .map((path) => {
      const item = allTools.value.find((tool) => tool.path === path);
      if (!item) return null;
      return {
        value: t(item.name),
        label: t(item.name),
        group: t(item.groupName),
        fullPath: item.parentPath + item.path,
        type: "tool",
      };
    })
    .filter(Boolean) as any[];
});

const findFirstCommand = (text: string) => {
  if (!text.trim().startsWith(">")) return null;
  const target = text.replace(">", "").trim().toLowerCase();
  return commandPalette.value.find((item) => {
    const value = item.value.toLowerCase();
    if (!target) return value.startsWith(">");
    return value.includes(target);
  });
};

const fetchSuggestions = (queryString: string, cb: (results: any[]) => void) => {
  const keyword = queryString.trim();

  if (keyword.startsWith(">")) {
    const target = keyword.replace(">", "").trim().toLowerCase();
    const seen = new Set<string>();
    const combined = [...commandPalette.value].filter((item) => {
      const value = item.value.toLowerCase();
      const match = target ? value.includes(target) : value.startsWith(">");
      if (!match) return false;
      if (seen.has(value)) return false;
      seen.add(value);
      return true;
    });
    cb(combined);
    return;
  }

  const list = matchTools(keyword).map((item) => ({
    value: t(item.name),
    label: t(item.name),
    group: t(item.groupName),
    fullPath: item.parentPath + item.path,
    type: "tool",
  }));

  if (!keyword) {
    cb([...helperItems.value, ...defaultMenus.value, ...list]);
    return;
  }

  cb(list);
};

const openFirst = () => {
  const text = searchText.value.trim();
  if (!text || text === ">") return; // 只查看列表时不触发执行

  const command = findFirstCommand(text);
  if (command) {
    handleSelect(command);
    return;
  }
  const first = matchTools(text)[0];
  if (first) {
    handleSelect({ fullPath: first.parentPath + first.path });
  }
};

const handleCommandHotkey = (e: KeyboardEvent) => {
  const key = e.key.toLowerCase();
  if ((e.metaKey || e.ctrlKey) && key === "k") {
    e.preventDefault();
    nextTick(() => {
      (searchRef.value as any)?.focus?.();
    });
    return;
  }
  if ((e.metaKey || e.ctrlKey) && e.shiftKey && key === "p") {
    e.preventDefault();
    searchText.value = ">";
    nextTick(() => {
      (searchRef.value as any)?.focus?.();
    });
  }
};

const macCommandKey = "\u2318";
const shortcutKeys = computed(() => (global.temp.isMacOS ? [macCommandKey, "K"] : ["Ctrl", "K"]));
</script>

<template>
  <div class="titlebar" :style="titlebarStyle">
    <div class="left">
      <div class="left-actions">
        <div class="nav-actions">
          <el-tooltip v-for="item in routerControl" :content="$t(item.label)">
            <el-button text class="custom-button" @click="item.action">
              <el-icon :size="16">
                <component :is="item.icon" />
              </el-icon>
            </el-button>
          </el-tooltip>
        </div>
      </div>
    </div>

    <div class="center" @dblclick="toggleMaximise">
      <div class="search-container" @dblclick.stop>
        <el-autocomplete ref="searchRef" v-model="searchText" :fetch-suggestions="fetchSuggestions"
          :placeholder="placeholder" :trigger-on-focus="true" size="small" :prefix-icon="Search" class="title-search"
          @select="handleSelect" @keyup.enter.native="openFirst">
          <template #default="{ item }">
            <div :class="['search-row', item.type]">
              <span class="name">
                {{ item.label }}
              </span>
              <span class="group" v-if="item.group || item.desc">{{ item.group || item.desc }}</span>
              <div class="helper-shortcut" v-if="item.type === 'helper' && item.hintKeys">
                <span class="keycap" v-for="(key, idx) in item.hintKeys" :key="key + idx">{{ key }}</span>
              </div>
            </div>
          </template>
          <template #suffix>
            <div class="shortcut-hint">
              <span class="keycap" v-for="(key, idx) in shortcutKeys" :key="key + idx">{{ key }}</span>
            </div>
          </template>
        </el-autocomplete>
      </div>
    </div>

    <div class="right flex">
      <div class="action-buttons">
        <el-tooltip :content="$t('titlebar.app_launcher')">
          <el-button class="custom-button" text @click="$router.push('/AppLauncher')">
            <template #icon>
              <el-icon :size="16">
                <runnerIcon />
              </el-icon>
            </template>
          </el-button>
        </el-tooltip>
      </div>
      <el-divider direction="vertical" />
      <el-button-group class="window-controls">
        <el-button v-for="item in windowsControl" :class="item.class!" text @click="item.action">
          <template #icon>
            <el-icon size="16">
              <component :is="item.icon" />
            </el-icon>
          </template>
        </el-button>
      </el-button-group>
    </div>
  </div>

  <el-drawer v-model="showLogger" :title="$t('titlebar.yx_log')" direction="rtl" size="50%">
    <div class="log-textarea" v-html="global.Logger.value"></div>
  </el-drawer>
</template>

<style scoped>
.titlebar {
  display: flex;
  align-items: center;
  gap: 10px;
  width: 100%;
  height: var(--titlebar-height);
  padding: 0 0 0 6px;
  background: var(--el-bg-color);
  border-bottom: 1px solid var(--el-border-color-lighter);
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.06);
  box-sizing: border-box;
  --wails-draggable: drag;
  cursor: default;

  .custom-button {
    margin: 0;
    height: 30px;
    width: 32px;
    border-radius: 6px;
    box-shadow: none;
    transition: background 0.12s ease;
    --wails-draggable: no-drag;
  }
}

.titlebar .el-button-group {
  gap: 0;
}

.nav-actions {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.action-buttons {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.window-controls {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.window-controls :deep(.el-button) {
  height: var(--titlebar-height);
  border-radius: 0;
}

.window-controls :deep(.el-button.is-text.close:hover) {
  background-color: red;
  color: #fff;
}

.custom-button:hover {
  background: var(--el-fill-color-light);
}

.left {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
  --wails-draggable: drag;
}

.left-actions {
  display: flex;
  align-items: center;
  gap: 8px;
  --wails-draggable: no-drag;
}

.center {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: row;
  --wails-draggable: drag;
}

.search-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  color: var(--el-text-color-secondary);
}

.search-label {
  font-weight: 600;
  font-size: 12px;
  letter-spacing: 0.2px;
}

.pill {
  padding: 4px 10px;
  border-radius: 10px;
  background: #252526;
  font-size: 11px;
  color: var(--el-text-color-primary);
  border: 1px solid #303233;
}

.search-container {
  --wails-draggable: no-drag;
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  max-width: 460px;
  min-width: 260px;
  padding: 0;
  background: transparent;
  border-radius: 0;
  border: none;
  box-shadow: none;
}

.title-search {
  width: 100%;
  max-width: 460px;
  min-width: 260px;
  flex: 1 1 auto;
  display: inline-flex;
}

.title-search :deep(.el-input__wrapper) {
  border-radius: 8px;
  height: 30px;
  padding: 0 10px;
  background: var(--el-fill-color-light);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.04), inset 0 0 0 1px var(--el-border-color);
}

.title-search :deep(.el-input__inner) {
  color: var(--el-text-color-primary);
  line-height: 30px;
}

.title-search :deep(.el-input__inner::placeholder) {
  color: var(--el-text-color-secondary);
}

.title-search :deep(.el-input__prefix),
.title-search :deep(.el-input__suffix) {
  color: var(--el-text-color-secondary);
}

.search-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 10px;
}

.search-row .name {
  color: var(--el-text-color-primary);
  font-weight: 600;
  font-size: 13px;
}

.search-row .group {
  color: var(--el-text-color-secondary);
  font-size: 12px;
}

.shortcut-hint {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  margin-left: 8px;
}

.keycap {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 16px;
  height: 16px;
  padding: 0 5px;
  border-radius: 3px;
  background: var(--el-fill-color);
  box-shadow: inset 0 -1px 0 rgba(0, 0, 0, 0.12);
  border: 1px solid var(--el-border-color-lighter);
  font-size: 10px;
  color: var(--el-text-color-regular);
  line-height: 1;
}

.helper-shortcut {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.right {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 6px;
  flex: 1;
  --wails-draggable: drag;
}

.action-buttons {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.window-controls {
  display: inline-flex;
  align-items: center;
  gap: 10px;
}

html.light .el-button.is-text:hover {
  background-color: #ededed;
}

.el-button.is-text.close:hover {
  background-color: red;
}
</style>
