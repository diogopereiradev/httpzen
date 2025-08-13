<template>
  <footer class="component--footer">
    <div id="docs" class="docs">
      <h2 class="section-title">{{ $t('docs.title') }}</h2>
      <p class="section-subtitle">{{ $t('docs.subtitle') }}</p>
      <div class="cta">
        <NuxtLink to="https://github.com/diogopereiradev/httpzen#readme" target="_blank" class="cta-btn">
          <Icon name="mdi:book-open-variant" size="20" />
          <span>{{ $t('docs.read') }}</span>
        </NuxtLink>
        <NuxtLink to="https://github.com/diogopereiradev/httpzen" target="_blank" class="cta-btn secondary">
          <Icon name="mdi:github" size="20" />
          <span>GitHub</span>
        </NuxtLink>
      </div>
    </div>
    <div class="copyright-message">
      <p>© {{ new Date().getFullYear() }} httpzen • {{ $t('footer.made') }}</p>
    </div>
  </footer>
</template>

<style lang="scss" scoped>
  .component--footer {
    position: relative;
    padding: 0 1rem;
    width: 100%;
    max-width: 1400px;
    margin: 0 auto;
    .section-title {
      text-align: center;
      font-size: 2rem;
      color: $on-primary;
      margin-bottom: .5rem;
    }
    .section-subtitle {
      text-align: center;
      color: rgba($on-primary, .85);
      margin-bottom: 2rem;
    }
    .docs {
      .cta { 
        display: flex;
        gap: 1rem;
        align-items: center;
        justify-content: center;
        margin-top: 1rem;
      }
      .cta-btn {
        display: inline-flex;
        gap: .5rem;
        align-items: center;
        justify-content: center;
        height: 3rem;
        padding: 0 1.25rem;
        border-radius: $rounded;
        text-decoration: none;
        cursor: pointer;
        background: $primary; 
        color: $on-primary;
        transition: .2s;
        border: 1px solid $primary;
        &:hover { background: $primary-2; }
        &.secondary {
          background: transparent;
          color: $on-primary; 
          border-color: $border-color; 
          &:hover { 
            border-color: $primary;
            color: $primary;
          }
        }
      }
    }
    .copyright-message {
      padding: 2rem 1rem;
      text-align: center;
      color: rgba($on-primary, .75);
    }
  }
</style>