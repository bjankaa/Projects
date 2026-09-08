<template>
  <header>
    <nav @click="redirectToFrontpage">
      <h1>ULTI</h1>
      <div class="card-suits">
        <div class="suit-row">
          <img src="../../assets/zold.png" alt="zöld" class="suit-icon">
          <img src="../../assets/makk.png" alt="makk" class="suit-icon">
        </div>
        <div class="suit-row">
          <img src="../../assets/tok.png" alt="tök" class="suit-icon">
          <img src="../../assets/piros.png" alt="piros" class="suit-icon">
        </div>
      </div>
    </nav>
    <ul>
      <li v-if="!isLoggedIn"><router-link to="/auth">Log In</router-link></li>
      <li v-else><button @click="logOut">Log Out</button></li>
    </ul>
  </header>
</template>

<script>

export default {
  computed: {
    isLoggedIn() {
      return this.$store.getters["auth/loggedIn"];
    },
    isInGame() {
      return this.$store.getters["ws/gameStarted"];
    }
  },
  methods: {
    logOut() {
      // If in a game, exit game and close WebSocket
      if (this.isInGame) {
        this.$store.dispatch("ws/exitGameAndLogout");
      }
      
      this.$store.dispatch("auth/logout");
      this.redirectToFrontpage();
    },
    redirectToFrontpage(){
      this.$router.push("/frontpage");
    }
  }
}


</script>


<style scoped>
header {
  background-color: var(--primary-color);
  color: var(--text-white);
  width: 100%;
  height: var(--header-height);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 2rem;
}

header nav {
  display: flex;
  align-items: center;
  gap: 1rem;
  cursor: pointer;
}

h1 {
  font-size: 2.5rem;
  font-weight: 800;
  letter-spacing: 0.15em;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.3);
  font-family: 'Georgia', serif;
}

.card-suits {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.suit-row {
  display: flex;
  gap: 0.25rem;
}

.suit-icon {
  width: 25px;
  height: 25px;
  object-fit: contain;
}

header ul {
  display: flex;
  align-items: center;
  margin: 0;
  padding: 0;
}

li {
  list-style: none;
  display: flex;
}

header a,
header button {
  color: var(--primary-color);
  background: var(--text-white);
  padding: 0.5rem 2rem;
  border: 2px solid transparent;
  text-decoration: none;
  display: inline-block;
  font-size: 1rem;
  font-weight: bold;
  font-family: system-ui, Avenir, Helvetica, Arial, sans-serif;
  border-radius: 1.5rem;
  white-space: nowrap;
  cursor: pointer;
  appearance: none;
}

a:hover,
button:hover {
  background-color: #f0f0f0;
  border: 2px solid var(--primary-hover);
}
</style>