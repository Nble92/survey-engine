import Vue from 'vue'

// initial state
const state = {
  status: '',
  images: []
}

// getters
const getters = {

}

// actions
const actions = {
  fetchImages({commit}) {
    return new Promise((resolve, reject) => {
      commit('fetch_request')
      
      Vue.axios.get('/images/all')
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
    console.log('images: ', payload)
    
    if (payload.images) {
      state.images = payload.images
    } else {
      state.images = []
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