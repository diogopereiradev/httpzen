<template>
  <div class="page--home">
    <div class="wrapper">
      <HomeHeader />
      <InteractiveTUI />
      <div class="content">
        <HomeFeatures />
        <HomeInstallation />
        <HomeFooter />
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
  .page--home {
    display: flex;
    flex-direction: column;
    width: 100%;
    position: relative;
    .content {
      display: flex;
      flex-direction: column;
      gap: 4rem;
      width: 100%;
      margin-top: 4rem;
    }
  }
</style>