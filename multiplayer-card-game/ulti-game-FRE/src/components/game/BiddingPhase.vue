<template>
  <div class="bidding-container">
    <div class="betting-pool-info">
      <h3>Betting Pool: {{ bettingPool }}</h3>
    </div>
    
    <div class="bidding-info">
      <h2>Bidding Phase</h2>
      <p v-if="gameColor" class="current-bid">
        Current Game Color: <strong>{{ gameColor.toUpperCase() }}</strong>
      </p>
      <p v-if="declarerIndex >= 0" class="declarer-info">
        Declarer: <strong>{{ playerNames[declarerIndex] }}</strong>
      </p>
    </div>

    <div v-if="lastBid" class="last-bid">
      <p>{{ playerNames[lastBid.bidderIndex] }} {{ lastBid.action === 'pass' ? 'passed' : `declared ${lastBid.color}` }}</p>
    </div>

    <div v-if="isYourBid" class="bidding-actions">
      <h3>Your turn to bid</h3>
      
      <button @click="passBid" class="bid-button pass-button">
        Pass
      </button>

      <div class="declare-section">
        <h4>Declare with color:</h4>
        <div class="color-buttons">
          <button 
            v-for="color in availableColors" 
            :key="color"
            @click="declareBid(color)" 
            :class="['bid-button', 'color-button', color]"
            :disabled="!canDeclareColor(color)"
          >
            {{ color.toUpperCase() }}
          </button>
        </div>
      </div>
    </div>

    <div v-else class="waiting-message">
      <p>Waiting for {{ playerNames[currentBidder] }} to bid...</p>
    </div>

    <div class="hand-display">
      <h3>Your Hand ({{ playerHand.length }} cards):</h3>
      <div class="hand-cards">
        <PlayingCard 
          v-for="cardId in playerHand" 
          :key="cardId" 
          :card-id="cardId"
          :disabled="true"
        />
      </div>
    </div>
  </div>
</template>

<script>
import PlayingCard from './PlayingCard.vue';

export default {
  components: {
    PlayingCard
  },
  computed: {
    isYourBid() {
      const value = this.$store.getters['ws/isYourBid'];
      console.log('BiddingPhase - isYourBid:', value);
      return value;
    },
    currentBidder() {
      const value = this.$store.getters['ws/currentBidder'];
      console.log('BiddingPhase - currentBidder:', value);
      return value;
    },
    lastBid() {
      return this.$store.getters['ws/lastBid'];
    },
    gameColor() {
      return this.$store.getters['ws/gameColor'];
    },
    declarerIndex() {
      return this.$store.getters['ws/declarerIndex'];
    },
    playerNames() {
      return this.$store.getters['ws/playerNames'];
    },
    bettingPool() {
      return this.$store.getters['ws/bettingPool'];
    },
    playerHand() {
      return this.$store.getters['ws/cards'];
    },
    availableColors() {
      return ['tok', 'makk', 'zold', 'piros'];
    },
    colorStrength() {
      return {
        'tok': 1,
        'makk': 2,
        'zold': 3,
        'piros': 4
      };
    }
  },
  mounted() {
    console.log('BiddingPhase mounted - isYourBid:', this.isYourBid, 'currentBidder:', this.currentBidder);
  },
  methods: {
    passBid() {
      this.$store.dispatch('ws/sendBidAction', {
        action: 'pass',
        color: ''
      });
    },
    declareBid(color) {
      this.$store.dispatch('ws/sendBidAction', {
        action: 'declare',
        color: color
      });
    },
    canDeclareColor(color) {
      // If no game color yet, all colors are available
      if (!this.gameColor) {
        return true;
      }
      
      // Must choose stronger color to one-up, or piros can be chosen again
      const currentStrength = this.colorStrength[this.gameColor];
      const newStrength = this.colorStrength[color];
      
      return newStrength > currentStrength || (newStrength === currentStrength && color === 'piros');
    }
  }
};
</script>

<style scoped>
.bidding-container {
  background-color: rgba(11, 59, 7, 0.9);
  border: 2px solid var(--primary-color);
  border-radius: 10px;
  padding: 2rem;
  max-width: 600px;
  margin: 2rem auto;
  color: var(--text-white);
}

.betting-pool-info {
  text-align: center;
  background-color: rgba(255, 255, 255, 0.15);
  border-radius: 8px;
  padding: 1rem;
  margin-bottom: 1.5rem;
  border: 2px solid rgba(255, 255, 255, 0.3);
}

.betting-pool-info h3 {
  margin: 0;
  font-size: 1.3rem;
  color: #ffd700;
}

.bidding-info {
  text-align: center;
  margin-bottom: 1.5rem;
}

.bidding-info h2 {
  margin-bottom: 1rem;
}

.current-bid {
  font-size: 1.2rem;
  margin: 0.5rem 0;
}

.declarer-info {
  font-size: 1.1rem;
  color: #ffd700;
}

.last-bid {
  background-color: rgba(255, 255, 255, 0.1);
  padding: 1rem;
  border-radius: 5px;
  text-align: center;
  margin-bottom: 1.5rem;
}

.bidding-actions {
  text-align: center;
}

.bidding-actions h3 {
  margin-bottom: 1rem;
  color: #ffd700;
}

.bid-button {
  padding: 0.75rem 1.5rem;
  font-size: 1rem;
  font-weight: bold;
  border: 2px solid var(--text-white);
  border-radius: 5px;
  cursor: pointer;
  transition: all 0.3s ease;
  margin: 0.5rem;
}

.pass-button {
  background-color: #dc3545;
  color: var(--text-white);
}

.pass-button:hover {
  background-color: #c82333;
  transform: scale(1.05);
}

.declare-section {
  margin-top: 1.5rem;
}

.declare-section h4 {
  margin-bottom: 1rem;
}

.color-buttons {
  display: flex;
  justify-content: center;
  gap: 1rem;
  flex-wrap: wrap;
}

.color-button {
  color: var(--text-white);
  text-shadow: 1px 1px 2px rgba(0, 0, 0, 0.8);
}

.color-button.tok {
  background-color: #8B4513;
}

.color-button.makk {
  background-color: #2E8B57;
}

.color-button.zold {
  background-color: #228B22;
}

.color-button.piros {
  background-color: #DC143C;
}

.color-button:hover:not(:disabled) {
  transform: scale(1.1);
  box-shadow: 0 0 10px rgba(255, 255, 255, 0.5);
}

.color-button:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.waiting-message {
  text-align: center;
  font-size: 1.2rem;
  padding: 2rem;
}

.waiting-message p {
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.hand-display {
  margin-top: 2rem;
  padding-top: 1.5rem;
  border-top: 2px solid rgba(255, 255, 255, 0.2);
}

.hand-display h3 {
  text-align: center;
  margin-bottom: 1rem;
  color: #ffd700;
}

.hand-cards {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 0.5rem;
  padding: 0 1rem;
}
</style>
