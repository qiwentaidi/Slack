<template>
    <div class="note" :class="noteClass" v-if="visible">
        <div class="note-header">
            <el-icon class="note-icon">
                <component :is="noteIcon" />
            </el-icon>
            <span class="note-title">{{ noteTitle }}</span>
        </div>
        <div class="note-content">
            <slot></slot>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { computed, ref } from "vue";
import { InfoFilled, WarningFilled, SuccessFilled, CircleCloseFilled } from "@element-plus/icons-vue";


const props = defineProps({
    type: {
        type: String as () => string, // 明确指定类型为 string
        default: "info",
        validator: (value: string) => ["info", "warning", "error", "success"].includes(value), // 明确指定参数类型为 string
    },
});

const visible = ref(true);


const noteClass = computed(() => `note-${props.type}`);

// 计算图标
const noteIcon = computed(() => {
    switch (props.type) {
        case "warning":
            return WarningFilled;
        case "error":
            return CircleCloseFilled;
        case "success":
            return SuccessFilled;
        default:
            return InfoFilled;
    }
});

// 计算标题
const noteTitle = computed(() => {
    switch (props.type) {
        case "warning":
            return "Warning";
        case "error":
            return "Error";
        case "success":
            return "Success";
        default:
            return "Note";
    }
});
</script>

<style scoped>
.note {
    padding: 12px;
    border-radius: 5px;
    font-size: 15px;
    font-weight: 500;
    border-left: 5px solid;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    transition: all 0.3s ease-in-out;
}

.note-header {
    display: flex;
    align-items: center;
    font-weight: bold;
    margin-bottom: 8px;
}

.note-title {
    font-weight: bold;
    font-size: large;
}

.note-icon {
    font-size: 18px;
    margin-right: 6px;
}
</style>