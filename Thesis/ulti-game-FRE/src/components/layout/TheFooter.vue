<template>
  <footer>
    <p>&copy; 2025 ULTI by @bjanka</p>
  </footer>
</template>

<script>
export default {
  name: 'TheFooter'
}
</script>

<style scoped>
footer {
  background-color: var(--primary-color);
  color: var(--text-white);
  width: 100%;
  padding: 0.75rem;
  display: flex;
  align-items: center;
  justify-content: center;
}

footer p {
  margin: 0;
  font-size: 0.75rem;
  opacity: 0.9;
}
</style>
