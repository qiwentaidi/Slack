<template>
    <div class="context-menu">
      <div class="click-div" @click="copyData(data)">Copy</div>
      <div class="click-div" @click="copySource(data)">Copy Source</div>
      <div class="click-div" @click="openURL(data)">Open Source</div>
    </div>
  </template>
  
  <script lang="ts">
  import { BrowserOpenURL } from '../../wailsjs/runtime/runtime';
  import { Copy } from '../util';
  export default {
    props: {
      data: {
        type: Object,
        required: true,
      },
    },
    methods: {
      copyData(data: any) {
        Copy(data)
      },
      copySource(data: any) {
        Copy(data)
      },
      openURL(data: any) {
        BrowserOpenURL(data)
      },
    },
  };
  </script>