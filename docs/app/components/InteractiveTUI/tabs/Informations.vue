<script lang="ts" setup>
  import { mockRequest } from '../mocks';

  const request = mockRequest();

  function parseRequestExecutionTime(time: number) {
    if (time < 300) { return { class: 'fast-response', label: 'fast' }; }
    return { class: 'slow-response', label: 'slow' };
  }
</script>

<template>
  <div class="component--tab">
    <p class="http-version">HTTP/{{ request.httpVersion }} GET {{ request.statusCode }}</p>
    <div class="request-info">
      <p class="info">URL: <span class="plain">{{ request.url }}</span></p>
      <p class="info">
        Response Time: 
        <span :class="parseRequestExecutionTime(request.responseTime).class">
          {{ request.responseTime }} ms ({{ parseRequestExecutionTime(request.responseTime).label }})
        </span>
      </p>
      <p class="info">Response Size: <span class="plain">{{ JSON.stringify(request.response).length }} bytes</span></p>
      <p class="info">Request Body: <span class="unavailable">No request body available.</span></p>
    </div>
  </div>
</template>

<style lang="scss" scoped>
  .component--tab {
    .http-version {
      font-size: .875rem;
      color: #5e6069;
    }
    .request-info {
      display: flex;
      flex-direction: column;
      margin-top: 1rem;
      .info {
        font-size: .875rem;
        color: var(--info);
        span {
          &.plain {
            color: var(--on-surface);
          }
          &.fast-response {
            color: #13ce66;
          }
          &.slow-response {
            color: #ff605c;
          }
          &.unavailable {
            color: #ffb366;
          }
        }
      }
    }
  }
</style>