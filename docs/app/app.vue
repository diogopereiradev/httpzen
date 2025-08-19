<script setup>
  useHead({
    htmlAttrs: {
      lang: 'en',
    },
    title: 'HTTPzen - Terminal HTTP Client',
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
        content: 'https://httpzen.diogopereira.site/images/banner.jpg'
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
        name: 'robots',
        content: 'index, follow'
      },
      {
        name: 'keywords',
        content: 'httpzen, http, api, rest, web development, software development, programming'
      },
      {
        name: 'charset',
        content: 'UTF-8'
      },
      { name: 'apple-mobile-web-app-title', content: 'Httpzen' },
      {
        name: 'theme-color',
        content: '#e86d01'
      },
      {
        name: 'twitter:card',
        content: 'summary_large_image'
      },
      {
        name: 'twitter:title',
        content: 'Httpzen - HTTP Client for APIs'
      },
      {
        name: 'twitter:description',
        content: 'Httpzen is a powerful HTTP client for building and testing APIs, designed to streamline your web development workflow.'
      },
      {
        name: 'twitter:image',
        content: 'https://httpzen.diogopereira.site/images/banner.jpg'
      },
      {
        name: 'twitter:url',
        content: 'https://httpzen.diogopereira.site'
      },
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