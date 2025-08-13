<script lang="ts" setup>
  const Tabs = {
    RESPONSE: 0,
    INFOS: 1,
    NETWORK: 2,
    REQUEST_HEADERS: 3,
    RESPONSE_HEADERS: 4,
  };

  const state = reactive({
    activeTab: Tabs.RESPONSE,
  });
</script>

<template>
  <div class="component--interactive-tui">
    <div class="interactive-tui">
      <div class="header">
        <div class="left">
          <div class="dots">
            <span class="dot" />
            <span class="dot" />
            <span class="dot" />
          </div>
        </div>
        <div class="right">
          <button title="Minimize"><Icon name="ic:round-minimize" size="16" /></button>
          <button title="Maximize"><Icon name="solar:maximize-bold" size="16" /></button>
          <button title="Close"><Icon name="ic:round-close" size="18" /></button>
        </div>
      </div>
      <NuxtImg class="ascii-logo" src="/images/ascii-logo.png" format="webp" />
      <div class="tabs">
        <button
          class="tab"
          :class="{ active: state.activeTab === Tabs.RESPONSE }"
          @click="state.activeTab = Tabs.RESPONSE"
        >
          <span class="label">Response</span>
          <span class="icon"><Icon name="mdi:file-document" size="20" /></span>
        </button>
        <button
          class="tab"
          :class="{ active: state.activeTab === Tabs.INFOS }"
          @click="state.activeTab = Tabs.INFOS"
        >
          <span class="label">Request Infos</span>
          <span class="icon"><Icon name="gravity-ui:server" size="20" /></span>
        </button>
        <button
          class="tab"
          :class="{ active: state.activeTab === Tabs.NETWORK }"
          @click="state.activeTab = Tabs.NETWORK"
        >
          <span class="label">Network</span>
          <span class="icon"><Icon name="ph:network-fill" size="20" /></span>
        </button>
        <button
          class="tab"
          :class="{ active: state.activeTab === Tabs.REQUEST_HEADERS }"
          @click="state.activeTab = Tabs.REQUEST_HEADERS"
        >
          <span class="label">Request Headers</span>
          <span class="icon"><Icon name="gravity-ui:layout-header-cells" size="20" /></span>
        </button>
        <button
          class="tab"
          :class="{ active: state.activeTab === Tabs.RESPONSE_HEADERS }"
          @click="state.activeTab = Tabs.RESPONSE_HEADERS"
        >
          <span class="label">Response Headers</span>
          <span class="icon"><Icon name="gravity-ui:layout-header-cells-large-thunderbolt" size="20" /></span>
        </button>
      </div>
      <InteractiveTUITabsResponse v-if="state.activeTab === Tabs.RESPONSE" />
      <InteractiveTUITabsInformations v-if="state.activeTab === Tabs.INFOS" />
      <InteractiveTUITabsNetwork v-if="state.activeTab === Tabs.NETWORK" />
      <InteractiveTUITabsRequestHeaders v-if="state.activeTab === Tabs.REQUEST_HEADERS" />
      <InteractiveTUITabsResponseHeaders v-if="state.activeTab === Tabs.RESPONSE_HEADERS" />
      <div class="navigation-infos">
        <p class="info">Use left/right arrows to navigate between tabs, 'q' to quit.</p>
        <p class="info">'c' to copy response, 'b' to benchmark, 'r' to resend request.</p>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
  .component--interactive-tui {
    position: relative;
    display: flex;
    align-items: center;
    justify-content: center;
    .interactive-tui {
      position: relative;
      width: 1000px;
      min-height: 200px;
      padding: 2rem;
      background-color: $surface;
      border-radius: $rounded;
      border: 10px solid rgba($surface-o3, .15);
      box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
      &::before {
        content: '';
        position: absolute;
        top: 1px;
        left: -1%;
        width: 102%;
        height: 10px;
        border-radius: $rounded $rounded 0 0;
        z-index: 2;
      }
      @media screen and (max-width: 1100px) {
        width: 700px;
      }
      @media screen and (max-width: 800px) {
        width: 100%;
      }
      @media screen and (max-width: 500px) {
        padding: 1rem;
      }
      .header {
        width: 100%;
        display: flex;
        align-items: center;
        justify-content: space-between;
        margin-bottom: 1.5rem;
        .left {
          .dots {
            display: flex;
            align-items: center;
            gap: .5rem;
            .dot {
              width: 16px;
              height: 16px;
              border-radius: 50%;
              &:nth-child(1) {
                background-color: #ff605c;
              }
              &:nth-child(2) {
                background-color: #ffbd44;
              }
              &:nth-child(3) {
                background-color: #00ca4e;
              }
            }
          }
        }
        .right {
          display: flex;
          align-items: center;
          button {
            width: 30px;
            height: 30px;
            display: flex;
            align-items: center;
            justify-content: center;
            border: none;
            background-color: transparent;
            color: $on-surface;
            border-radius: $rounded;
            transition: .2s;
            cursor: pointer;
            &:hover {
              background-color: $surface-2
            }
          }
        }
      }
      .ascii-logo {
        max-width: 450px;
        pointer-events: none;
        user-select: none;
        @media screen and (max-width: 800px) {
          max-width: 100%;
          width: 100%;
        }
      }
      .tabs {
        position: relative;
        display: flex;
        align-items: center;
        gap: .5rem;
        margin-top: 1rem;
        margin-bottom: 1.5rem;
        &:before {
          content: '';
          position: absolute;
          left: -2.5%;
          bottom: 0;
          width: 105%;
          height: 1px;
          background-color: $primary;
          z-index: 1;
        }
        .tab {
          padding: .725rem 1rem;
          border-radius: 6px;
          border-bottom-left-radius: 0px;
          border-bottom-right-radius: 0px;
          border: 1px solid $primary;
          background-color: transparent;
          color: $on-surface;
          font-size: .875rem;
          font-weight: 500;
          color: $on-surface;
          cursor: pointer;
          transition: .2s;
          z-index: 2;
          &.active {
            color: $primary;
            border-bottom: 1px solid $surface;
          }
          @media screen and (max-width: 500px) {
            padding: .5rem .875rem;
          }
          @media screen and (max-width: 400px) {
            padding: .425rem .65rem;
          }
          .label {
            @media screen and (max-width: 1100px) {
              display: none;
            }
          }
          .icon {
            display: none;
            @media screen and (max-width: 1100px) {
              display: unset;
            }
          }
        }
      }
      .navigation-infos {
        margin-top: 1rem;
        .info {
          font-size: .875rem;
          color: #5e6069;
          @media screen and (max-width: 500px) {
            font-size: .725rem;
            margin-bottom: .25rem;
          }
        }
      }
    }
  }
</style>