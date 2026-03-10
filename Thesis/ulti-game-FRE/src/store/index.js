import { createStore } from "vuex";

import authModule from './auth/auth_idx.js';
import wsModule from './ws/ws_idx.js';

const store = createStore({
    modules: {
        auth: authModule,
        ws: wsModule,
    }
});

export default store;