import 'babel-polyfill'
import Vue from 'vue'
import './plugins/vuetify'
import App from './App'
import router from './router'
import store from './store'

Vue.config.productionTip = false

Vue.prototype.$axios = store.axiosClient

Vue.prototype.$store = store.sessionStore
Vue.prototype.$localStore = store.localStore

new Vue({
  router,
  render: h => h(App)
}).$mount('#app')
