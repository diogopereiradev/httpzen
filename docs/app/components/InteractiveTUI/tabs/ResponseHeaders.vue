<script lang="ts" setup>
  import { mockRequest } from '../mocks';
  const request = mockRequest();
</script>

<template>
  <div class="component--tab">
    <p v-for="(header, index) in Object.keys(request.headers)" :key="index" class="info" >
      {{ header }}: <span class="plain">{{ request.headers[header] }}</span>
    </p>
  </div>
</template>

<style lang="scss" scoped>
  .component--tab {
    .info {
      font-size: .875rem;
      color: $primary;
      span {
        &.plain {
          color: $on-surface;
        }
      }
    }
  }
</style>