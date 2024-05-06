import Vue from 'vue'

// initial state
const state = {
  status: '',
  surveys: []
}

// getters
const getters = {

}

// actions
const actions = {
  fetchSurveys({commit}) {
    return new Promise((resolve, reject) => {
      commit('fetch_request')
      
      Vue.axios.get('/survey/list')
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
    if (payload.surveys) {
      state.surveys = payload.surveys
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