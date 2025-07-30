<script lang="ts" setup>
  import 'vue-json-pretty/lib/styles.css';
  import VueJsonPretty from 'vue-json-pretty';
  import { mockRequest } from '../mocks';

  const request = mockRequest();
</script>

<template>
  <div class="component--tab">
    <vue-json-pretty
      :show-line="false"
      :highlight-selected-node="false"
      :data="request.response"
    />
  </div>
</template>

<style lang="scss" scoped>
  .component--tab {
    :deep(span) {
      color: #f1f1f1;
    }
    :deep(.vjs-tree-node):hover {
      background-color: transparent;
    }
    :deep(.vjs-tree-brackets) {
      color: #f1f1f1;
    }
    :deep(.vjs-value-string) {
      color: #13ce66;
    }
    :deep(.vjs-value-number) {
      color: #d64363;
    }
    :deep(.vjs-value-boolean) {
      color: #973eb3;
    }
    :deep(.vjs-key) {
      color: var(--info);
    }
  }
</style>