import Vue from 'vue'

// initial state
const state = {
  status: '',
  logicFunctions: []
}

// getters
const getters = {

}

// actions
const actions = {
  fetchLogicFunctions({commit}) {
    return new Promise((resolve, reject) => {
      commit('fetch_request')
      
      Vue.axios.get('/logic/all')
      .then(response => {        
        commit('load_state', response.data)
        
        resolve(response)
      }).catch(function (error) {
        commit('fetch_error')
        
        reject(error)
      })
    })
  },
  
}

// mutations
const mutations = {
  fetch_request(state) {
    state.status = 'loading'
  },
  fetch_error(state) {
    state.status = 'error'
  },
  load_state(state, payload) {
    // eslint-disable-next-line
    console.log('logic: ', payload)
    
    if (payload.logicFunctions) {
      state.logicFunctions = payload.logicFunctions
    }
  }
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}