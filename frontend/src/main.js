import '@babel/polyfill'

import Vue from 'vue'

import axios from 'axios'
import VueAxios from 'vue-axios'

import store from './store'
import router from './router'
import vuetify from './plugins/vuetify';

import App from './App.vue'


Vue.use(VueAxios, axios)
Vue.config.productionTip = false

new Vue({
  store,
  router,
  vuetify,
  render: h => h(App)
}).$mount('#app')
