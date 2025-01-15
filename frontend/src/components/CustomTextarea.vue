<template>
    <div class="textarea-container">
        <el-input
            v-model="localValue"
            type="textarea"
            :rows="rows"
            :placeholder="placeholder"
            :readonly="readonly"
            :resize="resize"
            :minRows="2"
            @input="handleInput"
        ></el-input>

        <el-button-group v-if="showActions && !readonly" class="action-area">
            <el-button :icon="UploadIcon" size="small" @click="handleUpload">
                Upload
            </el-button>
            <el-button :icon="CloseIcon" size="small" @click="clearInput"></el-button>
        </el-button-group>
        <el-button v-else :icon="CopyIcon" size="small" @click="handleCopy" class="action-area">
            Copy
        </el-button>
    </div>
</template>

<script lang="ts">
import { defineComponent, ref, watch } from "vue";
import { ElMessage } from "element-plus";
import { Upload, CloseBold, DocumentCopy } from "@element-plus/icons-vue";
import { FileDialog, ReadFile } from "wailsjs/go/services/File";
import { Copy } from "@/util";

export default defineComponent({
    name: "CustomTextarea",
    props: {
        modelValue: {
            type: String,
            required: true,
        },
        rows: {
            type: Number,
            default: 6,
        },
        showActions: {
            type: Boolean,
            default: true,
        },
        readonly: {
            type: Boolean,
            default: false,
        },
        resize: {
            type: String,
            default: "vertical",
        },
        placeholder: {
            type: String,
            default: "",
        },
    },
    emits: ["update:modelValue"],
    setup(props, { emit }) {
        // 创建一个中间变量
        const localValue = ref(props.modelValue);

        // 监听父组件传入的 modelValue，当它发生变化时更新 localValue
        watch(
            () => props.modelValue,
            (newValue) => {
                localValue.value = newValue;
            }
        );

        const handleInput = (value: string) => {
            emit("update:modelValue", value); // 通知父组件更新
        };

        const handleUpload = async () => {
            try {
                const filepath = await FileDialog("*.txt");
                if (!filepath) return;

                const file = await ReadFile(filepath);
                if (file.Error) {
                    ElMessage({ type: "warning", message: file.Message });
                    return;
                }
                emit("update:modelValue", file.Content!);
            } catch (error) {
                console.error("Upload error:", error);
                ElMessage.error("Failed to upload file.");
            }
        };

        const clearInput = () => {
            localValue.value = "";
            emit("update:modelValue", "");
        };

        const handleCopy = () => {
            Copy(props.modelValue);
        };

        return {
            localValue,
            handleInput,
            handleUpload,
            clearInput,
            handleCopy,
            UploadIcon: Upload,
            CloseIcon: CloseBold,
            CopyIcon: DocumentCopy,
        };
    },
});
</script>