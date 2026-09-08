<template>
    <div class="game-container">
        <!-- Left Player Name -->
        <div class="side-player left-player">
            <h3>{{ getPlayerName(getLeftPlayerIndex()) }}</h3>
        </div>
        
        <!-- Right Player Name -->
        <div class="side-player right-player">
            <h3>{{ getPlayerName(getRightPlayerIndex()) }}</h3>
        </div>
        
        <!-- Round Winner Notification -->
        <div v-if="showRoundWinnerNotification" class="round-winner-notification">
            {{ getPlayerName(lastRoundWinner) }} won the round!
        </div>
        
        <!-- Played Cards Area - Centered -->
        <div class="played-cards-area">
            <h2>Played Cards:</h2>
            <div class="cards-display">
                <!-- Left position -->
                <div class="card-position left">
                    <div v-if="leftCard" class="card-content">
                        <p>Card: {{ leftCard.cardId }}</p>
                    </div>
                </div>
                
                <!-- Center position (your card) -->
                <div class="card-position center">
                    <div v-if="centerCard" class="card-content">
                        <p>Card: {{ centerCard }}</p>
                    </div>
                </div>
                
                <!-- Right position -->
                <div class="card-position right">
                    <div v-if="rightCard" class="card-content">
                        <p>Card: {{ rightCard.cardId }}</p>
                    </div>
                </div>
            </div>
        </div>
        
        <div class="betting-pool-display">
            <h3>Betting Pool: {{ bettingPool }}</h3>
        </div>
        
        <div class="points-display">
            <h3>Your Points: {{ yourPoints }}</h3>
        </div>
        
        <div class="game-color-display" v-if="gameColor">
            <p>{{ gameColor.toUpperCase() }} - {{ getPlayerName(declarerIndex) }}</p>
        </div>
        
        <div class="player-section">
            <p class="current-player-name">{{ Name }} (You)</p>
            <div class="hand-background">
                <ul class="hand">
                    <li v-for="cardId in playerHand" :key="cardId" class="pcard">
                        <playing-card
                            :card-id="cardId"
                            :disabled="!yourTurn"
                            @clicked="handleCardClick"
                        ></playing-card>
                    </li>
                </ul>
            </div>
        </div>
    </div>
</template>

<script>
import PlayingCard from "./PlayingCard.vue";

export default {
    components: {
        PlayingCard
    },
    data() {
        return {
            card: 0,
            isClicked: false,
            showRoundWinnerNotification: false,
            lastRoundWinner: -1,
        };
    },
    computed: {
        gameid() {
            return this.$store.getters["ws/gameID"];
        },
        playerHand() {
            return this.$store.getters["ws/cards"];
        },
        yourTurn() {
            return this.$store.getters["ws/yourturn"];
        },
        yourPoints() {
            return this.$store.getters["ws/yourPoints"];
        },
        bettingPool() {
            return this.$store.getters["ws/bettingPool"];
        },
        gameColor() {
            return this.$store.getters["ws/gameColor"];
        },
        declarerIndex() {
            return this.$store.getters["ws/declarerIndex"];
        },
        playerIndex() {
            return this.$store.getters["ws/idx"];
        },
        playedCards() {
            return this.$store.getters["ws/playedCards"];
        },
        roundWinner() {
            return this.$store.getters["ws/roundWinner"];
        },
        newRound() {
            return this.$store.getters["ws/roundcounter"];
        },
        centerCard() {
            return this.isClicked ? this.card : 0;
        },
        leftCard() {
            return this.getCardAtPosition('left');
        },
        rightCard() {
            return this.getCardAtPosition('right');
        },
        playerNames() {
            return this.$store.getters["ws/playerNames"];
        },
        Name() {
            const idx = this.playerIndex;
            return this.playerNames[idx] || `Player ${idx + 1}`;
        }
    },
    watch: {
        newRound(val) {
            this.card = 0;
            this.isClicked = false;
            // Hide notification when new round starts
            this.showRoundWinnerNotification = true;
        },
        roundWinner(newWinner, oldWinner) {
            // Show notification when round winner changes (and it's valid)
            if (newWinner !== -1 && newWinner !== oldWinner) {
                this.lastRoundWinner = newWinner;
                this.showRoundWinnerNotification = false;
            }
        }
    },
    methods: {
        getPlayerName(playerIndex) {
            return this.playerNames[playerIndex] || `Player ${playerIndex + 1}`;
        },
        
        getCardAtPosition(position) {
            const roundStarter = this.roundWinner;
            if (this.playerIndex === -1 || roundStarter === -1) return null;
            
            // Calculate relative position in turn order
            const myPosition = (this.playerIndex - roundStarter + 3) % 3;
            
            let targetPlayerIndex;
            
            if (myPosition === 0) {
                // I'm the starter (first)
                if (position === 'left') targetPlayerIndex = (roundStarter + 2) % 3;
                if (position === 'right') targetPlayerIndex = (roundStarter + 1) % 3;
            } else if (myPosition === 1) {
                // I'm second
                if (position === 'left') targetPlayerIndex = (roundStarter + 2) % 3;
                if (position === 'right') targetPlayerIndex = roundStarter;
            } else if (myPosition === 2) {
                // I'm third (last)
                if (position === 'left') targetPlayerIndex = roundStarter;
                if (position === 'right') targetPlayerIndex = (roundStarter + 1) % 3;
            }
            
            return this.playedCards.find(card => card.playerIndex === targetPlayerIndex) || null;
        },
        
        handleCardClick(cardId) {
            if (!this.yourTurn) {
                console.log("it not your turn you can't put a card down");
                return;
            }
            console.log("Player choose: ", cardId);
            this.card = cardId;

            this.$store.commit("ws/deleteCard", cardId);

            // Load the card in the middle
            this.isClicked = true;

            const event = {
                type: "card_clicked",
                id: this.gameid,
                payload: { cardId },
            };

            this.$store.dispatch("ws/sendEvent", event);
            this.$store.commit("ws/setYourTurn", false);
        },
        
        getLeftPlayerIndex() {
            const roundStarter = this.roundWinner;
            if (this.playerIndex === -1 || roundStarter === -1) return -1;
            
            const myPosition = (this.playerIndex - roundStarter + 3) % 3;
            
            if (myPosition === 0) return (roundStarter + 2) % 3;
            if (myPosition === 1) return (roundStarter + 2) % 3;
            if (myPosition === 2) return (roundStarter + 1) % 3;
            return -1;
        },
        
        getRightPlayerIndex() {
            const roundStarter = this.roundWinner;
            if (this.playerIndex === -1 || roundStarter === -1) return -1;
            
            const myPosition = (this.playerIndex - roundStarter + 3) % 3;
            
            if (myPosition === 0) return (roundStarter + 1) % 3;
            if (myPosition === 1) return roundStarter;
            if (myPosition === 2) return roundStarter;
            return -1;
        }
    }
}
</script>

<style scoped>
.game-container {
    position: relative;
    width: 100%;
    height: 100%;
}

/* Side Player Names */
.side-player {
    position: fixed;
    top: 50%;
    transform: translateY(-50%);
    background-color: rgba(11, 59, 7, 0.9);
    border: 2px solid var(--primary-color);
    border-radius: 10px;
    padding: 1.5rem 1rem;
    z-index: 5;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
}

.left-player {
    left: 20px;
}

.right-player {
    right: 20px;
}

.side-player h3 {
    margin: 0;
    color: var(--text-white);
    font-size: 1.1rem;
    writing-mode: vertical-rl;
    text-orientation: mixed;
    white-space: nowrap;
}

/* Round Winner Notification */
.round-winner-notification {
    position: fixed;
    top: calc(50% - 200px);
    left: 50%;
    transform: translateX(-50%);
    color: maroon;
    padding: 1rem 2rem;
    font-size: 1rem;
    font-style: italic;
    z-index: 15;
    animation: slideInFade 0.4s ease-out;
    text-align: center;
    min-width: 250px;
}

@keyframes slideInFade {
    0% {
        opacity: 0;
        transform: translateX(-50%) translateY(-20px);
    }
    100% {
        opacity: 1;
        transform: translateX(-50%) translateY(0);
    }
}

@keyframes slideOutFade {
    0% {
        opacity: 1;
        transform: translateX(-50%) translateY(0);
    }
    100% {
        opacity: 0;
        transform: translateX(-50%) translateY(-20px);
    }
}

@media (max-width: 768px) {
    .side-player {
        padding: 1rem 0.5rem;
    }
    
    .side-player h3 {
        font-size: 0.9rem;
    }
    
    .round-winner-notification {
        font-size: 1.1rem;
        padding: 0.8rem 1.5rem;
        min-width: 200px;
        top: calc(50% - 180px);
    }
    
    .left-player {
        left: 10px;
    }
    
    .right-player {
        right: 10px;
    }
}

.played-cards-area {
    position: fixed;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    background-color: rgba(11, 59, 7, 0.95);
    border: 3px solid var(--primary-color);
    border-radius: 15px;
    padding: 2rem;
    z-index: 10;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.5);
    min-width: 600px;
}

.played-cards-area h2 {
    color: var(--text-white);
    text-align: center;
    margin: 0 0 1.5rem 0;
    font-size: 1.5rem;
    text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.5);
}

@media (max-width: 768px) {
    .played-cards-area {
        min-width: 90%;
        padding: 1.5rem 1rem;
    }
    
    .played-cards-area h2 {
        font-size: 1.2rem;
        margin: 0 0 1rem 0;
    }
}

.cards-display {
    display: flex;
    justify-content: center;
    align-items: center;
    gap: 2rem;
    margin-top: 1rem;
    min-height: 150px;
}

.card-position {
    flex: 1;
    max-width: 200px;
    min-height: 100px;
    border: 2px solid rgba(255, 255, 255, 0.5);
    border-radius: 8px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 0.5rem;
    background-color: rgba(255, 255, 255, 0.1);
}

.card-content {
    width: 100%;
    padding: 1rem;
    color: var(--text-white);
}

.card-position div {
    text-align: center;
}

.card-position p {
    margin: 0.25rem 0;
}

.betting-pool-display {
    position: fixed;
    top: calc(5rem + 30px);
    left: 20px;
    background-color: rgba(255, 255, 255, 0.95);
    border: 2px solid var(--primary-color);
    border-radius: 8px;
    padding: 1rem 1.5rem;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    z-index: 10;
}

.betting-pool-display h3 {
    margin: 0;
    color: #333;
    font-size: 1.2rem;
}

.points-display {
    position: fixed;
    top: calc(5rem + 90px);
    left: 20px;
    background-color: rgba(255, 255, 255, 0.95);
    border: 2px solid var(--primary-color);
    border-radius: 8px;
    padding: 1rem 1.5rem;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    z-index: 10;
}

.points-display h3 {
    margin: 0;
    color: #333;
    font-size: 1.2rem;
}

.game-color-display {
    position: fixed;
    top: calc(5rem + 150px);
    left: 20px;
    z-index: 10;
}

.game-color-display p {
    margin: 0;
    color: #000;
    font-size: 1rem;
}

@media (max-width: 768px) {
    .betting-pool-display {
        top: calc(5rem + 20px);
        left: 10px;
        padding: 0.75rem 1rem;
    }
    
    .betting-pool-display h3 {
        font-size: 1rem;
    }
    
    .points-display {
        top: calc(5rem + 75px);
        left: 10px;
        padding: 0.75rem 1rem;
    }
    
    .points-display h3 {
        font-size: 1rem;
    }
    
    .game-color-display {
        top: calc(5rem + 130px);
        left: 10px;
    }
    
    .game-color-display p {
        font-size: 0.9rem;
    }
}

/* Player Section with Background */
.player-section {
    position: fixed;
    bottom: 42px;
    left: 0;
    right: 0;
    display: flex;
    flex-direction: column;
    align-items: center;
    z-index: 15;
}

@media (max-width: 768px) {
    .player-section {
        bottom: 38px;
    }
}

.current-player-name {
    font-weight: bold;
    color: var(--text-white);
    font-size: 1.3rem;
    margin: 0 0 0.5rem 0;
    text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.5);
    background-color: rgba(11, 59, 7, 0.9);
    padding: 0.5rem 2rem;
    border-radius: 10px 10px 0 0;
    border: 2px solid var(--primary-color);
    border-bottom: none;
}

@media (max-width: 768px) {
    .current-player-name {
        font-size: 1.1rem;
        padding: 0.4rem 1.5rem;
    }
}

.hand-background {
    background: linear-gradient(to bottom, rgba(11, 59, 7, 0.95), rgba(11, 59, 7, 0.98));
    border-top: 3px solid var(--primary-color);
    width: 100%;
    padding: 1rem 0 1.5rem 0;
    box-shadow: 0 -4px 20px rgba(0, 0, 0, 0.4);
}

.hand {
    display: flex;
    justify-content: center;
    gap: 0.5rem;
    list-style: none;
    padding: 0;
    margin: 0;
}

@media (max-width: 768px) {
    .hand {
        gap: 0.3rem;
    }
    
    .hand-background {
        padding: 0.75rem 0 1rem 0;
    }
}

.pcard {
    display: flex;
    align-items: flex-end;
    justify-content: space-around;
}
</style>
