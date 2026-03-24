import { createApp } from 'vue'
import { createPinia } from 'pinia'
import vuetify from './plugins/vuetify'
import router from './router'
import i18n from './i18n'
import App from './App.vue'
import './styles/main.css'

const app = createApp(App)

app.use(createPinia())
app.use(vuetify)
app.use(router)
app.use(i18n)

app.mount('#app')
