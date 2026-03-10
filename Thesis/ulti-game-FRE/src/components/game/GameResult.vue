<template>
  <div class="result-container">
    <div class="result-card">
      <h2 class="result-title">Game Over!</h2>
      
      <div class="winner-section">
        <h3 v-if="gameResult.isDeclarerWin">
          🎉 {{ playerNames[declarerIndex] }} (Declarer) Wins! 🎉
        </h3>
        <h3 v-else>
          🎉 Defenders Win! 🎉
        </h3>
      </div>

      <div class="points-section">
        <h4>Your Result:</h4>
        <p class="points" :class="pointsClass">
          {{ pointsDisplay }}
        </p>
      </div>

      <div v-if="waitingForPlayers" class="waiting-section">
        <p>Waiting for other players to decide...</p>
        <p class="player-status">{{ playersReady }} / {{ totalPlayers }} ready</p>
      </div>

      <div class="action-buttons">
        <button 
          @click="playAgain" 
          class="action-button play-again"
          :disabled="hasVoted"
        >
          {{ hasVoted ? 'Waiting...' : 'Play Again' }}
        </button>
        <button 
          @click="exitToHub" 
          class="action-button exit"
        >
          Exit to Game Hub
        </button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      hasVoted: false
    };
  },
  computed: {
    gameResult() {
      return this.$store.getters['ws/gameResult'];
    },
    declarerIndex() {
      return this.$store.getters['ws/declarerIndex'];
    },
    playerNames() {
      return this.$store.getters['ws/playerNames'];
    },
    myIndex() {
      return this.$store.getters['ws/idx'];
    },
    playAgainVotes() {
      return this.$store.getters['ws/playAgainVotes'] || {};
    },
    waitingForPlayers() {
      return this.hasVoted && Object.keys(this.playAgainVotes).length < Object.keys(this.playerNames).length;
    },
    playersReady() {
      return Object.values(this.playAgainVotes).filter(v => v === true).length;
    },
    totalPlayers() {
      return Object.keys(this.playerNames).length;
    },
    pointsClass() {
      if (this.gameResult.points > 0) return 'positive';
      if (this.gameResult.points < 0) return 'negative';
      return 'neutral';
    },
    pointsDisplay() {
      const points = this.gameResult.points;
      if (points > 0) return `+${points} points`;
      if (points < 0) return `${points} points`;
      return 'No points';
    }
  },
  methods: {
    playAgain() {
      this.hasVoted = true;
      const event = {
        type: 'play_again',
        id: this.$store.getters['ws/gameID'],
        payload: { vote: true }
      };
      this.$store.dispatch('ws/sendEvent', event);
    },
    exitToHub() {
      if (this.hasVoted) {
        // Send denial vote
        const event = {
          type: 'play_again',
          id: this.$store.getters['ws/gameID'],
          payload: { vote: false }
        };
        this.$store.dispatch('ws/sendEvent', event);
      }
      
      // Exit and go to game hub
      this.$store.dispatch('ws/exitGameAndLogout');
      this.$router.push('/game');
    }
  }
};
</script>

<style scoped>
.result-container {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.8);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.result-card {
  background-color: rgba(11, 59, 7, 0.95);
  border: 3px solid var(--primary-color);
  border-radius: 15px;
  padding: 3rem;
  max-width: 500px;
  width: 90%;
  color: var(--text-white);
  text-align: center;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.5);
}

.result-title {
  font-size: 2.5rem;
  margin-bottom: 2rem;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.5);
}

.winner-section {
  background-color: rgba(255, 215, 0, 0.1);
  border: 2px solid gold;
  border-radius: 10px;
  padding: 1.5rem;
  margin-bottom: 2rem;
}

.winner-section h3 {
  font-size: 1.8rem;
  color: #ffd700;
  margin: 0;
}

.points-section {
  margin-bottom: 2rem;
}

.points-section h4 {
  font-size: 1.3rem;
  margin-bottom: 1rem;
}

.points {
  font-size: 2rem;
  font-weight: bold;
}

.points.positive {
  color: #28a745;
}

.points.negative {
  color: #dc3545;
}

.points.neutral {
  color: #ffc107;
}

.waiting-section {
  background-color: rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  padding: 1rem;
  margin-bottom: 2rem;
}

.player-status {
  font-size: 1.2rem;
  color: #ffd700;
  margin-top: 0.5rem;
}

.action-buttons {
  display: flex;
  gap: 1rem;
  justify-content: center;
}

.action-button {
  padding: 1rem 2rem;
  font-size: 1.1rem;
  font-weight: bold;
  border: 2px solid var(--text-white);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.play-again {
  background-color: #28a745;
  color: var(--text-white);
}

.play-again:hover:not(:disabled) {
  background-color: #218838;
  transform: scale(1.05);
}

.play-again:disabled {
  background-color: #6c757d;
  cursor: not-allowed;
  opacity: 0.6;
}

.exit {
  background-color: #dc3545;
  color: var(--text-white);
}

.exit:hover {
  background-color: #c82333;
  transform: scale(1.05);
}
</style>
