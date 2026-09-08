import mutations from "./mutations";
import actions from "./actions";
import getters from "./getters";



export default {
    namespaced: true,
    state() {
        return {
            loggedIn: false,
            token: null,
            email: '',
            score: 0,
        };
    },
    getters,
    actions,
    mutations
}