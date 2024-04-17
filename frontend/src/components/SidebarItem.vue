<script setup lang="ts">
import { RouteRecordRaw } from 'vue-router';
// 做类型限制，解决ts类型报错
type CustomRouteRecordRaw = RouteRecordRaw & {
	meta: {
		isShow?: boolean;
	};
};
const props = defineProps({
	// 拿到父组件传递过来的路由列表进行渲染
	routerList: {
		type: Array as () => CustomRouteRecordRaw[],
		required: true
	}
});
</script>
<template>
	<template v-for="item in props.routerList" :key="item.path">
		<!-- 当该菜单项有子菜单时 -->
		<el-sub-menu :index="item.path" v-if="item.children && item.children.length > 0">
			<template #title v-if="item.meta.icon">
				<!-- 菜单项图标，我此处用的是全局封装的 svg组件 -->
				<el-icon><svg-icon :icon-class="item.meta.icon" /></el-icon>
				<!-- 菜单项名称，在路由中定义好 -->
				<span>{{ item.meta.title }}</span>
			</template>
			<!-- 若路由中未定义菜单项icon，则仅展示名称--（我的仅一级菜单有图标） -->
			<template #title v-else>{{ item.meta.title }}</template>

			<!-- 递归遍历-自己调用自己（核心代码） -->
			<sidebarItem :routerList="( item.children as CustomRouteRecordRaw[])" />
		</el-sub-menu>

		<!-- 当前菜单项无子菜单 -->
		<el-menu-item :index="item.path" v-else>
			<!-- 与上面注释大致相同，不多做额外注释 -->
			<template v-if="item.meta.icon">
				<el-icon><svg-icon :icon-class="item.meta.icon" /></el-icon>
				<span>{{ item.meta.title }}</span>
			</template>
			<template v-else>
				{{ item.meta.title }}
			</template>
		</el-menu-item>
	</template>
</template>

