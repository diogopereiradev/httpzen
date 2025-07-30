<script setup lang="ts">
  const theme = useTheme();
  const { locale, setLocale } = useI18n();
</script>

<template>
  <div>
    <button @click="theme.set(theme.current.value === Themes.DARK ? Themes.LIGHT : Themes.DARK)">change theme</button>
    <button @click="setLocale(locale === 'en' ? 'pt' : 'en')">change language</button>
    <InteractiveTUI />
  </div>
</template>

<style lang="scss" scoped>
  div {
    height: 100vh;
    display: flex;
    flex-direction: column;
    align-items: center;
  }
</style>