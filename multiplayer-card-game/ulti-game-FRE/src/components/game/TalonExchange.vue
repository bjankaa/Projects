<template>
  <div class="talon-container">
    <div v-if="isDeclarer" class="talon-exchange">
      <h2>Talon Exchange</h2>
      <p class="instructions">You received the talon! Reveal it to see what cards you got, then select exactly 2 cards to discard.</p>
      
      <div class="talon-cards-section">
        <h3>Talon Cards (received):</h3>
        <div v-if="!talonRevealed" class="hidden-talon">
          <button @click="revealTalon" class="reveal-button">
            🎴 Reveal Talon Cards
          </button>
          <p class="skip-hint">Or keep your original hand:</p>
          <button @click="skipExchange" class="skip-button">
            Skip Exchange (Discard Talon)
          </button>
        </div>
        <div v-else class="revealed-talon">
          <div class="card-list">
            <PlayingCard 
              v-for="cardId in talonCards" 
              :key="'talon-' + cardId"
              :card-id="cardId"
              :disabled="true"
            />
          </div>
        </div>
      </div>

      <div class="original-hand" v-if="!talonRevealed">
        <h3>Your Original Hand ({{ originalHand.length }} cards):</h3>
        <div class="card-list">
          <PlayingCard 
            v-for="cardId in originalHand" 
            :key="'original-' + cardId"
            :card-id="cardId"
            :disabled="true"
          />
        </div>
      </div>

      <div class="your-hand" v-if="talonRevealed">
        <h3>Your Hand ({{ allCards.length }} cards - select 2 to discard):</h3>
        <div class="card-list">
          <div 
            v-for="cardId in allCards" 
            :key="'hand-' + cardId"
            @click="toggleCardSelection(cardId)"
            :class="['card-wrapper', { 'selected': isSelected(cardId), 'from-talon': isTalonCard(cardId) }]">
            <PlayingCard 
              :card-id="cardId"
              :disabled="false"
            />
            <span v-if="isTalonCard(cardId)" class="talon-badge">New</span>
          </div>
        </div>
      </div>

      <div class="selection-info" v-if="talonRevealed">
        <p>Selected for discard: {{ selectedCards.length }} / 2 cards</p>
        <p class="final-count" :class="{ 'valid': finalHandSize === 10, 'invalid': finalHandSize !== 10 }">
          Final hand size: {{ finalHandSize }} cards (must be 10)
        </p>
      </div>

      <button 
        @click="confirmExchange" 
        :disabled="!canConfirm"
        class="confirm-button"
        v-if="talonRevealed"
      >
        Confirm Exchange
      </button>
    </div>

    <div v-else class="waiting-declarer">
      <h2>Talon Exchange Phase</h2>
      <p>Waiting for {{ playerNames[declarerIndex] }} to exchange cards with the talon...</p>
    </div>
  </div>
</template>

<script>
import PlayingCard from './PlayingCard.vue';

export default {
  components: {
    PlayingCard
  },
  data() {
    return {
      talonRevealed: false,
      selectedCards: []
    };
  },
  computed: {
    isDeclarer() {
      const myIndex = this.$store.getters['ws/idx'];
      const declarerIndex = this.$store.getters['ws/declarerIndex'];
      return myIndex === declarerIndex;
    },
    declarerIndex() {
      return this.$store.getters['ws/declarerIndex'];
    },
    playerNames() {
      return this.$store.getters['ws/playerNames'];
    },
    talonCards() {
      return this.$store.getters['ws/talonCards'] || [];
    },
    allCards() {
      return this.$store.getters['ws/cards'];
    },
    originalHand() {
      // Filter out talon cards to show only original 10 cards
      return this.allCards.filter(cardId => !this.talonCards.includes(cardId));
    },
    finalHandSize() {
      return this.allCards.length - this.selectedCards.length;
    },
    canConfirm() {
      return this.talonRevealed && this.selectedCards.length === 2;
    }
  },
  methods: {
    revealTalon() {
      this.talonRevealed = true;
    },
    skipExchange() {
      console.log('TalonExchange: Skipping exchange, discarding talon cards:', this.talonCards);
      
      // Discard both talon cards to keep original hand
      this.$store.dispatch('ws/sendTalonExchange', this.talonCards);
      
      // Reset state
      this.talonRevealed = false;
      this.selectedCards = [];
    },
    toggleCardSelection(cardId) {
      const index = this.selectedCards.indexOf(cardId);
      if (index > -1) {
        this.selectedCards.splice(index, 1);
      } else {
        if (this.selectedCards.length < 2) {
          this.selectedCards.push(cardId);
        }
      }
    },
    isSelected(cardId) {
      return this.selectedCards.includes(cardId);
    },
    isTalonCard(cardId) {
      return this.talonCards.includes(cardId);
    },
    confirmExchange() {
      if (!this.canConfirm) {
        console.log('Cannot confirm - must select exactly 2 cards');
        return;
      }

      console.log('TalonExchange: Discarding cards:', this.selectedCards);
      
      this.$store.dispatch('ws/sendTalonExchange', this.selectedCards);
      
      // Reset state
      this.talonRevealed = false;
      this.selectedCards = [];
    }
  },
  mounted() {
    console.log('TalonExchange mounted - isDeclarer:', this.isDeclarer, 'talonCards:', this.talonCards);
  }
};
</script>

<style scoped>
.talon-container {
  background-color: rgba(11, 59, 7, 0.9);
  border: 2px solid var(--primary-color);
  border-radius: 10px;
  padding: 2rem;
  max-width: 900px;
  margin: 2rem auto;
  color: var(--text-white);
}

.talon-exchange h2,
.waiting-declarer h2 {
  text-align: center;
  margin-bottom: 1rem;
}

.instructions {
  text-align: center;
  font-size: 1.1rem;
  margin-bottom: 1.5rem;
  color: #ffd700;
}

.talon-cards-section {
  background-color: rgba(0, 0, 0, 0.3);
  border: 2px solid #DC143C;
  border-radius: 10px;
  padding: 1.5rem;
  margin-bottom: 2rem;
}

.talon-cards-section h3 {
  margin-bottom: 1rem;
  color: #ffd700;
}

.hidden-talon {
  text-align: center;
  padding: 2rem;
}

.reveal-button {
  padding: 1.5rem 3rem;
  font-size: 1.3rem;
  font-weight: bold;
  background: linear-gradient(135deg, #DC143C, #8B0000);
  color: white;
  border: 3px solid #FFD700;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 5px 15px rgba(220, 20, 60, 0.4);
}

.reveal-button:hover {
  transform: scale(1.05);
  box-shadow: 0 8px 25px rgba(220, 20, 60, 0.6);
}

.skip-hint {
  margin: 1.5rem 0 0.5rem 0;
  font-size: 1rem;
  color: #ffd700;
}

.skip-button {
  padding: 1rem 2rem;
  font-size: 1.1rem;
  font-weight: bold;
  background: linear-gradient(135deg, #666, #444);
  color: white;
  border: 2px solid #999;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 3px 10px rgba(0, 0, 0, 0.3);
}

.skip-button:hover {
  background: linear-gradient(135deg, #777, #555);
  transform: scale(1.05);
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.5);
}

.revealed-talon {
  animation: fadeIn 0.5s ease;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(-10px); }
  to { opacity: 1; transform: translateY(0); }
}

.exchange-hint {
  text-align: center;
  font-size: 0.9rem;
  color: #ffd700;
  margin-top: 1rem;
  font-style: italic;
}

.your-hand {
  margin: 1.5rem 0;
}

.original-hand {
  margin-top: 1.5rem;
  padding: 1rem;
  background-color: rgba(255, 255, 255, 0.05);
  border-radius: 8px;
  border: 2px solid rgba(255, 255, 255, 0.2);
}

.original-hand h3,
.your-hand h3 {
  margin-bottom: 1rem;
  color: #ffd700;
}

.card-list {
  display: flex;
  gap: 1rem;
  flex-wrap: wrap;
  justify-content: center;
}

.card-wrapper {
  position: relative;
  cursor: pointer;
  transition: all 0.3s ease;
}

.card-wrapper:hover {
  transform: translateY(-5px);
}

.card-wrapper.selected {
  transform: translateY(-10px);
}

.card-wrapper.selected .card-btn {
  background-color: #ffd700;
  border-color: #ffd700;
  color: #000;
  box-shadow: 0 8px 20px rgba(255, 215, 0, 0.5);
}

.card-wrapper.from-talon .card-btn {
  background-color: rgba(220, 20, 60, 0.3);
  border-color: #DC143C;
}

.card-wrapper.from-talon:hover .card-btn {
  background-color: rgba(220, 20, 60, 0.5);
}

.talon-badge {
  position: absolute;
  top: -8px;
  right: -8px;
  background-color: #DC143C;
  color: white;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-size: 0.7rem;
  font-weight: bold;
  z-index: 10;
}

.card-item {
  position: relative;
  padding: 1rem 1.5rem;
  background-color: rgba(255, 255, 255, 0.1);
  border: 2px solid var(--text-white);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
  min-width: 80px;
  text-align: center;
}

.card-item:hover {
  transform: translateY(-5px);
  box-shadow: 0 5px 15px rgba(255, 255, 255, 0.3);
}

.card-item.selected {
  background-color: #ffd700;
  color: #000;
  border-color: #ffd700;
  transform: translateY(-10px);
  box-shadow: 0 8px 20px rgba(255, 215, 0, 0.5);
}

.card-item.talon-card {
  background-color: rgba(220, 20, 60, 0.4);
  border-color: #DC143C;
}

.card-item.talon-card:hover {
  background-color: rgba(220, 20, 60, 0.6);
}

.card-item.talon-card.selected {
  background-color: #00ff00;
  color: #000;
  border-color: #00ff00;
}

.selection-info {
  text-align: center;
  margin: 1.5rem 0;
  padding: 1rem;
  background-color: rgba(0, 0, 0, 0.3);
  border-radius: 8px;
}

.selection-info p {
  margin: 0.5rem 0;
  font-size: 1rem;
}

.final-count {
  font-size: 1.2rem;
  font-weight: bold;
  margin-top: 1rem;
}

.final-count.valid {
  color: #00ff00;
}

.final-count.invalid {
  color: #ff4444;
}

.confirm-button {
  display: block;
  margin: 2rem auto 0;
  padding: 1rem 2rem;
  font-size: 1.1rem;
  font-weight: bold;
  background-color: #28a745;
  color: var(--text-white);
  border: 2px solid var(--text-white);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.confirm-button:hover:not(:disabled) {
  background-color: #218838;
  transform: scale(1.05);
}

.confirm-button:disabled {
  background-color: #6c757d;
  cursor: not-allowed;
  opacity: 0.5;
}

.waiting-declarer {
  text-align: center;
  padding: 3rem;
}

.waiting-declarer p {
  font-size: 1.2rem;
  color: #ffd700;
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}
</style>
