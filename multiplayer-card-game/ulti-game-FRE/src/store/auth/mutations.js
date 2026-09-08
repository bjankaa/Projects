export default {
    setUser(state, payload) {
        state.token = payload.token;
        state.email = payload.email || '';
        state.loggedIn = true;
    },
    setUserLogout(state) {
        state.token = null;
        state.email = '';
        state.loggedIn = false;
    },
    setEmail(state, email) {
        state.email = email;
    }
}