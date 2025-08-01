<script setup>
  useHead({
    meta: [
      {
        name: 'description',
        content: 'Httpzen is a powerful HTTP client for building and testing APIs, designed to streamline your web development workflow.'
      },
      {
        property: 'og:title',
        content: 'Httpzen - HTTP Client for APIs'
      },
      {
        property: 'og:description',
        content: 'Httpzen is a powerful HTTP client for building and testing APIs, designed to streamline your web development workflow.'
      },
      {
        property: 'og:image',
        content: ''
      },
      {
        property: 'og:type',
        content: 'website'
      },
      {
        property: 'og:url',
        content: 'https://httpzen.diogopereira.site'
      },
      {
        name: 'keywords',
        content: 'httpzen, http, api, rest, web development, software development, programming'
      },
      {
        name: 'charset',
        content: 'UTF-8'
      },
      { name: 'apple-mobile-web-app-title', content: 'Httpzen' }
    ],
    link: [
      { rel: 'icon', type: 'image/png', href: '/favicons/favicon-96x96.png', sizes: '96x96' },
      { rel: 'icon', type: 'image/svg+xml', href: '/favicons/favicon.svg' },
      { rel: 'shortcut icon', href: '/favicons/favicon.ico' },
      { rel: 'apple-touch-icon', sizes: '180x180', href: '/favicons/apple-touch-icon.png' },
      { rel: 'manifest', href: '/favicons/site.webmanifest' }
    ]
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
  }
</style>