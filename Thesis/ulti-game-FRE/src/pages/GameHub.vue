<template>
    <h1>Game Hub</h1>
    <div class="button-container" v-if="isConnected">
        <div class="border">
            <button @click="gameStart"><h3>Let's Game</h3></button>
        </div>
        <div class="border">
            <button @click="disconnect"><h3>Exit</h3></button>
        </div>
    </div>
    <div v-else class="error-container">
        <div class="error-message">
            <h2>Couldn't connect to server</h2>
            <p>Please check your connection and try again</p>
            <button class="refresh-button" @click="reconnect">
                Refresh <span class="refresh-icon">⟳</span>
            </button>
        </div>
    </div>

    <div class="rules-container">
        <div class="rules-header" @click="toggleRules">
            <h2>Game Rules</h2>
            <span class="toggle-icon">{{ showRules ? '▲' : '▼' }}</span>
        </div>
        <div v-if="showRules" class="rules-content">
            <p>Game rules will be added here...</p>
            <ul>
                <li>Rule 1: ...</li>
                <li>Rule 2: ...</li>
                <li>Rule 3: ...</li>
            </ul>
        </div>
    </div>
</template>


<script>

export default {
    data() {
        return {
            gameId: 0,
            showRules: false,
        };
    },
    computed: {
        isConnected(){
            return this.$store.getters["ws/isConnected"];
        },
        storedGameId(){
            return this.$store.getters["ws/gameID"];
        }
    },
    watch: {
        isConnected(val) {
            if (val) {
                console.log("WebSocket connected!");
            } else {
                console.log("WebSocket disconnected!");
                this.gameExist = false;
                this.gameId = 0;
            }
        },
        storedGameId (newId){
            console.log("GameID from store changed:", newId);
            this.gameId = newId;

            if (newId && newId !== 0) {
                this.gameExist = true;

                this.$router.push(`/game/${newId}`);
            } else {
                this.gameExist = false;
            }
        }
    },
    mounted(){
		console.log("mounted");
        this.websocketRequest();
    },
    methods: {
        async websocketRequest(){
            await this.$store.dispatch("ws/connect");
        },
        gameStart(){
            if(!this.isConnected){
                console.log("Not connected to a websocket")
                return;

            };

            const event = {
                type: "game_init",
                id: 0,
                payload: {}
            };

            this.$store.dispatch("ws/sendEvent", event);
            console.log("Game init was sent");
            
        },
        disconnect(){
            this.gameId = 0;
            this.$router.push("/frontpage");
			this.$store.dispatch("ws/closeSocket");
			console.log('Not connewcted anymore');
			
		},
        toggleRules() {
            this.showRules = !this.showRules;
        },
        reconnect() {
            this.websocketRequest();
        }
    }
};

</script>

<style scoped>
h1 {
    text-align: center;
    margin: 2rem auto 1rem;
    padding: 0 1rem;
}

.button-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0;
    margin: 1rem auto 2rem;
}

.border {
    border-radius: 2rem;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.26);
    padding: 0;
    margin: 1rem auto;
    max-width: 15rem;
    width: 15rem;
}

button {
    width: 100%;
    text-decoration: none;
    display: flex;
    justify-content: center;
    align-items: center;
    font-size: large;
    color: var(--primary-color);
    padding: 1rem 2rem;
    border-radius: 2rem;
    transition: all 0.3s ease;
    background: transparent;
    border: none;
    cursor: pointer;
}

button:hover {
    color: var(--text-white);
    background-color: var(--primary-color);
}

h3 {
    margin: 0.5rem 0;
}

h2 {
    text-align: center;
    margin: 2rem auto;
    padding: 0 1rem;
}

.error-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1rem;
    margin: 1rem auto 2rem;
    padding: 0 1rem;
}

.error-message {
    background-color: #f8d7da;
    border: 2px solid #f5c6cb;
    color: #721c24;
    padding: 1.5rem;
    border-radius: 1rem;
    max-width: 400px;
    width: 90%;
    text-align: center;
}

.error-message h2 {
    margin: 0 0 0.5rem 0;
    font-size: 1.25rem;
}

.error-message p {
    margin: 0 0 1rem 0;
    font-size: 0.95rem;
}

.refresh-button {
    padding: 0.4rem 0.8rem;
    border-radius: 2rem;
    background-color: var(--primary-color);
    color: var(--text-white);
    border: 2px solid var(--primary-color);
    cursor: pointer;
    transition: all 0.3s ease;
    display: inline-flex;
    align-items: center;
    gap: 0.4rem;
    font-size: 0.9rem;
    width: fit-content;
}

.refresh-button:hover {
    background-color: var(--text-white);
    color: var(--primary-color);
}

.refresh-icon {
    font-size: 1.25rem;
    font-weight: bold;
    line-height: 1;
}

.rules-container {
    max-width: 600px;
    margin: 2rem auto;
    padding: 0 1rem;
}

.rules-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1rem 1.5rem;
    background-color: var(--primary-color);
    color: var(--text-white);
    border-radius: 1rem;
    cursor: pointer;
    transition: all 0.3s ease;
}

.rules-header:hover {
    opacity: 0.9;
}

.rules-header h2 {
    margin: 0;
    padding: 0;
    font-size: 1.25rem;
}

.toggle-icon {
    font-size: 1.5rem;
    transition: transform 0.3s ease;
}

.rules-content {
    padding: 1.5rem;
    background-color: #f5f5f5;
    border-radius: 0 0 1rem 1rem;
    margin-top: -0.5rem;
    animation: slideDown 0.3s ease;
}

@keyframes slideDown {
    from {
        opacity: 0;
        transform: translateY(-10px);
    }
    to {
        opacity: 1;
        transform: translateY(0);
    }
}

.rules-content p {
    margin: 0 0 1rem 0;
}

.rules-content ul {
    margin: 0;
    padding-left: 1.5rem;
}

.rules-content li {
    margin: 0.5rem 0;
}
</style>