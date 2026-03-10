import { createRouter, createWebHistory } from "vue-router";

import FrontPage from "./pages/FrontPage.vue";
import AuthPage from "./pages/AuthPage.vue";
import GameHub from "./pages/GameHub.vue";
import Profile from "./pages/Profile.vue";
import Settings from "./pages/Settings.vue";
import NotFound from "./pages/NotFound.vue";
import GamePage from "./pages/GamePage.vue";


const router = createRouter({
    history: createWebHistory(),
    routes: [
        {path: '/', redirect: '/frontpage'},
        {path: '/frontpage', component: FrontPage},
        {path: '/auth', component: AuthPage},
        {path: '/game', component: GameHub,},
        {path: '/game/:id', component: GamePage},
        {path: '/profile', component: Profile},
        {path: '/settings', component: Settings},
        {path: '/:notFound(.*)', component: NotFound}
    ]
});

export default router;