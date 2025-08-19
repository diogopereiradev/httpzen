<template>
  <section id="install" class="component--installation">
    <div class="content">
      <div class="left">
        <h2 class="section-title">{{ $t('install.title') }}</h2>
        <p class="section-subtitle">{{ $t('install.subtitle') }}</p>
        <div class="cards">
          <HomeInstallCard
            :icon="'mdi:rocket-launch-outline'"
            :title="$t('install.binary.title')"
            :desc="$t('install.binary.desc')"
            link="https://github.com/diogopereiradev/httpzen/releases/latest/download/httpzen"
            :link-label="$t('install.open-releases')"
            :hint="$t('install.after-download')"
            :code="'chmod +x ./httpzen && mv ./httpzen /usr/local/bin'"
            :code-to-copy="'chmod +x ./httpzen && mv ./httpzen /usr/local/bin'"
          />
          <HomeInstallCard
            :icon="'mdi:package-variant-closed'"
            :title="$t('install.deb.title')"
            :desc="$t('install.deb.desc')"
            link="https://github.com/diogopereiradev/httpzen/releases/latest/download/httpzen.deb"
            :link-label="$t('install.open-releases')"
            :hint="$t('install.after-download')"
            :code="'apt install ./httpzen.deb'"
            :code-to-copy="'sudo apt install ./httpzen.deb'"
            :alt="true"
          />
          <HomeInstallCard
            :icon="'mdi:package-variant'"
            :title="$t('install.rpm.title')"
            :desc="$t('install.rpm.desc')"
            link="https://github.com/diogopereiradev/httpzen/releases/latest/download/httpzen.rpm"
            :link-label="$t('install.open-releases')"
            :hint="$t('install.after-download')"
            :code="'dnf install ./httpzen.rpm'"
            :code-to-copy="'sudo dnf install ./httpzen.rpm'"
          />
          <HomeInstallCard
            :icon="'mdi:usb-flash-drive'"
            :title="$t('install.flatpak.title')"
            :desc="$t('install.flatpak.desc')"
            link="https://github.com/diogopereiradev/httpzen/releases/latest/download/httpzen.flatpak"
            :link-label="$t('install.open-releases')"
            :hint="$t('install.after-download')"
            :code="'flatpak install ./httpzen.flatpak'"
            :code-to-copy="'flatpak install ./httpzen.flatpak'"
            :alt="true"
          />
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