<template>
    <div @contextmenu.prevent="showMenu" style="width: 100%;">
        <slot>

        </slot>
    </div>

    <div v-if="visible" :style="menuStyle" class="context-menu">
        <slot name="menu">

        </slot>
    </div>
</template>

<script lang="ts">
import { defineComponent, ref, onMounted, onBeforeUnmount } from 'vue';

interface MenuItem {
    label: string;
    action: () => void;
}

export default defineComponent({
    name: 'ContextMenu',
    props: {

    },
    setup(props) {
        const visible = ref(false);
        const menuStyle = ref({
            top: '0px',
            left: '0px',
        });

        const showMenu = (event: MouseEvent) => {
            event.preventDefault();
            visible.value = true;
            menuStyle.value = {
                top: `${event.clientY}px`,
                left: `${event.clientX}px`,
            };
        };

        const hideMenu = () => {
            visible.value = false;
        };

        const handleItemClick = (item: MenuItem) => {
            item.action();
            hideMenu();
        };

        onMounted(() => {
            window.addEventListener('click', hideMenu);
        });

        onBeforeUnmount(() => {
            window.removeEventListener('click', hideMenu);
        });

        return {
            visible,
            menuStyle,
            handleItemClick,
            showMenu,
        };
    },
});
</script>

<style scoped>
.context-menu {
    position: absolute;
    background-color: white;
    border: 1px solid #ccc;
    border-radius: 5px;
    padding: 5px;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
    z-index: 1000;
}

.context-menu ul {
    list-style-type: none;
    padding: 0;
    margin: 0;
}

.context-menu li {
    padding: 8px 12px;
    cursor: pointer;
}

.context-menu li:hover {
    background-color: #f5f5f5;
}
</style>