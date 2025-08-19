<script lang="ts" setup>
  const props = defineProps<{
    icon: string;
    title: string;
    desc: string;
    link: string;
    linkLabel: string;
    hint: string;
    code: string;
    codeToCopy: string;
    alt?: boolean;
  }>();
  const { copyToClipboard } = useUtils();
</script>

<template>
  <div class="card">
    <div class="card-head">
      <div class="icon"><Icon :name="props.icon" size="20" /></div>
      <h3>{{ props.title }}</h3>
    </div>
    <p class="desc">{{ props.desc }}</p>
    <div class="links">
      <NuxtLink class="releases" :to="props.link" target="_blank">
        <Icon name="mdi:download" size="18" />
        <span>{{ props.linkLabel }}</span>
      </NuxtLink>
    </div>
    <div class="hint">{{ props.hint }}</div>
    <div class="code-block" :class="{ alt: props.alt }">
      <span>$ {{ props.code }}</span>
      <button class="copy" @click="copyToClipboard(props.codeToCopy)">
        <Icon name="mdi:content-copy" size="18" />
      </button>
    </div>
  </div>
</template>

<style lang="scss" scoped>
  .card {
    background: $surface;
    border: 1px solid $border-color;
    border-radius: $rounded;
    padding: 1.25rem;
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
  }
</style>