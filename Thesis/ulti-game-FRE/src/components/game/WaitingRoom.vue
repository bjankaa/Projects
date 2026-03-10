<template>
    <base-card v-if="!gameStarted">
        <div class="waiting-room">
            <h2>Waiting for Players...</h2>
            <p class="player-count">{{ playerCount }} / 3 players joined</p>
            <div class="loading-animation">
                <div class="dot"></div>
                <div class="dot"></div>
                <div class="dot"></div>
            </div>
            <p class="waiting-text">Waiting for all players to join....</p>
        </div>
    </base-card>
</template>

<script>
export default {
    computed: {
        playerCount(){
            const playerNames = this.$store.getters["ws/playerNames"];
            return Object.keys(playerNames).length;
        },
        gameStarted(){
            return this.$store.getters["ws/gameStarted"];
        },
    },
};



</script>

<style scoped>
    

.waiting-room {
    text-align: center;
    padding: 2rem;
}

.waiting-room h2 {
    color: #333;
    margin-bottom: 1rem;
}

.player-count {
    font-size: 1.5rem;
    font-weight: bold;
    color: #ffa500;
    margin: 1rem 0;
}

.loading-animation {
    display: flex;
    justify-content: center;
    gap: 0.5rem;
    margin: 2rem 0;
}

.dot {
    width: 12px;
    height: 12px;
    background-color: #ffa500;
    border-radius: 50%;
    animation: bounce 1.4s infinite ease-in-out both;
}

.dot:nth-child(1) {
    animation-delay: -0.32s;
}

.dot:nth-child(2) {
    animation-delay: -0.16s;
}

@keyframes bounce {
    0%, 80%, 100% {
        transform: scale(0);
    }
    40% {
        transform: scale(1);
    }
}

.waiting-text {
    color: #666;
    font-style: italic;
}
</style>