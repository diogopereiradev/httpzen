<script lang="ts" setup>
  import { mockRequest } from '../mocks';

  const request = mockRequest();
</script>

<template>
  <div class="component--tab">
    <p class="info">Protocol: <span class="plain">{{ request.network.protocol }}</span></p>
    <p class="info">IP Address: <span class="plain">{{ request.network.ip }}</span></p>
    <p class="info">Country: <span class="plain">{{ request.network.country }}</span></p>
    <p class="info">Hostname: <span class="plain">{{ request.network.hostname }}</span></p>
    <p class="info">City: <span class="plain">{{ request.network.city }}</span></p>
    <p class="info">ASN: <span class="plain">{{ request.network.asn }}</span></p>
    <p class="info">ISP: <span class="plain">{{ request.network.isp }}</span></p>
    <p class="info">Coordinates: <span class="plain">{{ request.network.coordinates }}</span></p>
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