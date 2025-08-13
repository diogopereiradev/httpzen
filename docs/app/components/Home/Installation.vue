<script lang="ts" setup>
  const { copyToClipboard } = useUtils();
</script>

<template>
  <section id="install" class="component--installation">
    <div class="content">
      <div class="left">
        <h2 class="section-title">{{ $t('install.title') }}</h2>
        <p class="section-subtitle">{{ $t('install.subtitle') }}</p>
        <div class="cards">
          <div class="card">
            <div class="card-head">
              <div class="icon"><Icon name="mdi:rocket-launch-outline" size="20" /></div>
              <h3>{{ $t('install.binary.title') }}</h3>
            </div>
            <p class="desc">{{ $t('install.binary.desc') }}</p>
            <div class="links">
              <NuxtLink class="releases" to="https://github.com/diogopereiradev/httpzen/releases" target="_blank">
                <Icon name="mdi:download" size="18" />
                <span>{{ $t('install.open-releases') }}</span>
              </NuxtLink>
            </div>
            <div class="hint">{{ $t('install.after-download') }}</div>
            <div class="code-block">
              <span>$ chmod +x ./httpzen && mv ./httpzen /usr/local/bin</span>
              <button class="copy" @click="copyToClipboard('chmod +x ./httpzen && mv ./httpzen /usr/local/bin')">
                <Icon name="mdi:content-copy" size="18" />
              </button>
            </div>
          </div>
          <div class="card">
            <div class="card-head">
              <div class="icon"><Icon name="mdi:package-variant-closed" size="20" /></div>
              <h3>{{ $t('install.deb.title') }}</h3>
            </div>
            <p class="desc">{{ $t('install.deb.desc') }}</p>
            <div class="links">
              <NuxtLink class="releases" to="https://github.com/diogopereiradev/httpzen/releases" target="_blank">
                <Icon name="mdi:download" size="18" />
                <span>{{ $t('install.open-releases') }}</span>
              </NuxtLink>
            </div>
            <div class="hint">{{ $t('install.after-download') }}</div>
            <div class="code-block alt">
              <span>$ apt install ./httpzen.deb</span>
              <button class="copy" @click="copyToClipboard('sudo apt install ./httpzen.deb')">
                <Icon name="mdi:content-copy" size="18" />
              </button>
            </div>
          </div>
          <div class="card">
            <div class="card-head">
              <div class="icon"><Icon name="mdi:package-variant" size="20" /></div>
              <h3>{{ $t('install.rpm.title') }}</h3>
            </div>
            <p class="desc">{{ $t('install.rpm.desc') }}</p>
            <div class="links">
              <NuxtLink class="releases" to="https://github.com/diogopereiradev/httpzen/releases" target="_blank">
                <Icon name="mdi:download" size="18" />
                <span>{{ $t('install.open-releases') }}</span>
              </NuxtLink>
            </div>
            <div class="hint">{{ $t('install.after-download') }}</div>
            <div class="code-block">
              <span>$ dnf install ./httpzen.rpm</span>
              <button class="copy" @click="copyToClipboard('sudo dnf install ./httpzen.rpm')">
                <Icon name="mdi:content-copy" size="18" />
              </button>
            </div>
          </div>
          <div class="card">
            <div class="card-head">
              <div class="icon"><Icon name="mdi:usb-flash-drive" size="20" /></div>
              <h3>{{ $t('install.flatpak.title') }}</h3>
            </div>
            <p class="desc">{{ $t('install.flatpak.desc') }}</p>
            <div class="links">
              <NuxtLink class="releases" to="https://github.com/diogopereiradev/httpzen/releases" target="_blank">
                <Icon name="mdi:download" size="18" />
                <span>{{ $t('install.open-releases') }}</span>
              </NuxtLink>
            </div>
            <div class="hint">{{ $t('install.after-download') }}</div>
            <div class="code-block alt">
              <span>$ flatpak install ./httpzen.flatpak</span>
              <button class="copy" @click="copyToClipboard('flatpak install ./httpzen.flatpak')">
                <Icon name="mdi:content-copy" size="18" />
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<style lang="scss" scoped>
  .component--installation {
    width: 100%;
    position: relative;
    padding: 0 1rem;
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
    .cards {
      display: grid;
      grid-template-columns: repeat(2, 1fr);
      gap: 1rem;
      @media screen and (max-width: 900px) {
        grid-template-columns: 1fr;
      }
    }
    .card {
      background: $surface;
      border: 1px solid $border-color;
      border-radius: $rounded;
      padding: 1.25rem;
    }
    .card-head { 
      display: flex;
      align-items: center;
      gap: .5rem;
      margin-bottom: .25rem;
      h3 { color: $on-primary; font-size: 1.05rem; }
      .icon { color: $primary; }
    }
    .desc { 
      color: rgba($on-surface, .9);
      font-size: .95rem;
      margin-bottom: .75rem;
    }
    .links { 
      margin-bottom: .5rem;
      .releases { 
        display: inline-flex;
        gap: .5rem;
        align-items: center;
        text-decoration: none;
        color: $primary;
        border: 1px solid $primary;
        padding: .4rem .75rem;
        border-radius: 10px;
        transition: .2s;
        &:hover { background: rgba($primary, .1); }
      }
    }
    .hint { 
      color: rgba($on-surface, .8);
      font-size: .85rem;
      margin: 1rem 0 .5rem; 
    }
    .code-block {
      background: $surface-2;
      border: 1px solid $border-color;
      border-radius: $rounded;
      padding: .85rem 1rem;
      color: $on-surface;
      display: flex;
      align-items: center;
      justify-content: space-between;
      margin-bottom: .5rem;
      &.alt { background: $surface; }
      span { font-family: "JetBrains Mono", monospace; font-size: .9rem; }
      .copy {
        border: none;
        background: transparent;
        color: $on-surface;
        cursor: pointer;
        transition: .2s;
        &:hover { color: $primary; }
      }
    }
    .checks {
      list-style: none;
      display: flex; flex-direction: column; gap: .5rem;
      li { 
        display: flex;
        align-items: center;
        gap: .5rem;
        color: rgba($on-primary, .9);
      }
    }
  }
</style>