export default {
    setUser(state, payload) {
        state.token = payload.token;
        state.email = payload.email || '';
        state.loggedIn = true;
        if (typeof payload.score === 'number') {
            state.score = payload.score;
        }
    },
    setUserLogout(state) {
        state.token = null;
        state.email = '';
        state.loggedIn = false;
        state.score = 0;
    },
    setEmail(state, payload) {
        state.email = payload;
    },
    setScore(state, payload) {
        state.score = payload;
    },
    addScore(state, payload) {
        state.score += payload;
    }
}