<template>
  <section id="features" class="component--features">
    <h2 class="section-title">{{ $t('features.title') }}</h2>
    <p class="section-subtitle">{{ $t('features.subtitle') }}</p>
    <div class="grid">
      <div class="card">
        <div class="icon"><Icon name="mdi:rocket-launch-outline" size="22" /></div>
        <h3>{{ $t('features.fast') }}</h3>
        <p>{{ $t('features.fast-desc') }}</p>
      </div>
      <div class="card">
        <div class="icon"><Icon name="mdi:chart-timeline-variant" size="22" /></div>
        <h3>{{ $t('features.bench') }}</h3>
        <p>{{ $t('features.bench-desc') }}</p>
      </div>
      <div class="card">
        <div class="icon"><Icon name="mdi:security" size="22" /></div>
        <h3>{{ $t('features.inspect') }}</h3>
        <p>{{ $t('features.inspect-desc') }}</p>
      </div>
      <div class="card">
        <div class="icon"><Icon name="mdi:ip-network" size="22" /></div>
        <h3>{{ $t('features.iplookup') }}</h3>
        <p>{{ $t('features.iplookup-desc') }}</p>
      </div>
      <div class="card">
        <div class="icon"><Icon name="mdi:keyboard" size="22" /></div>
        <h3>{{ $t('features.shortcuts') }}</h3>
        <p>{{ $t('features.shortcuts-desc') }}</p>
      </div>
      <div class="card">
        <div class="icon"><Icon name="mdi:code-json" size="22" /></div>
        <h3>{{ $t('features.pretty') }}</h3>
        <p>{{ $t('features.pretty-desc') }}</p>
      </div>
    </div>
  </section>
</template>

<style lang="scss" scoped>
  .component--features {
    position: relative;
    width: 100%;
    max-width: 1400px;
    padding: 0 1rem;
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
    .grid {
      display: grid;
      grid-template-columns: repeat(3, 1fr);
      gap: 1rem;
      @media screen and (max-width: 1100px) {
        grid-template-columns: repeat(2, 1fr);
      }
      @media screen and (max-width: 700px) {
        grid-template-columns: 1fr;
      }
    }
    .card {
      background: $surface;
      border: 1px solid $border-color;
      border-radius: $rounded;
      padding: 1.25rem;
      color: $on-surface;
      .icon { color: $primary; margin-bottom: .5rem; }
      h3 { color: $on-primary; font-size: 1.1rem; margin-bottom: .25rem; }
      p { color: rgba($on-surface, .9); font-size: .95rem; }
    }
  }
</style>