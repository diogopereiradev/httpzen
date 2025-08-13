<script lang="ts" setup>
  import pkg from '@/../package.json';

  const utils = useUtils();
  const state = reactive({
    copied: false,
  });

  const copyToClipboard = (text: string) => {
    utils.copyToClipboard(text);
    state.copied = true;
    setTimeout(() => {
      state.copied = false;
    }, 2000);
  };
</script>

<template>
  <div class="component--home--header">
    <nav class="navbar">
      <div class="left">
        <h3>httpzen</h3>
        <div class="version"><p>v{{ pkg.version.split('.')[0] }}.{{ pkg.version.split('.')[1] }}</p></div>
      </div>
      <div class="right">
        <NuxtLink to="#features" class="nav-link">{{ $t('nav.features') }}</NuxtLink>
        <NuxtLink to="#install" class="nav-link">{{ $t('nav.install') }}</NuxtLink>
        <NuxtLink to="#docs" class="nav-link">{{ $t('nav.docs') }}</NuxtLink>
        <a href="https://github.com/diogopereiradev/httpzen" target="_blank" title="GitHub">
          <Icon name="f7:logo-github" size="24" />
        </a>
      </div>
    </nav>
    <div class="glow-beam-effect" />
    <div class="glow-lamp-effect" />
    <div class="infos">
      <a class="version-tag" href="https://github.com/diogopereiradev/httpzen/releases" target="_blank" title="Version">
        <span>v{{ pkg.version }} {{ $t('is-out') }}</span>
        <Icon name="material-symbols:arrow-right-alt-rounded" size="18" />
      </a>
      <h1 class="title">{{ $t('title') }}</h1>
      <p class="description">{{ $t('description') }}</p>
      <div class="actions">
        <div class="example-command">
          <Icon name="octicon:command-palette-16" size="20" />
          <span>httpzen GET https://google.com</span>
          <button class="copy-btn" :disabled="state.copied" @click="copyToClipboard('httpzen GET https://google.com')">
            <Icon v-if="!state.copied" name="mdi:content-copy" size="20" />
            <Icon v-else name="material-symbols:check-rounded" size="20" />
          </button>
        </div>
          <NuxtLink to="#install" class="install-button">
          <span>{{ $t('installation') }}</span>
          <Icon name="material-symbols:arrow-right-alt-rounded" size="18" />
        </NuxtLink>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
  .component--home--header {
    display: flex;
    flex-direction: column;
    position: relative;
    width: 100%;
    min-height: 600px;
    z-index: 1;
    overflow: hidden;
    .navbar {
      width: 100%;
      height: 80px;
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding: 0 1rem;
      .left {
        display: flex;
        align-items: center;
        gap: .5rem;
        h3 {
          font-size: 1.25rem;
          color: $on-background;
        }
        .version {
          background-color: rgba($primary, .15);
          border: 1px solid $primary;
          border-radius: $rounded;
          padding: .125rem .5rem;
          p {
            color: $primary;
            font-size: .725rem;
          }
        }
      }
      .right {
        display: flex;
        align-items: center;
        gap: .875rem;
        .nav-link {
          color: $on-background;
          text-decoration: none;
          font-size: .9rem;
          padding: .25rem .5rem;
          border-radius: 8px;
          transition: .2s;
          margin-top: -8px;
          &:hover {
            color: $primary;
            background: rgba($primary, .08);
          }
        }
        .theme-toggle {
          border: none;
          background: transparent;
          color: $on-background;
          cursor: pointer;
          transition: .2s;
          &:hover {
            color: $primary;
          }
        }
        a {
          border: none;
          background: transparent;
          color: $on-background;
          cursor: pointer;
          transition: .2s;
          &:hover {
            color: $primary;
          }
        }
      }
    }
    .glow-beam-effect {
      position: absolute;
      top: -100px;
      left: 50%;
      width: 150px;
      height: 900px;
      background-color: $primary;
      transform: translateX(-50%);
      border-radius: 50%;
      filter: blur(90px);
      opacity: .3;
    }
    .glow-lamp-effect {
      position: absolute;
      top: -200px;
      left: 50%;
      width: 0;
      height: 0;
      border-left: 600px solid transparent;
      border-right: 600px solid transparent;
      border-top: 1000px solid $primary;
      transform: translateX(-50%) rotate(180deg);
      opacity: .015;
      pointer-events: none;
      @media screen and (max-width: 1100px) {
        border-left: 330px solid transparent;
        border-right: 330px solid transparent;
        border-top: 800px solid $primary;
      }
      &::after {
        content: '';
        position: absolute;
        left: 50%;
        bottom: 800px;
        transform: translateX(-50%);
        width: 1200px;
        height: 300px;
        background-color: $background;
        z-index: 1;
        opacity: 1;
      }
    }
    .infos {
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 1.25rem;
      z-index: 2;
      margin: auto 0;
      .version-tag {
        display: flex;
        align-items: center;
        justify-content: center;
        gap: .5rem;
        background-color: rgba($primary, .2);
        padding: 8px 1rem;
        border-radius: 50px;
        text-decoration: none;
        color: $on-primary;
        transition: .2s;
        border: 1px solid $primary-2;
        &:hover {
          background-color: rgba($primary-2, .3);
          :deep(.iconify) {
            left: 4px;
          }
        }
        :deep(.iconify) {
          position: relative;
          left: 0px;
          transition: .2s;
        }
        span {
          color: $primary;
          font-size: .875rem;
          font-weight: 600;
        }
      }
      .title {
        margin-top: 1rem;
        font-size: 2.5rem;
        color: $on-primary;
      }
      .description {
        max-width: 700px;
        width: 100%;
        text-align: center;
        font-size: 1rem;
        color: rgba($on-primary, .9);
      }
      .actions {
        display: flex;
        flex-direction: row-reverse;
        align-items: center;
        gap: 1rem;
        .example-command {
          display: flex;
          align-items: center;
          gap: .5rem;
          background-color: $surface-o3;
          padding: 0 2rem;
          height: 3.5rem;
          color: $on-surface;
          border-radius: $rounded;
          box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
          span {
            font-size: .875rem;
            color: rgba($on-surface, .8);
          }
          .copy-btn {
            width: 30px;
            height: 30px;
            display: flex;
            align-items: center;
            justify-content: center;
            color: $on-surface;
            border: none;
            background-color: transparent;
            margin-left: 1.5rem;
            transition: .2s;
            cursor: pointer;
            &:hover {
              color: $primary;
            }
            &:disabled {
              color: $success;
              cursor: not-allowed;
            }
          }
        }
        .install-button {
          display: flex;
          align-items: center;
          justify-content: center;
          gap: .5rem;
          background-color: $primary;
          color: $on-primary;
          border: none;
          border-radius: $rounded;
          padding: 0 2rem;
          height: 3.5rem;
          text-decoration: none;
          cursor: pointer;
          transition: .2s;
          box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
          &:hover {
            background-color: $primary-2;
            :deep(.iconify) {
              left: 4px;
            }
          }
          :deep(.iconify) {
            position: relative;
            left: 0px;
            transition: .2s;
          }
        }
      }
    }
  }
</style>