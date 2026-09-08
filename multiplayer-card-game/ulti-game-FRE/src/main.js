import { createApp } from 'vue'

import App from './App.vue'
import router from './route';
import store from './store/index.js';
import BaseCard from './components/layout/BaseCard.vue';
import BaseButton from './components/layout/BaseButton.vue';

const app = createApp(App);

app.component('base-card', BaseCard);
app.component('base-button', BaseButton);

app.use(router);
app.use(store);

app.mount('#app');
