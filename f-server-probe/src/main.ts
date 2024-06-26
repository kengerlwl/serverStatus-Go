import './assets/main.css'

import { createApp } from 'vue'
import App from './App.vue'
import store from './store'

import Antd from 'ant-design-vue';
import 'ant-design-vue/dist/reset.css';
// import 'ant-design-vue/dist/antd.css';

const app = createApp(App)

app.use(store)
app.use(Antd)
app.mount('#app')
