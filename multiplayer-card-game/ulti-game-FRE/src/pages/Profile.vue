<template>
    <div class="profile-container">
        <base-card>
            <div class="profile-header">
                <div class="avatar-placeholder">
                    {{ userInitial }}
                </div>
                <h1 class="username">{{ username }}</h1>
            </div>

            <div class="profile-section">
                <h2>Game History</h2>
                <div class="game-history-placeholder">
                    <p class="placeholder-text">No game history yet</p>
                    <p class="placeholder-subtext">Your completed games will appear here</p>
                </div>
            </div>
        </base-card>

        <back-button @click="$router.go(-1)">
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
            error: null
        };
    },
    computed: {
        userInitial() {
            return this.username && this.username !== 'Loading...' ? this.username.charAt(0).toUpperCase() : '?';
        },
        token() {
            return this.$store.getters['auth/token'];
        }
    },
    async mounted() {
        await this.loadProfile();
    },
    methods: {
        async loadProfile() {
            try {
                const data = await this.$store.dispatch('auth/fetchProfile');
                this.username = data.name;
                this.isLoading = false;
            } catch (error) {
                this.error = 'Failed to load profile';
                this.username = 'Error loading name';
                this.isLoading = false;
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