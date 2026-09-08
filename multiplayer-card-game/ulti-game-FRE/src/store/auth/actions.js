export default {
    async login(context, payload) {

        const res = await fetch('/api/auth', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                email: payload.email,
                password: payload.password,
                state: "login"
            })
        })

        const resData = await res.json();

        if (!res.ok) {
            const error = new Error(resData.message || 'Failed to Login in!');
            throw error;
        }

        localStorage.setItem('token', resData.token);


        context.commit('setUser', {
            token: resData.token
        });

        try {
            const profile = await context.dispatch('fetchProfile');
            if (profile && typeof profile.score === 'number') {
                context.commit('setScore', profile.score);
            }
        } catch (e) {
            console.error('Failed to fetch profile after login:', e);
        }

    },
    async signup(context, payload) {

        const res = await fetch('/api/auth', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                name: payload.name,
                email: payload.email,
                password: payload.password,
                state: "signup"
            })
        })
        const resData = await res.json();

        if (!res.ok) {
            const error = new Error(resData.message || 'Failed to sign up!');
            throw error;
        }
    },

    async logout(context) {
        const token = context.getters.token;

        if (token) {
            await fetch('/api/logout', {
                method: 'POST',
                headers: {
                    'Authorization': 'Bearer ' + token
                }
            }).catch((err) => {
                console.error('logout error:', err);
            });
        }

        localStorage.removeItem('token');
        context.commit('setUserLogout');
    },

    async fetchProfile(context) {
        const token = context.getters.token;

        const res = await fetch('/api/profile', {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });

        if (!res.ok) {
            const error = new Error('Failed to fetch profile');
            throw error;
        }

        const data = await res.json();

        if (data && typeof data.email === 'string') {
            context.commit('setEmail', data.email);
        }
        if (data && typeof data.score === 'number') {
            context.commit('setScore', data.score);
        }
        return data;
    },

    async changeEmail(context, payload) {
        const token = context.getters.token;

        const res = await fetch('/api/change-email', {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({
                newEmail: payload.newEmail,
                password: payload.password
            })
        });

        const data = await res.json();

        if (!res.ok) {
            const error = new Error(data.message || 'Failed to update email');
            throw error;
        }

        // Update store with new email
        context.commit('setEmail', payload.newEmail);
        return data;
    },

    async changePassword(context, payload) {
        const token = context.getters.token;

        const res = await fetch('/api/change-password', {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({
                currentPassword: payload.currentPassword,
                newPassword: payload.newPassword
            })
        });

        const data = await res.json();

        if (!res.ok) {
            const error = new Error(data.message || 'Failed to update password');
            throw error;
        }

        return data;
    }
}
