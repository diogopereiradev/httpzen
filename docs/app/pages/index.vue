<script setup lang="ts">
</script>

<template>
  <div class="component--home">
    <div class="wrapper">
      <div class="header">
        <HomeHeader />
      </div>
      <div class="interactive-tui">
        <InteractiveTUI />
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
  .component--home {
    display: flex;
    flex-direction: column;
    width: 100%;
    position: relative;
    .header {
      position: relative;
      z-index: 1;
    }
    .interactive-tui {
      position: relative;
      padding: 0 1rem;
      margin-top: -10px;
      z-index: 2;
    }
  }
</style>