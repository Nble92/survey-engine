import Vue from 'vue'

// initial state
const state = {
  status: '',
  surveyType: '',
  config: {},
  surveyTypes: [
    {type: 'linear', label: 'Single Linear Survey'},
    {type: 'calendar', label: 'Calendar Based Repeated Measures'},
  ]
}

// getters
const getters = {
  guestSurveys: state => {
    if (state.config && state.config.allowGuests && state.config.allowGuests == true) {
      return true
    }
    return false
  }
}

// actions
const actions = {
  fetchState({commit}) {
    return new Promise((resolve, reject) => {
      commit('fetch_request')
      
      Vue.axios.get('/survey/settings')
      .then(response => {
        commit('load_state', response.data)
        
        resolve(response)
      }).catch(function (error) {
        commit('fetch_error')
        
        reject(error)
      })
    })
  },
  update({commit}, settings) {
    return new Promise((resolve, reject) => {
      commit('fetch_request')
      
      Vue.axios.post('/survey/settings', {type: settings.surveyType, config: settings.config})
      .then(response => {
        commit('load_state', response.data)
        
        resolve(response)
      }).catch(function (error) {
        commit('fetch_error')
        
        reject(error)
      })
    })
  }
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
    console.log('surveySettings', payload)
    
    state.surveyType = payload.type
    state.config = payload.config
    
    state.status = 'success'
  }
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}