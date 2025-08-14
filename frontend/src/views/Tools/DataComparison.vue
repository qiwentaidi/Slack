<template>
  <vue-monaco-diff-editor
    theme="vs-dark"
    original="还没实装"
    modified="// the modified code"
    language="javascript"
    :options="OPTIONS"
    @mount="handleMount"
  />
</template>

<script lang="ts" setup>
import { ref, shallowRef } from 'vue'

const OPTIONS = {
  automaticLayout: true,
  formatOnType: true,
  formatOnPaste: true,
  readOnly: false,
}

const diffEditor = shallowRef()
const handleMount = diffEditorInstance => (diffEditor.value = diffEditorInstance)

// get the original value
function getOriginalValue() {
  return diffEditor.value.getOriginalEditor().getValue()
}

// get the modified value
function getModifiedValue() {
  return diffEditor.value.getModifiedEditor().getValue()
}
</script>