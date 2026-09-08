<template>
    <div class="border">
        <router-link to="/game"><h3>Game Hub</h3></router-link>
    </div>
    <div class="border">
        <router-link to="/profile"><h3>Profile</h3></router-link>
    </div>
     <div class="border">
        <router-link to="/settings"><h3>Settings</h3></router-link>
    </div>
</template>

<script>


</script>

<style scoped>
a {
    text-decoration: none;
    display: flex;
    justify-content: center;
    align-items: center;
    font-size: large;
    color: var(--primary-color);
    padding: 1rem 2rem;
    border-radius: 2rem;
    transition: all 0.3s ease;
}

h2{
    display: flex;
    justify-content: center;
    align-items: center;
}

a:hover {
    color: var(--text-white);
    background-color: var(--primary-color);
}

.border {
    border-radius: 2rem;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.26);
    padding: 0;
    margin: 2rem auto;
    max-width: 15rem;
    width: 15rem;
}

h3 {
    margin: 0.5rem 0;
}
</style>