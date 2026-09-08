<template>
  <div class="app-wrapper">
    <the-header></the-header>
    <main class="main-content">
      <router-view v-slot="slotProps">
        <component :is="slotProps.Component"></component>
      </router-view>
    </main>
    <the-footer></the-footer>
  </div>
</template>

<script>
import TheHeader from "./components/layout/TheHeader.vue";
import TheFooter from "./components/layout/TheFooter.vue";

export default {
  components: {
    TheHeader,
    TheFooter
  }
}

</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

:root {
  font-family: system-ui, Avenir, Helvetica, Arial, sans-serif;
  --primary-color: #1f4e1b;
  --primary-hover: #02220a;
  --text-white: white;
  --header-height: 5rem;
}

html, body {
  margin: 0;
  height: 100%;
}

#app {
  min-height: 100%;
  display: flex;
  flex-direction: column;
}

.app-wrapper {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

.main-content {
  flex: 1;
}
</style>