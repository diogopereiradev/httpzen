<script setup>
  const theme = useTheme();

  useHead({
    htmlAttrs: {
      'data-theme': theme.current.value
    }
  });
</script>

<template>
  <div>
    <NuxtLoadingIndicator />
    <div class="main-container">
      <NuxtLayout>
        <NuxtPage />
      </NuxtLayout>
    </div>
  </div>
</template>

<style lang="scss" scoped>
  .main-container {
    width: 100%;
    min-height: 100vh;
    padding: 0 1rem;
  }
</style>