<template>
    <div class="profile-container">
        <base-card>
            <div class="profile-header">
                <div class="avatar-placeholder">
                    {{ userInitial }}
                </div>
                <h1 class="username">{{ username }}</h1>
                <p class="score" v-if="typeof score === 'number'">Score: {{ score }}</p>
            </div>

            <div class="profile-section">
                <h2>Game History</h2>
                <div v-if="historyLoading" class="game-history-placeholder">
                    <p class="placeholder-text">Loading game history...</p>
                </div>
                <div v-else-if="games.length === 0" class="game-history-placeholder">
                    <p class="placeholder-text">No game history yet</p>
                    <p class="placeholder-subtext">Your completed games will appear here</p>
                </div>
                <div v-else class="game-history-list">
                    <div v-for="game in games" :key="game.id" class="game-item">
                        <div class="game-row">
                            <span class="date">{{ formatDate(game.created_at) }}</span>
                            <span class="players">P1: {{ game.player1_name }} · P2: {{ game.player2_name }} · P3: {{ game.player3_name }}</span>
                        </div>
                        <div class="game-row">
                            <span class="declarer">Declarer: {{ game.declarer_name }}</span>
                            <span class="result">Declarer won? : {{ game.declarer_win ? 'Yes' : 'No' }}</span>
                        </div>
                        <div class="game-row points">
                            <span class="dec-points">Declarer Points: {{ game.declarer_points }}</span>
                            <span class="def-points">Defenders Points: {{ game.defenders_points }}</span>
                        </div>
                    </div>
                </div>
            </div>
        </base-card>

        <back-button>
            <span class="arrow">←</span> Back
        </back-button>
    </div>
</template>

<script>
import BackButton from '../components/layout/BackButton.vue';

export default {
    components: {
        BackButton
    },
    data() {
        return {
            username: 'Loading...',
            isLoading: true,
            error: null,
            hasLoadedOnce: false,
            games: [],
            historyLoading: false
        };
    },
    computed: {
        userInitial() {
            return this.username && this.username !== 'Loading...' ? this.username.charAt(0).toUpperCase() : '?';
        },
        token() {
            return this.$store.getters['auth/token'];
        },
        score() {
            return this.$store.getters['auth/score'];
        },
        userId() {
            return this.$store.getters['auth/userId'];
        }
    },
    watch: {
        userId(newId, oldId) {
            if (newId !== oldId && newId) {
                this.hasLoadedOnce = false;
                this.username = 'Loading...';
                this.loadProfile();
            }
        }
    },
    async mounted() {
        if (!this.hasLoadedOnce) {
            await this.loadProfile();
            await this.loadHistory();
        }
    },
    methods: {
        async loadProfile() {
            try {
                const data = await this.$store.dispatch('auth/fetchProfile');
                this.username = data.name;
                if (typeof data.score === 'number') {
                    this.$store.commit('auth/setScore', data.score, { root: true });
                }
                this.isLoading = false;
                this.hasLoadedOnce = true;
            } catch (error) {
                // If unauthorized, redirect to auth page
                if (error && error.message && /Failed to fetch profile|unauthorized|expired/i.test(error.message)) {
                    this.$router.push('/auth');
                    return;
                }
                this.error = 'Failed to load profile';
                this.username = 'Error loading name';
                this.isLoading = false;
            }
        },
        async loadHistory() {
            this.historyLoading = true;
            try {
                const res = await fetch('/api/games', {
                    method: 'GET',
                    headers: {
                        'Authorization': `Bearer ${this.token}`
                    }
                });
                if (!res.ok) throw new Error('Failed to fetch games');
                const data = await res.json();
                this.games = Array.isArray(data) ? data : [];
            } catch (e) {
                console.error('History load failed:', e);
                this.games = [];
            } finally {
                this.historyLoading = false;
            }
        },
        formatDate(ts) {
            try {
                const d = new Date(ts);
                return d.toLocaleString();
            } catch {
                return ts;
            }
        }
    }
};
</script>

<style scoped>
.profile-container {
    max-width: 900px;
    margin: 2rem auto;
    padding: 1rem;
}

.profile-header {
    text-align: center;
    padding: 2rem 0;
    border-bottom: 2px solid #e0e0e0;
    margin-bottom: 2rem;
}

.avatar-placeholder {
    width: 120px;
    height: 120px;
    border-radius: 50%;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    font-size: 3rem;
    font-weight: bold;
    display: flex;
    align-items: center;
    justify-content: center;
    margin: 0 auto 1.5rem;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.username {
    font-size: 2.5rem;
    color: #333;
    margin: 0;
    font-weight: 600;
}

.score {
    margin-top: 0.5rem;
    font-size: 1.2rem;
    color: #555;
}

.profile-section {
    padding: 2rem 0;
}

.profile-section h2 {
    color: #555;
    font-size: 1.8rem;
    margin-bottom: 1.5rem;
    text-align: center;
}

.game-history-placeholder {
    background-color: #f8f9fa;
    border: 2px dashed #ddd;
    border-radius: 12px;
    padding: 4rem 2rem;
    text-align: center;
    min-height: 300px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
}

.game-history-list {
    display: flex;
    flex-direction: column;
    gap: 1rem;
}

.game-item {
    border: 1px solid #e5e7eb;
    border-radius: 10px;
    padding: 1rem;
    background-color: #fff;
    box-shadow: 0 2px 8px rgba(0,0,0,0.05);
}

.game-row {
    display: flex;
    justify-content: space-between;
    font-size: 0.95rem;
    color: #444;
}

.game-row.points {
    margin-top: 0.5rem;
    font-weight: 500;
}

.placeholder-text {
    font-size: 1.3rem;
    color: #999;
    margin: 0 0 0.5rem 0;
    font-weight: 500;
}

.placeholder-subtext {
    font-size: 1rem;
    color: #aaa;
    margin: 0;
}

.back-button {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.5rem 1rem;
    background-color: transparent;
    border: 1px solid #ddd;
    border-radius: 6px;
    color: #666;
    font-size: 0.9rem;
    cursor: pointer;
    transition: all 0.3s ease;
    margin-top: 2rem;
}

.back-button:hover {
    background-color: #f5f5f5;
    border-color: #999;
    color: #333;
}

.arrow {
    font-size: 1.2rem;
    font-weight: bold;
}
</style>