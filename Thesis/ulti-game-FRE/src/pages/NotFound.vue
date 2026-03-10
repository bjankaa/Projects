<template>
  <section>
    <base-card>
      <h2>Page not found</h2>
      <p>
        This page could not be found - maybe check out all our
        <router-link to="/frontpage">Go back to frontpage</router-link>.
      </p>
    </base-card>
  </section>
</template>

<style scoped>
section {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 80vh;
}
</style>