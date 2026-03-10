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
        };
    },
    getters,
    actions,
    mutations
}