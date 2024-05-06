import Vue from 'vue'

// initial state
const state = {
  status: '',
  token: localStorage.getItem('survey-token') || '',
  username: '',
  isAdmin: false,
}

// getters
const getters = {
  isLoggedIn: state => !!state.token,
  isAdmin: state => state.isAdmin,
  authStatus: state => state.status,
}

// actions
const actions = {
  register({commit}, user) {
    return new Promise((resolve, reject) => {
      commit('auth_request')
            
      Vue.axios.post('/register', user)
      .then(response => {
        const token = response.data.token
                
        localStorage.setItem('survey-token', token)
        
        // set auth token
        Vue.axios.defaults.headers.common['Authorization'] = 'Bearer ' + token
        
        commit('auth_success', response.data)
        
        resolve(response)
      }).catch(function (error) {
        commit('auth_error')
        
        localStorage.removeItem('survey-token')
        
        reject(error)
      })
    })
  },
  login({commit}, user){
    return new Promise((resolve, reject) => {
      commit('auth_request')
      
      Vue.axios.post('/login', user)
      .then(response => {
        const token = response.data.token
        
        localStorage.setItem('survey-token', token)
        
        // set auth token
        Vue.axios.defaults.headers.common['Authorization'] = 'Bearer ' + token
                
        commit('auth_success', response.data)
        
        resolve(response)
      }).catch(function (error) {
        commit('auth_error')
        localStorage.removeItem('survey-token')
        
        reject(error)
      })
    })
  },
  tokenLogin({commit}) {
    return new Promise((resolve, reject) => {
      commit('auth_request')
      
      // set our default Authorization header
      Vue.axios.defaults.headers.common['Authorization'] = 'Bearer ' + localStorage.getItem('survey-token')
      
      Vue.axios.get('/user/info')
      .then(response => {
        
        commit('auth_success', response.data)
        
        resolve(response)
      }).catch(function (error) {
        commit('logout')
        
        localStorage.removeItem('survey-token')
        
        reject(error)
      })
    })
  },
  logout({commit}) {
    return new Promise((resolve) => {
      commit('logout')
      
      localStorage.removeItem('survey-token')
      
      // remove auth token
      delete Vue.axios.defaults.headers.common['Authorization']
            
      resolve()
    })
  },
  

}

// mutations
const mutations = {
  auth_request(state) {
    state.status = 'loading'
  },
  auth_success(state, payload) {
    state.username = payload.Username
    state.isAdmin = payload.isAdmin
    
    if (payload.token) {
      state.token = payload.token
    }
    
    state.status = 'success'    
  },
  auth_error(state) {
    state.status = 'error'
  },
  logout(state) {
    state.status = ''
    state.token = ''
    state.username = ''
    state.isAdmin = false
  },
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}