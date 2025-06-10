<script setup lang="ts">
import MenuList from "@/router/menu";
import { ref } from "vue";
const activeNames = ref(['1', '2', '3', '4'])
</script>

<template>
    <el-collapse v-model="activeNames" class="el-collapse-parent">
        <template v-for="(groups, index) in MenuList">
            <el-collapse-item v-if="groups.children" :name="index.toString()">
                <template #title>
                    <div class="custom-header">
                        <el-icon :size="24">
                            <component :is="groups.icon" />
                        </el-icon>
                        <span>{{ $t(groups.name) }}</span>
                    </div>
                </template>

                <div style="display: flex; gap: 10px; flex-wrap: wrap; padding-top: 10px;">
                    <el-card shadow="hover" class="card" v-for="item in groups.children" :key="item.path"
                        @click="$router.push(groups.path + item.path)">
                        <div class="flex">
                            <div class="card-content">
                                <!-- 左侧图片 -->
                                <el-image :src="item.icon" style="width: 24px;" />
                                <!-- 右侧内容 -->
                                <span>{{ $t(item.name) }}</span>
                            </div>
                            <el-icon :size="30" class="location-icon">
                                <DArrowRight />
                            </el-icon>
                        </div>
                    </el-card>
                </div>
            </el-collapse-item>
        </template>
    </el-collapse>
</template>

<style scoped>
.el-collapse-parent {
    border: none;
    width: 100%;
    :deep(.el-collapse-item__header) {
        background: var(--sidebar-bg-color);
        border-radius: 4px 4px 0px 0px;
        border-bottom: 1px solid var(--collapse-border-color);
    }
    :deep(.el-collapse-item__content) {
        padding-bottom: 10px;
    }
}

.custom-header {
    display: flex;
    align-items: center;
    font-size: large;
    font-weight: bold;
    padding: 10px 10px;
}

.custom-header span {
    margin-left: 10px;
}

.card {
    width: calc(25% - 10px);
    position: relative;
    /* 确保图标相对于此容器进行绝对定位 */

    .card-content {
        display: flex;
        align-items: center;
        text-align: left;
    }

    span {
        margin-left: 5px;
        font-weight: bold;
        font-size: 15px;
    }

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


/* 跳转图标 */
.location-icon {
    position: absolute;
    opacity: 0;
    right: 10px;
    transform: translateX(20px);
    /* 初始位置在右侧外 */
    transition: transform 0.3s ease, opacity 0.3s ease;
}

:deep(.el-card__body) {
    padding: var(--el-card-padding);
}
</style>
