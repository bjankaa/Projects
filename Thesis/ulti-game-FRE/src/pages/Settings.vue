<template>
    <div class="settings-container">
        <base-card>
            <h1>Settings</h1>
            
            <!-- Change Email Section -->
            <section class="settings-section">
                <h2>Change Email</h2>
                <form @submit.prevent="changeEmail">

                    <div class="form-group">
                        <label for="new-email">New Email</label>
                        <input 
                            type="email" 
                            id="new-email" 
                            v-model="emailForm.newEmail" 
                            required
                            class="input-field"
                            placeholder="Enter new email"
                        />
                    </div>
                    <div class="form-group">
                        <label for="confirm-email">Confirm New Email</label>
                        <input 
                            type="email" 
                            id="confirm-email" 
                            v-model="emailForm.confirmEmail" 
                            required
                            class="input-field"
                            placeholder="Confirm new email"
                        />
                    </div>
                    <div class="form-group">
                        <label for="email-password">Current Password</label>
                        <input 
                            type="password" 
                            id="email-password" 
                            v-model="emailForm.password" 
                            required
                            class="input-field"
                            placeholder="Enter password to confirm"
                        />
                    </div>
                    <p v-if="emailMessage" :class="emailError ? 'error-message' : 'success-message'">
                        {{ emailMessage }}
                    </p>
                    <submit-button :disabled="emailLoading">{{ emailLoading ? 'Updating...' : 'Update Email' }}</submit-button>
                </form>
            </section>

            <!-- Change Password Section -->
            <section class="settings-section">
                <h2>Change Password</h2>
                <form @submit.prevent="changePassword">
                    <div class="form-group">
                        <label for="current-password">Current Password</label>
                        <input 
                            type="password" 
                            id="current-password" 
                            v-model="passwordForm.currentPassword" 
                            required
                            class="input-field"
                            placeholder="Enter current password"
                        />
                    </div>
                    <div class="form-group">
                        <label for="new-password">New Password</label>
                        <input 
                            type="password" 
                            id="new-password" 
                            v-model="passwordForm.newPassword" 
                            required
                            class="input-field"
                            placeholder="Enter new password"
                        />
                    </div>
                    <div class="form-group">
                        <label for="confirm-password">Confirm New Password</label>
                        <input 
                            type="password" 
                            id="confirm-password" 
                            v-model="passwordForm.confirmPassword" 
                            required
                            class="input-field"
                            placeholder="Confirm new password"
                        />
                    </div>
                    <p v-if="passwordMessage" :class="passwordError ? 'error-message' : 'success-message'">
                        {{ passwordMessage }}
                    </p>
                    <submit-button :disabled="passwordLoading">{{ passwordLoading ? 'Updating...' : 'Update Password' }}</submit-button>
                </form>
            </section>
        </base-card>
        <back-button></back-button>
    </div>
</template>

<script>
import BackButton from '../components/layout/BackButton.vue';
import SubmitButton from '../components/layout/SubmitButton.vue';

export default {
    components: {
        BackButton,
        SubmitButton
    },
    data() {
        return {
            emailForm: {
                currentEmail: '',
                newEmail: '',
                confirmEmail: '',
                password: ''
            },
            passwordForm: {
                currentPassword: '',
                newPassword: '',
                confirmPassword: ''
            },
            emailMessage: '',
            emailError: false,
            emailLoading: false,
            passwordMessage: '',
            passwordError: false,
            passwordLoading: false
        };
    },
    computed: {
        token() {
            return this.$store.getters['auth/token'];
        },
        userEmail() {
            return this.$store.getters['auth/email'];
        }
    },
    mounted() {
        this.emailForm.currentEmail = this.userEmail;
    },
    methods: {
        async changeEmail() {
            this.emailMessage = '';
            this.emailError = false;

            // Validation
            if (this.emailForm.newEmail !== this.emailForm.confirmEmail) {
                this.emailMessage = 'New emails do not match!';
                this.emailError = true;
                return;
            }

            if (this.emailForm.newEmail === this.emailForm.currentEmail) {
                this.emailMessage = 'New email must be different from current email!';
                this.emailError = true;
                return;
            }

            this.emailLoading = true;

            try {
                await this.$store.dispatch('auth/changeEmail', {
                    newEmail: this.emailForm.newEmail,
                    password: this.emailForm.password
                });

                this.emailMessage = 'Email updated successfully!';
                this.emailError = false;
                this.emailForm.currentEmail = this.emailForm.newEmail;
                this.emailForm.newEmail = '';
                this.emailForm.confirmEmail = '';
                this.emailForm.password = '';
            } catch (error) {
                this.emailMessage = error.message || 'Network error. Please try again.';
                this.emailError = true;
            } finally {
                this.emailLoading = false;
            }
        },

        async changePassword() {
            this.passwordMessage = '';
            this.passwordError = false;

            // Validation
            if (this.passwordForm.newPassword !== this.passwordForm.confirmPassword) {
                this.passwordMessage = 'New passwords do not match!';
                this.passwordError = true;
                return;
            }

            if (this.passwordForm.newPassword.length < 5) {
                this.passwordMessage = 'New password must be at least 5 characters long!';
                this.passwordError = true;
                return;
            }

            this.passwordLoading = true;

            try {
                await this.$store.dispatch('auth/changePassword', {
                    currentPassword: this.passwordForm.currentPassword,
                    newPassword: this.passwordForm.newPassword
                });

                this.passwordMessage = 'Password updated successfully!';
                this.passwordError = false;
                this.passwordForm.currentPassword = '';
                this.passwordForm.newPassword = '';
                this.passwordForm.confirmPassword = '';
            } catch (error) {
                this.passwordMessage = error.message || 'Network error. Please try again.';
                this.passwordError = true;
            } finally {
                this.passwordLoading = false;
            }
        }
    }
};
</script>

<style scoped>
.settings-container {
    max-width: 800px;
    margin: 2rem auto;
    padding: 1rem;
}

h1 {
    color: #333;
    margin-bottom: 2rem;
    text-align: center;
}

.settings-section {
    margin-bottom: 3rem;
    padding-bottom: 2rem;
    border-bottom: 1px solid #e0e0e0;
}

.settings-section:last-child {
    border-bottom: none;
}

h2 {
    color: #555;
    margin-bottom: 1.5rem;
    font-size: 1.5rem;
}

.form-group {
    margin-bottom: 1.5rem;
}

label {
    display: block;
    margin-bottom: 0.5rem;
    color: #666;
    font-weight: 500;
}

.input-field {
    width: 100%;
    padding: 0.75rem;
    border: 1px solid #ddd;
    border-radius: 6px;
    font-size: 1rem;
    transition: border-color 0.3s ease;
}

.input-field:focus {
    outline: none;
    border-color: var(--primary-color);
}

.input-field:disabled {
    background-color: #f5f5f5;
    cursor: not-allowed;
}

.error-message {
    color: #dc3545;
    margin: 1rem 0;
    padding: 0.75rem;
    background-color: #f8d7da;
    border: 1px solid #f5c6cb;
    border-radius: 6px;
}

.success-message {
    color: #28a745;
    margin: 1rem 0;
    padding: 0.75rem;
    background-color: #d4edda;
    border: 1px solid #c3e6cb;
    border-radius: 6px;
}
</style>