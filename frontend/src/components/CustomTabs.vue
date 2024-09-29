<script lang="ts">
import { defineComponent, PropType } from 'vue';
export default defineComponent({
    name: 'CustomTabs',
    props: {
        filterLeftPosition: {
            type: [Number, String] as PropType<number | string>,
            default: 100,
        },
    },
})
</script>


<template>
    <div class="custom-tabs-container">
        <slot></slot>
        <div class="filter-container" :style="{ left: typeof filterLeftPosition === 'number' ? `${filterLeftPosition}px` : filterLeftPosition }">
            <slot name="filter"></slot>
        </div>
        <div class="custom-tabs-titlebar">
            <slot name="ctrl">
            </slot>
        </div>
    </div>
</template>

<style scoped>
.custom-tabs-container {
    position: relative;
}

.filter-container {
    position: absolute;
    top: 4px;
}

.custom-tabs-titlebar {
    position: absolute;
    right: 5px;
    top: 4px;
}
</style>