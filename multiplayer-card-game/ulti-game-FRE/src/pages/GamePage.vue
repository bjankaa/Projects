<template>
    <waiting-room v-if="!gameStarted"></waiting-room>
    <interruption-message :show="showInterruption" />
    <template v-if="gameStarted">
        <button class="exit-button" @click="showExitConfirmation">✕</button>
        
        <exit-dialog 
            :show="showExitDialog" 
            @confirm="confirmExit" 
            @cancel="showExitDialog = false"
        />
        
        <bidding-phase v-if="gamePhase === 'bidding'" />
        <talon-exchange v-if="gamePhase === 'talon_exchange'" />
        <playing-area v-if="gamePhase === 'playing'" />
        <game-result v-if="gamePhase === 'ended'" />
    </template>
</template>

<script>
import WaitingRoom from "../components/game/WaitingRoom.vue";
import InterruptionMessage from "../components/game/InterruptionMessage.vue";
import ExitDialog from "../components/game/ExitDialog.vue";
import PlayingArea from "../components/game/PlayingArea.vue";
import BiddingPhase from "../components/game/BiddingPhase.vue";
import TalonExchange from "../components/game/TalonExchange.vue";
import GameResult from "../components/game/GameResult.vue";

export default {
    components: {
        WaitingRoom,
        InterruptionMessage,
        ExitDialog,
        PlayingArea,
        BiddingPhase,
        TalonExchange,
        GameResult
    },
    data() {
        return {
            showExitDialog: false,
            showInterruption: false,
        };
    },
    computed: {
        gameid() {
            return this.$store.getters["ws/gameID"];
        },
        gameStarted(){
            return this.$store.getters["ws/gameStarted"];
        },
        gamePhase(){
            return this.$store.getters["ws/gamePhase"];
        }
    },
    watch: {
        gameStarted(newVal, oldVal) {
            if (oldVal === true && newVal === false) {
                this.showInterruption = true;
                setTimeout(() => {
                    this.showInterruption = false;
                    this.$router.push('/game');
                }, 3000);
            }
        },
        gameid(newVal, oldVal) {
            if (oldVal !== 0 && newVal === 0 && this.gamePhase === 'waiting') {
                this.$router.push('/game');
            }
        }
    },

    methods: {
        showExitConfirmation() {
            this.showExitDialog = true;
        },

        confirmExit() {
            this.showExitDialog = false;
            
            const event = {
                type: "game_exit",
                id: this.gameid,
                payload: {},
            };

            this.$store.dispatch("ws/sendEvent", event);
        },

    },
};
</script>

<style scoped>
.exit-button {
    position: fixed;
    top: calc(5rem + 20px);
    right: 20px;
    padding: 0.5rem;
    width: 40px;
    height: 40px;
    background-color: #dc3545;
    color: white;
    border: none;
    border-radius: 50%;
    font-size: 1.5rem;
    font-weight: bold;
    cursor: pointer;
    z-index: 50;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 2px 8px rgba(220, 53, 69, 0.3);
    transition: all 0.3s ease;
    line-height: 1;
}

@media (max-width: 768px) {
    .exit-button {
        top: calc(5rem + 10px);
        right: 10px;
        width: 35px;
        height: 35px;
        font-size: 1.2rem;
    }
}

.exit-button:hover {
    background-color: #c82333;
    transform: scale(1.1);
    box-shadow: 0 4px 12px rgba(220, 53, 69, 0.4);
}
</style>

