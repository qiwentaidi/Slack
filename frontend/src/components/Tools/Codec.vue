<script lang="ts" setup>
import Timestamp from './transcode/Timestamp.vue';
import Hash from './transcode/Hash.vue';
import Transcode from './transcode/Transcode.vue';
</script>

<template>
    <div style="display: flex;">
        <div>
            <Timestamp></Timestamp>
            <Hash style="margin-top: 10px;"></Hash>
        </div>
        <Transcode style="margin-left: 10px;"></Transcode>
    </div>
</template>