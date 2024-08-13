<script setup lang="ts">
import MenuList from "@/router/menu";
</script>

<template>
  <el-scrollbar height="92vh">
    <div v-for="groups in MenuList">
      <div v-if="groups.children">
        <el-divider><span style="font-size: large;">{{ $t(groups.name) }}</span></el-divider>
        <div style="display: flex; gap: 10px; flex-wrap: wrap;">
          <el-card shadow="hover" class="card" v-for="item in groups.children" @click="$router.push(groups.path + item.path)">
            <div style="display: flex;">
              <div class="card-content">
                <!-- 左侧图片 -->
                <el-image class="appNavIcon" :src="item.icon" />
                <!-- 右侧内容 -->
                <span class="card-text">{{ $t(item.name) }}</span>
              </div>
              <el-icon :size="30" class="location-icon">
                <DArrowRight />
              </el-icon>
            </div>
          </el-card>
        </div>
      </div>
    </div>
  </el-scrollbar>
</template>

<style>
.card {
  width: 24%;
  position: relative; /* 确保图标相对于此容器进行绝对定位 */
}

.card:hover {
  cursor: pointer;

  .location-icon {
    opacity: 1;
    transform: translateX(0);
    /* 从右向左移动到正常位置 */
  }

  .card-content {
    filter: blur(1px);
    /* 图片虚化 */
  }
}


.card-content {
  display: flex;
  align-items: center;
  text-align: left;
}


.card-text {
  margin-left: 5px;
  font-weight: bold;
}

/* 跳转图标 */
.location-icon {
  position: absolute;
  opacity: 0;
  right: 10px;
  transform: translateX(20px);
  /* 初始位置在右侧外 */
  transition: transform 0.3s ease, opacity 0.3s ease;
}

.appNavIcon .el-image__inner {
  height: 25px;
  width: 25px;
}
</style>
