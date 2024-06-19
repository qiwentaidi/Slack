<script setup lang="ts">
import { InitConfigFile } from "../config";
import global from "../global";
import { onMounted, ref, computed } from "vue";
import MenuList from "../router/menu";
import { Search, Box, Plus, CircleCheckFilled, InfoFilled } from "@element-plus/icons-vue";
import { BrowserOpenURL } from "../../wailsjs/runtime/runtime";
import LocalNavigation from "../components/LocalNavigation.vue"
import { ElMessage } from "element-plus";
// 初始化时调用
onMounted(async () => {
  await InitConfigFile(1000);
});

const searchFilter = ref("");
const activeCard = ref("app");
const filteredOptions = computed(() => {
  if (!searchFilter.value) {
    return global.onlineOptions;
  }
  return global.onlineOptions.map((group) => ({
    ...group,
    value: group.value.filter((item) =>
      item.name.toLowerCase().includes(searchFilter.value.toLowerCase())
    ),
  }));
});


const localNavigationRef = ref();

function addGroup() {
  if (localNavigationRef.value) {
    (localNavigationRef.value as any).addGroup();
  }
}

function save() {
  console.log(global.webscan)
  localStorage.setItem('webscan', JSON.stringify(global.webscan))
  ElMessage({
    type: 'success',
    message: '保存成功',
  })
}
</script>

<template>
  <div style="position: relative">
    <el-tabs v-model="activeCard">
      <!-- 应用导航 -->
      <!-- 应用导航 -->
      <!-- 应用导航 -->
      <el-tab-pane :label="$t('aside.app_navigation')" name="app">
        <el-scrollbar height="85vh">
          <template v-for="groups in MenuList">
            <div v-if="groups.children" class="card-group">
              <el-card>
                <div style="margin-bottom: 10px">
                  <span style="font-weight: bold">{{ $t(groups.name) }}</span>
                </div>
                <el-row>
                  <el-col :lg="4" v-for="item in groups.children">
                    <el-result :title="$t(item.name)" @click="$router.push(item.path)">
                      <template #icon>
                        <el-image class="appNavIcon" :src="item.icon" />
                      </template>
                    </el-result>
                  </el-col>
                </el-row>
              </el-card>
            </div>
          </template>
        </el-scrollbar>
      </el-tab-pane>
      <!-- 在线导航 -->
      <!-- 在线导航 -->
      <!-- 在线导航 -->
      <el-tab-pane :label="$t('aside.navigation')" name="online">
        <el-scrollbar height="85vh">
          <div v-for="groups in filteredOptions" class="card-group">
            <el-card>
              <div style="margin-bottom: 10px">
                <span style="font-weight: bold">{{ $t(groups.label) }}</span>
              </div>
              <el-row>
                <el-col :lg="4" v-for="item in groups.value">
                  <el-tooltip :content="item.url" placement="top">
                    <el-result :title="item.name" @click="BrowserOpenURL(item.url)">
                      <template #icon>
                        <el-image class="onlineNavIcon" :src="item.icon" />
                      </template>
                    </el-result>
                  </el-tooltip>
                </el-col>
              </el-row>
            </el-card>
          </div>
        </el-scrollbar>
      </el-tab-pane>
      <!-- 本地导航 -->
      <!-- 本地导航 -->
      <!-- 本地导航 -->
      <el-tab-pane name="local">
        <template #label>
            <span class="custom-tabs-label">
              <span>{{ $t('aside.local_navigation') }}</span>
              <el-tooltip placement="right">
            <template #content>
              1、添加组后会出现卡片组，名称不支持重复<br />
              2、卡片新增后可以从桌面或者文件夹拖入元素或者右上角添加元素<br />
              3、右键元素可以编辑已有信息、打开文件夹、删除<br />
              4、Java Verison 代表运行jar包的java环境变量
            </template>
              <el-icon>
                <InfoFilled />
              </el-icon>
            </el-tooltip>
            </span>
        </template>
        <LocalNavigation ref="localNavigationRef" />
      </el-tab-pane>
    </el-tabs>
    <div class="custom_eltabs_titlebar" v-show="activeCard == 'online'">
      <el-input class="input-filter" :suffix-icon="Search" v-model="searchFilter" @input="filteredOptions"
        placeholder="Filter search"></el-input>
    </div>
    <div class="custom_eltabs_titlebar" v-show="activeCard == 'local'">
      <el-input v-model="global.webscan.java" style="width: 250px; margin-right: 5px;">
        <template #prepend>
          Java Version
        </template>
        <template #suffix>
          <el-button @click="save" :icon="CircleCheckFilled" link>
          </el-button>
        </template>
      </el-input>
      <el-button-group>
        <el-button :icon="Box" @click="addGroup">{{ $t("navigator.add_group") }}</el-button>
        <el-button :icon="Plus" @click="global.temp.localAddItem = true">{{ $t("navigator.add_item") }}</el-button>
      </el-button-group>
    </div>
  </div>
</template>

<style>
.input-filter {
  width: 240px;
  margin-right: 1%;
}

.onlineNavIcon .el-image__inner {
  height: 25px;
  width: 25px;
}

.el-result:hover {
  cursor: pointer;
  background-color: #f5f7fa;
}

.appNavIcon .el-image__inner {
  height: 30px;
  width: 30px;
}

.card-group {
  margin-top: 10px;
  width: 99%;
}

.custom-tabs-label .el-icon {
  vertical-align: middle;
}

.custom-tabs-label span {
  vertical-align: middle;
  margin-left: 4px;
}
</style>
