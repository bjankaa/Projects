<template>
    <div v-if="!isLoggedIn">
        <logged-out></logged-out>
    </div>
    <div v-else class="centered-container">
        <logged-in></logged-in>
    </div>
</template>


<script>

import LoggedIn from '../components/frontpage/LoggedIn.vue';
import LoggedOut from '../components/frontpage/LoggedOut.vue';

export default {
    components: {
        LoggedOut,
        LoggedIn
    },
    computed: {
        isLoggedIn() {
            const login = this.$store.getters["auth/loggedIn"];
            return login;
        }
    }
}

</script>

<style scoped>
.centered-container {
    display: flex;
    flex-direction: column;
    justify-content: center;
    min-height: 70vh;
}
</style>