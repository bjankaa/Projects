<template>
    <div>
        <base-card>
            <form @submit.prevent="submitForm">
                <div v-if=singupInfo class="form.control" :class="{ invalid: !isFormValid }">
                    <label for="name">Name</label>
                    <input type="text" id="name" v-model.trim="name">
                </div>
                <div class="form.control" :class="{ invalid: !isFormValid }">
                    <label for="email">E-mail</label>
                    <input type="email" id="email" v-model.trim="email">
                </div>
                <div class="form.control" :class="{ invalid: !isFormValid }">
                    <label for="password">Password</label>
                    <input type="password" id="password" v-model.trim="password">
                    <small v-if="singupInfo" class="password-hint">Password must be at least 5 characters long</small>
                </div>
                <p v-if="!isFormValid" class="error-message">Please fill all the required areas in the form correctly.</p>
                <p v-if="error" class="error-message">{{ error }}</p>
                <submit-button>{{ submitButtonCaption }}</submit-button>
                <button @click.prevent="switchButtonFunc" class="link">{{ otherButtonCaption }}</button>
            </form>
        </base-card>
    </div>
</template>

<script>
import SubmitButton from '../components/layout/SubmitButton.vue';

export default {
    components: {
        SubmitButton
    },
    data() {
        return {
            name: "",
            email: "",
            password: "",
            mode: 'login',
            isFormValid: true,
            error: null
        } 
    },
    computed: {
        singupInfo() {
            return this.mode === 'signup';
        },
        submitButtonCaption() {
            return this.mode === 'login' ? 'Log in' : 'Sign up';
        },
        otherButtonCaption() {
            if (this.mode === 'login') {
                return 'Sign up here';
            } else {
                return 'Log in here';
            }
        }
    },
    methods: {
        switchButtonFunc() {
            if (this.mode === 'login') {
                console.log(this.mode);
                 return this.mode = 'signup';
             } else {
                console.log(this.mode);
                return this.mode = 'login';
             }

        },
        validation() {
            this.isFormValid = true;
            console.log(this.isFormValid);
            if (this.mode === 'signup') {
                if (this.name === '' || this.email === '' || !this.email.includes('@') || this.password.length < 5) {
                    this.isFormValid = false;
                    console.log(this.isFormValid);
                    return false;
                } else {
                    return true;
                }
            } else {
                if (this.email === '' || !this.email.includes('@') || this.password === '') {
                    this.isFormValid = false;
                    console.log(this.isFormValid);
                    return false;
                } else {
                    return true;
              
                }
            }
        },
        async submitForm() {
            if (!this.validation()) {
                return;
            }

            this.error = null; // Clear previous errors

            const loginPayload = {
                email: this.email,
                password: this.password
            };
            const signupPayload = {
                name: this.name,
                email: this.email,
                password: this.password
            };

            try {
                if (this.mode === 'login') {
                    await this.$store.dispatch("auth/login", loginPayload);
                } else {
                    await this.$store.dispatch("auth/signup", signupPayload);
                }
                this.$router.push("/frontpage");
            } catch (err) {
                this.error = err.message || 'Failed to authenticate.';
            }

        }
        }
    }


</script>

<style scoped>
form {
    margin: 1rem;
    padding: 1rem;
}

.form-control {
    margin: 0.5rem 0;
}

label {
    font-weight: bold;
    margin-bottom: 0.5rem;
    display: block;
}

input,
textarea {
    display: block;
    width: 100%;
    font: inherit;
    border: 1px solid #ccc;
    padding: 0.15rem;
    border-radius: 2rem;
}

input:focus,
textarea:focus {
    border-color: var(--primary-color);
    background-color: transparent;
    outline: none;
    border-radius: 2rem;
}

.link {
    text-decoration: none;
    font-size: large;
    margin-left: 1rem;
    background-color: transparent;
    border: 1px solid transparent;
    color: var(--primary-color);
    border-radius: 2rem;
    max-width: 15rem;
    display: inline-block;
    padding: 0.2rem 1rem;
    cursor: pointer;
}

.link:hover {
    color: var(--primary-color);
    border-color: var(--primary-color);
}

.invalid label {
    color: red;
}

.invalid input,
.invalid textarea {
    border: 1px solid red;
}

.error-message {
    color: #dc3545;
    background-color: #f8d7da;
    border: 1px solid #f5c6cb;
    border-radius: 8px;
    padding: 0.75rem;
    margin: 0.5rem 0;
    font-size: 0.95rem;
}

.password-hint {
    display: block;
    margin-top: 0.25rem;
    color: #666;
    font-size: 0.85rem;
    font-style: italic;
}
</style>