<script lang="ts">
import { defineComponent } from 'vue';
export default defineComponent({
    name: 'CustomTabs',
})
</script>


<template>
    <div style="position: relative;">
        <slot></slot>
        <div class="custom_eltabs_titlebar">
            <slot name="ctrl">
            </slot>
        </div>
    </div>
</template>

<style scoped>


</style>