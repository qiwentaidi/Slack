<template>
  <!-- 标签页 -->
  <el-tabs v-model="activeTab" type="card" editable @edit="handleTabsEdit">
    <el-tab-pane v-for="tab in tabs" :key="tab.name" :label="tab.title" :name="tab.name">
      <!-- 选择目录 -->
      <div v-if="tab.status === 'select'" class="select-dir-container">
        <el-card class="select-dir-card">
          <div class="select-dir-header">
            <el-icon class="icon-folder">
              <Folder />
            </el-icon>
            <h2 class="title">选择文件夹</h2>
            <el-link @click="settingDialog = true">存在敏感词的文件会标红, 更改检索规则<el-icon class="el-icon--right">
                <Setting />
              </el-icon></el-link>
          </div>
          <div class="select-dir-body">
            <el-input v-model="tab.path" placeholder="请选择或输入文件夹路径" size="large">
              <template #suffix>
                <el-button link :icon="FolderOpened" @click="selectDirectory(tab)" class="icon-btn" />
              </template>
            </el-input>
            <el-button type="primary" size="large" class="btn-generate" @click="generateTree(tab)">
              开始检索
            </el-button>
          </div>
        </el-card>
      </div>

      <!-- 文件树 -->
      <div v-else class="tree-container">
        <div></div>
        <splitpanes class="default-theme">
          <pane min-size="30">
            <div class="my-header" style="margin-bottom: 20px; margin-right: 10px;">
              <el-tooltip :content="tab.path">
                <span class="dirtips">物理路径: {{ tab.path }}</span>
              </el-tooltip>
              <el-button-group>
                <el-tooltip content="返回选择">
                  <el-button :icon="Back" link @click="backToSelect(tab)" />
                </el-tooltip>
                <el-tooltip content="刷新资源管理器">
                  <el-button :icon="RefreshLeft" link @click="generateTree(tab); ElMessage.success('刷新成功')" />
                </el-tooltip>
                <el-tooltip content="展开/折叠所有文件夹">
                  <el-button :icon="Film" link @click="changeNodeAll(tab)"></el-button>
                </el-tooltip>
              </el-button-group>
            </div>
            <el-tree-v2 ref="treeRef" :data="tab.treeData" :props="treeProps" node-key="id" @node-click="nodeClick"
              :height="treeHeight" :highlight-current="true">
              <template #default="{ node }">
                <el-icon style="margin-right: 5px;">
                  <FolderOpened class="file-color" v-if="node.data.isDir" />
                  <Document v-else />
                </el-icon>
                <span :class="{ 'highlight-file': node.data.hits }">{{ node.label }}</span>
              </template>
            </el-tree-v2>
          </pane>
          <pane min-size="20" style="margin-top: -12px; height: calc(100% + 12px);">
            <div v-if="currentFileHits && Object.keys(currentFileHits).length" class="keyword-matches">
              <span class="match-title">关键词匹配:</span>
              <el-tag v-for="(count, keyword) in currentFileHits" :key="keyword" type="danger" class="match-tag">
                {{ keyword }}: {{ count }} 次
              </el-tag>
            </div>
            <highlightjs :code="code"></highlightjs>
          </pane>
        </splitpanes>
      </div>
    </el-tab-pane>
  </el-tabs>
  <el-drawer v-model="settingDialog" size="50%">
    <template #header>
      <span class="drawer-title">设置检索规则</span>
    </template>
    <el-form :modle="global.fileRetrieval" label-width="auto">
      <el-form-item label="检索关键词:">
        <el-input v-model="global.fileRetrieval.keywords" type="textarea" :rows="5"></el-input>
        <span class="form-item-tips">以,分割关键词</span>
      </el-form-item>
      <el-form-item label="文件后缀黑名单:">
        <el-input v-model="global.fileRetrieval.blackList" type="textarea" :rows="5"></el-input>
        <span class="form-item-tips">以,分割关键词</span>
      </el-form-item>
      <el-form-item style="float: right;">
        <el-button type="primary" @click="SaveConfig">保存</el-button>
      </el-form-item>
    </el-form>
  </el-drawer>
</template>

<script lang="ts" setup>
import { onMounted, onUnmounted, ref } from 'vue';
import { Document, FolderOpened, Back, RefreshLeft, Film, Setting } from '@element-plus/icons-vue';
import { ElMessage } from 'element-plus';
import { DirectoryDialog, BuildTree, CheckFileStat, ReadFile, ListDir } from 'wailsjs/go/services/File';
import { TreeNodeData } from 'element-plus/es/components/tree-v2/src/types';
import { Splitpanes, Pane } from 'splitpanes'
import { DirectoryTab, TreeNode } from '@/stores/interface';
import global from '@/stores/index';
import { SaveConfig } from '@/config';

// ================   动态计算文件树高度 ================

onMounted(() => {
  calculateTreeHeight();
  window.addEventListener("resize", calculateTreeHeight);
});

onUnmounted(() => {
  window.removeEventListener("resize", calculateTreeHeight);
});

const treeHeight = ref(0);
const calculateTreeHeight = () => {
  const fixedOffset = 190;
  treeHeight.value = window.innerHeight - fixedOffset;
};

// ================   动态计算文件树高度 ================

const activeTab = ref('tab-1');
const tabIndex = ref(1);
const treeProps = { value: 'id', label: 'label', children: 'children' };

const tabs = ref<DirectoryTab[]>([
  { name: 'tab-1', title: '标签页 1', path: '', status: 'select', treeData: [], isCollapse: false },
]);
function handleTabsEdit(targetName: string | undefined, action: 'remove' | 'add') {
  if (action === 'add') {
    const newTabName = `tab-${++tabIndex.value}`;
    const newTab: DirectoryTab = {
      name: newTabName,
      title: `标签页 ${tabIndex.value}`,
      path: '',
      status: 'select',
      treeData: [],
      isCollapse: false
    };
    tabs.value.push(newTab);
    activeTab.value = newTabName;
  } else if (action === 'remove') {
    const targetIndex = tabs.value.findIndex(tab => tab.name === targetName);
    if (tabs.value.length === 1) {
      ElMessage.warning('至少保留一个标签页');
      return;
    }

    // 保证删除后有一个活动的标签页
    let activeTabName = activeTab.value;
    if (activeTabName === targetName) {
      const nextTab = tabs.value[targetIndex + 1] || tabs.value[targetIndex - 1];
      if (nextTab) {
        activeTabName = nextTab.name;
      }
    }
    activeTab.value = activeTabName;
    tabs.value = tabs.value.filter(tab => tab.name !== targetName);
  }
};
async function selectDirectory(tab: DirectoryTab) {
  let path = await DirectoryDialog();
  if (!path) return;
  tab.path = path;
};

async function generateTree(tab: DirectoryTab) {
  if (!tab.path) {
    ElMessage.warning('请选择文件夹路径');
    return;
  }

  try {
    const data = await BuildTree(tab.path, global.fileRetrieval.keywords.split(','), global.fileRetrieval.blackList.split(','));
    tab.treeData = Array.isArray(data) ? data : [data]; // 确保是数组
    tab.status = 'tree';
  } catch (error) {
    ElMessage.error('构建文件树失败: ' + error);
  }
};

function backToSelect(tab: DirectoryTab) {
  tab.status = 'select';
  tab.treeData = [];
};

const code = ref('');

const currentFileHits = ref<Record<string, number> | null>(null);
async function nodeClick(data: TreeNodeData, node: TreeNode, e: MouseEvent) {
  if (data.isDir) return;

  let isStat = await CheckFileStat(data.id);
  if (!isStat) {
    ElMessage.error('文件不存在');
    return;
  }

  let file = await ReadFile(data.id);
  if (file.Content.length > 100000) {
    code.value = '文件内容过大，不支持查看';
    currentFileHits.value = null;
    return;
  }

  // 设置当前文件的匹配关键词
  currentFileHits.value = data.hits || {};
  code.value = file.Content;
}

let treeRef = ref();
// 折叠/展开所有节点
async function changeNodeAll(tab: DirectoryTab) {
  // 获取当前标签页索引
  let index = tabs.value.findIndex(tab => tab.name === activeTab.value);
  let treeKeys = await ListDir(tab.path)
  for (const treeKey of treeKeys) {
    let node = treeRef.value[index]?.getNode(treeKey)
    tab.isCollapse ? treeRef.value[index]?.collapseNode(node) : treeRef.value[index]?.expandNode(node)
  }
  tab.isCollapse = !tab.isCollapse
}

const settingDialog = ref(false);
</script>

<style scoped>
.tree-container {
  border-radius: 8px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  padding: 10px;
}

.select-dir-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 600px;
}

.select-dir-card {
  width: 600px;
  border-radius: 8px;
  padding: 20px;
}

.select-dir-header {
  text-align: center;
  margin-bottom: 20px;
}

.select-dir-header .icon-folder {
  font-size: 60px;
  color: #3b82f6;
}

.select-dir-header .title {
  font-size: 24px;
  font-weight: 600;
}

.select-dir-header .description {
  color: #6b7280;
  margin-top: 10px;
}

.select-dir-body {
  display: flex;
  gap: 12px;
  align-items: center;
}

.btn-generate {
  font-size: 16px;
  font-weight: 500;
}

.file-color {
  color: var(--el-color-primary);
}

::v-deep(pre code.hljs) {
  display: block;
  overflow-x: auto;
  padding: 10px;
  margin-left: 5px;
  white-space: pre-wrap;
  word-wrap: break-word;
  word-break: break-all;
  border-radius: 5px;
  height: 100%;
}

::v-deep(code.hljs) {
  padding: 0
}

.splitpanes {
  height: calc(100vh - 155px);
}

.dirtips {
  color: var(--el-text-color);
  font-size: 13px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 85%;
}

.highlight-file {
  color: #e74c3c;
  /* 红色高亮 */
  font-weight: bold;
}

.keyword-matches {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 10px;
}

.match-title {
  font-weight: bold;
  color: #333;
  margin-left: 5px;
}

.match-tag {
  font-size: 12px;
}
</style>